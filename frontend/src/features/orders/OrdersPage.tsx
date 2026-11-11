import { useEffect, useState } from 'react'
import { api } from '../../shared/api'
import type { Order } from '../../shared/types'
import { euros } from '../../shared/money'
import { useAuth } from '../auth/AuthContext'
export default function OrdersPage() {
  const { user, ready } = useAuth()
  const [items, setItems] = useState<Order[] | null>(null)
  const [error, setError] = useState('')
  useEffect(() => {
    if (!user) return
    const c = new AbortController()
    api<{
      items: Order[]
    }>('/orders', { signal: c.signal })
      .then((d) => setItems(d.items))
      .catch((e) => {
        if (!c.signal.aborted) setError(e.message)
      })
    return () => c.abort()
  }, [user])
  if (!ready) return <div className="page-container loading">Checking your account…</div>
  if (!user)
    return (
      <div className="page-container narrow">
        <div className="empty">
          <h2>Orders are private</h2>
          <p>Sign in to view your own purchases.</p>
          <a className="button button-yellow" href="/login?next=%2Forders">
            Sign in
          </a>
        </div>
      </div>
    )
  return (
    <div className="page-container narrow">
      <div className="page-heading">
        <span className="eyebrow">YOUR COLLECTION</span>
        <h1>
          Order history<span className="yellow-dot">.</span>
        </h1>
        <p>Completed demo purchases and your delivered keys.</p>
      </div>
      {error && <div className="message error">{error}</div>}
      {!items && !error && <div className="loading">Loading orders…</div>}
      {items?.length === 0 && (
        <div className="empty">
          <h3>No purchases yet</h3>
          <a href="/store" className="button button-yellow">
            Browse the store
          </a>
        </div>
      )}
      {items?.map((o) => (
        <article className="history-order" key={o.id}>
          <div className="history-head">
            <div>
              <strong>Order #{o.id}</strong>
              <small>{new Date(o.created_at).toLocaleString()}</small>
            </div>
            <div>
              <span className={`order-status ${o.status}`}>
                {o.status === 'paid' ? 'Paid' : 'Declined'}
              </span>
              <strong>{euros(o.total_cents)}</strong>
            </div>
          </div>
          <div className="history-lines">
            {o.items.map((item) => (
              <div key={item.id} className="history-line">
                <div>
                  <a href={`/product/${item.slug}`}>{item.title}</a>
                  <small>Sold by {item.seller}</small>
                </div>
                <span>{euros(item.price_cents)}</span>
                {item.key_code && (
                  <div className="delivered-key">
                    <span>DEMONSTRATION KEY</span>
                    <code>{item.key_code}</code>
                    <button onClick={() => navigator.clipboard?.writeText(item.key_code!)}>
                      Copy
                    </button>
                  </div>
                )}
              </div>
            ))}
          </div>
        </article>
      ))}
    </div>
  )
}
