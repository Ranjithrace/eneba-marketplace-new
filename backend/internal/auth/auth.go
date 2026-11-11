package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"eneba-marketplace/backend/internal/httpjson"
	"golang.org/x/crypto/bcrypt"
)

const cookieName = "marketplace_session"
const maxAge = 7 * 24 * time.Hour

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Service struct {
	DB            *sql.DB
	SecureCookies bool
}
type contextKey struct{}

func userFrom(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(contextKey{}).(User)
	return u, ok
}
func UserFrom(ctx context.Context) (User, bool) { return userFrom(ctx) }
func normalize(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	p, e := mail.ParseAddress(email)
	if e != nil || p.Address != email || len(email) > 254 || !strings.Contains(email, ".") {
		return "", errors.New("invalid email")
	}
	return email, nil
}
func Validate(email, password string) (string, error) {
	n, err := normalize(email)
	if err != nil {
		return "", err
	}
	if len(password) < 10 || len(password) > 72 {
		return "", errors.New("password must have 10 to 72 characters")
	}
	return n, nil
}
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
func (s Service) cookie(w http.ResponseWriter, value string, maxAgeSeconds int) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: value, Path: "/", HttpOnly: true, Secure: s.SecureCookies, SameSite: http.SameSiteLaxMode, MaxAge: maxAgeSeconds})
}
func (s Service) startSession(w http.ResponseWriter, r *http.Request, id int64) error {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return err
	}
	raw := hex.EncodeToString(token)
	_, err := s.DB.ExecContext(r.Context(), `INSERT INTO auth_sessions(user_id,token_hash,expires_at) VALUES($1,$2,$3)`, id, hashToken(raw), time.Now().UTC().Add(maxAge))
	if err != nil {
		return err
	}
	s.cookie(w, raw, int(maxAge.Seconds()))
	return nil
}
func (s Service) Register(w http.ResponseWriter, r *http.Request) {
	var c Credentials
	if !httpjson.Decode(w, r, &c) {
		return
	}
	email, err := Validate(c.Email, c.Password)
	if err != nil {
		httpjson.Error(w, 400, err.Error())
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		httpjson.Error(w, 500, "Could not create account")
		return
	}
	var id int64
	err = s.DB.QueryRowContext(r.Context(), `INSERT INTO users(email,password_hash) VALUES($1,$2) ON CONFLICT DO NOTHING RETURNING id`, email, string(hashed)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		httpjson.Error(w, 409, "Email already registered")
		return
	}
	if err != nil {
		log.Printf("register: %v", err)
		httpjson.Error(w, 500, "Could not register")
		return
	}
	if err = s.startSession(w, r, id); err != nil {
		log.Printf("register session: %v", err)
		httpjson.Error(w, 500, "Could not start session")
		return
	}
	httpjson.Write(w, 201, User{ID: id, Email: email})
}
func (s Service) Login(w http.ResponseWriter, r *http.Request) {
	var c Credentials
	if !httpjson.Decode(w, r, &c) {
		return
	}
	email, err := normalize(c.Email)
	if err != nil {
		httpjson.Error(w, 401, "Invalid email or password")
		return
	}
	var id int64
	var hashed string
	err = s.DB.QueryRowContext(r.Context(), `SELECT id,password_hash FROM users WHERE LOWER(email)=$1`, email).Scan(&id, &hashed)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hashed), []byte(c.Password)) != nil {
		httpjson.Error(w, 401, "Invalid email or password")
		return
	}
	if err = s.startSession(w, r, id); err != nil {
		log.Printf("login session: %v", err)
		httpjson.Error(w, 500, "Could not start session")
		return
	}
	httpjson.Write(w, 200, User{ID: id, Email: email})
}
func (s Service) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		_, _ = s.DB.ExecContext(r.Context(), `DELETE FROM auth_sessions WHERE token_hash=$1`, hashToken(c.Value))
	}
	s.cookie(w, "", -1)
	httpjson.Write(w, 200, map[string]bool{"ok": true})
}
func (s Service) Authenticate(r *http.Request) (User, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil || len(c.Value) != 64 {
		return User{}, false
	}
	var u User
	err = s.DB.QueryRowContext(r.Context(), `SELECT u.id,u.email FROM auth_sessions a JOIN users u ON u.id=a.user_id WHERE a.token_hash=$1 AND a.expires_at>NOW()`, hashToken(c.Value)).Scan(&u.ID, &u.Email)
	return u, err == nil
}
func (s Service) Me(w http.ResponseWriter, r *http.Request) {
	u, ok := s.Authenticate(r)
	if !ok {
		httpjson.Error(w, 401, "Please sign in")
		return
	}
	httpjson.Write(w, 200, u)
}
func (s Service) Required(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := s.Authenticate(r)
		if !ok {
			httpjson.Error(w, 401, "Please sign in")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, u)))
	}
}
