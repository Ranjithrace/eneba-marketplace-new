import { useState } from 'react'
import { useAuth } from './AuthContext'
import { safeReturn } from '../../shared/navigation'
export default function AuthPage({ register }: { register: boolean }) {
  const { signIn } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  const after = safeReturn(new URLSearchParams(window.location.search).get('next'))
  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setBusy(true)
    try {
      await signIn(email, password, register)
      window.location.assign(after)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not sign in')
    } finally {
      setBusy(false)
    }
  }
  return (
    <div className="auth-shell">
      <div className="auth-box">
        <span className="eyebrow">YOUR PLAYER ACCOUNT</span>
        <h1>{register ? 'Join the marketplace.' : 'Welcome back.'}</h1>
        <p>
          {register
            ? 'Create an account to purchase and access demonstration keys.'
            : 'Sign in to view your purchases and complete checkout.'}
        </p>
        <form onSubmit={submit}>
          <label>
            Email address
            <input
              type="email"
              autoComplete="email"
              required
              maxLength={254}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
            />
          </label>
          <label>
            Password
            <input
              type="password"
              autoComplete={register ? 'new-password' : 'current-password'}
              required
              minLength={register ? 10 : 1}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="At least 10 characters"
            />
          </label>
          {error && (
            <div className="message error" role="alert">
              {error}
            </div>
          )}
          <button className="button button-yellow full" disabled={busy} type="submit">
            {busy ? 'Working…' : register ? 'Create account' : 'Sign in'} →
          </button>
        </form>
        <div className="auth-switch">
          {register ? 'Already have an account?' : 'New to the marketplace?'}{' '}
          <a href={`${register ? '/login' : '/register'}?next=${encodeURIComponent(after)}`}>
            {register ? 'Sign in' : 'Create account'}
          </a>
        </div>
        <small>For demo use only. No real payments or real game keys.</small>
      </div>
    </div>
  )
}
