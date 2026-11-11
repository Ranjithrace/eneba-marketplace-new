import { useEffect, useState } from 'react'
import { useAuth } from '../features/auth/AuthContext'
import { useCart } from '../features/cart/CartContext'
import { itemCount } from '../features/cart/cartLogic'
import HomePage from '../features/home/HomePage'
import StorePage from '../features/catalog/StorePage'
import ProductPage from '../features/product/ProductPage'
import CartPage from '../features/cart/CartPage'
import AuthPage from '../features/auth/AuthPage'
import CheckoutPage from '../features/orders/CheckoutPage'
import OrdersPage from '../features/orders/OrdersPage'
import { api } from '../shared/api'
import { storeSearch } from '../shared/navigation'
import Icon from '../shared/Icon'
import type { Product } from '../shared/types'

function Brand() {
  return (
    <a className="logo" href="/" aria-label="Eneba internship demo home">
      <span className="logo-mark">
        <Icon name="gamepad" />
      </span>
      <span className="logo-word">eneba</span>
      <span className="demo-tag">DEMO</span>
    </a>
  )
}

function Header() {
  const { user, logout } = useAuth()
  const { items } = useCart()
  const [q, setQ] = useState(new URLSearchParams(window.location.search).get('q') || '')
  const [suggest, setSuggest] = useState<Product[]>([])
  const [open, setOpen] = useState(false)
  useEffect(() => {
    const c = new AbortController()
    if (q.trim().length < 2) {
      setSuggest([])
      return () => c.abort()
    }
    const id = setTimeout(
      () =>
        api<{ items: Product[] }>(`/search/suggest?q=${encodeURIComponent(q.trim())}`, {
          signal: c.signal,
        })
          .then((d) => setSuggest(d.items))
          .catch(() => {
            if (!c.signal.aborted) setSuggest([])
          }),
      220,
    )
    return () => {
      clearTimeout(id)
      c.abort()
    }
  }, [q])
  return (
    <>
      <div className="announcement">
        INTERNSHIP DEMO <span>· Simulated payments and non-redeemable keys only</span>
      </div>
      <header className="site-header">
        <div className="header-inner">
          <Brand />
          <form
            className="global-search"
            role="search"
            onSubmit={(e) => {
              e.preventDefault()
              window.location.assign(storeSearch(q))
            }}
          >
            <Icon name="search" className="search-icon" />
            <input
              value={q}
              onFocus={() => setOpen(true)}
              onBlur={() => setTimeout(() => setOpen(false), 150)}
              onChange={(e) => {
                setQ(e.target.value)
                setOpen(true)
              }}
              placeholder="Search games, editions, platforms..."
              aria-label="Search games"
              maxLength={100}
              autoComplete="off"
            />
            <button type="submit" aria-label="Submit search">
              <Icon name="arrow" />
            </button>
            {open && suggest.length > 0 && (
              <div className="suggestions">
                {suggest.slice(0, 5).map((p) => (
                  <a key={p.id} href={`/product/${encodeURIComponent(p.slug)}`}>
                    <Icon name="search" />
                    <span>
                      {p.title}
                      <small>
                        {p.platform} · {p.region}
                      </small>
                    </span>
                  </a>
                ))}
              </div>
            )}
          </form>
          <div className="header-actions">
            {user ? (
              <>
                <a className="header-action" href="/orders">
                  <Icon name="user" />
                  <span>My orders</span>
                </a>
                <button
                  className="header-action signout-action"
                  type="button"
                  onClick={() => logout().then(() => window.location.assign('/'))}
                  title="Sign out"
                >
                  <Icon name="logout" />
                  <span>Sign out</span>
                </button>
              </>
            ) : (
              <a className="header-action" href="/login">
                <Icon name="user" />
                <span>Sign in</span>
              </a>
            )}
            <a
              className="header-action cart-link"
              href="/cart"
              aria-label={`Cart with ${itemCount(items)} items`}
            >
              <Icon name="cart" />
              <span className="cart-count">{itemCount(items)}</span>
              <span className="cart-text">Cart</span>
            </a>
          </div>
        </div>
        <nav className="nav-row" aria-label="Marketplace navigation">
          <div className="nav-inner">
            <a className={window.location.pathname === '/' ? 'active' : ''} href="/">
              Discover
            </a>
            <a
              className={
                window.location.pathname === '/store' &&
                !window.location.search.includes('platform=')
                  ? 'active'
                  : ''
              }
              href="/store"
            >
              All games
            </a>
            <a href="/store?platform=PC">PC games</a>
            <a href="/store?platform=Xbox">Xbox games</a>
            <a href="/store?platform=Nintendo+Switch">Switch games</a>
            <a href="/store?platform=PlayStation+5">PS5 games</a>
            <a className="nav-featured" href="/#bestsellers">
              Best sellers
            </a>
          </div>
        </nav>
      </header>
    </>
  )
}

function Content() {
  const path = decodeURIComponent(window.location.pathname)
  if (path === '/') return <HomePage />
  if (path === '/store') return <StorePage />
  if (path.startsWith('/product/')) return <ProductPage slug={path.substring('/product/'.length)} />
  if (path === '/cart') return <CartPage />
  if (path === '/login') return <AuthPage register={false} />
  if (path === '/register') return <AuthPage register />
  if (path === '/checkout') return <CheckoutPage />
  if (path === '/orders') return <OrdersPage />
  return (
    <div className="page-container empty">
      <h1>Page not found</h1>
      <a href="/store">Return to the store</a>
    </div>
  )
}

export default function App() {
  return (
    <div className="app-shell">
      <Header />
      <main>
        <Content />
      </main>
      <footer className="site-footer">
        <div className="footer-inner">
          <div>
            <Brand />
            <p>
              An independent internship demonstration inspired by Eneba's interface. Not an official
              Eneba website.
            </p>
          </div>
          <div>
            <strong>Explore</strong>
            <a href="/store">All games</a>
            <a href="/#bestsellers">Best sellers</a>
            <a href="/cart">Your cart</a>
            <a href="/orders">Order history</a>
          </div>
          <div>
            <strong>Demo only</strong>
            <p>
              All sellers, prices and keys are fictional. No real payments or redeemable game keys.
            </p>
          </div>
        </div>
        <div className="footer-bottom">
          Internship project · React + Go + PostgreSQL · Game art belongs to its respective owners
        </div>
      </footer>
    </div>
  )
}
