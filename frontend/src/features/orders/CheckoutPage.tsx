import { useState } from 'react'
import { api } from '../../shared/api'
import { useAuth } from '../auth/AuthContext'
import { useCart } from '../cart/CartContext'
import { cartTotal, itemCount } from '../cart/cartLogic'
import { euros } from '../../shared/money'
import type { Order } from '../../shared/types'
import { checkoutKey, clearCheckoutKey } from './idempotency'
export default function CheckoutPage() {
  const { user, ready } = useAuth()
  const { items, clear } = useCart()
  const [method, setMethod] = useState<'approve' | 'decline'>('approve')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [order, setOrder] = useState<Order | null>(null)
  const [key, setKey] = useState(() => checkoutKey(items, method))
  if (!ready) return <div className="page-container loading">Checking your account…</div>
  if (!user)
    return (
      <div className="page-container narrow">
        <div className="empty">
          <h2>Sign in to continue</h2>
          <p>Checkout is available only to registered customers.</p>
          <a className="button button-yellow" href="/login?next=%2Fcheckout">
            Sign in to checkout
          </a>
        </div>
      </div>
    )
  if (order)
    return (
      <div className="page-container narrow">
        <div className={`checkout-result ${order.status === 'paid' ? 'successful' : 'declined'}`}>
          <span className="result-icon">{order.status === 'paid' ? '✓' : '!'}</span>
          <h1>{order.status === 'paid' ? 'Purchase successful!' : 'Payment declined'}</h1>
          <p>
            Order #{order.id} · {euros(order.total_cents)}
          </p>
          <p>
            {order.status === 'paid'
              ? 'Your demonstration keys are ready. You can view them in your order history.'
              : 'This was a simulated declined payment. No keys were issued and no money changed hands.'}
          </p>
          <a href="/orders" className="button button-yellow">
            View order history →
          </a>
        </div>
        {order.status === 'paid' && (
          <div className="keys-list">
            {order.items.map((item) => (
              <div key={item.id}>
                <strong>{item.title}</strong>
                <code>{item.key_code}</code>
              </div>
            ))}
          </div>
        )}
      </div>
    )
  if (items.length === 0)
    return (
      <div className="page-container narrow">
        <div className="empty">
          <h2>Your cart is empty</h2>
          <a className="button button-yellow" href="/store">
            Browse games
          </a>
        </div>
      </div>
    )
  async function checkout() {
    setError('')
    setBusy(true)
    try {
      const result = await api<Order>('/checkout', {
        method: 'POST',
        body: JSON.stringify({
          items: items.map((i) => ({ offer_id: i.offer_id, quantity: i.quantity })),
          payment_method: method,
          idempotency_key: key,
        }),
      })
      setOrder(result)
      clearCheckoutKey()
      if (result.status === 'paid') clear()
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Checkout failed')
    } finally {
      setBusy(false)
    }
  }
  function switchMethod(value: 'approve' | 'decline') {
    setMethod(value)
    setKey(checkoutKey(items, value))
    setError('')
  }
  return (
    <div className="page-container narrow">
      <div className="page-heading">
        <span className="eyebrow">SECURE DEMONSTRATION CHECKOUT</span>
        <h1>
          Almost there<span className="yellow-dot">.</span>
        </h1>
        <p>Signed in as {user.email}</p>
      </div>
      <div className="checkout-layout">
        <div className="checkout-main">
          <h2>Payment simulation</h2>
          <p>Select a result to test the checkout flow. No payment card or real money is used.</p>
          <label className={`payment-option ${method === 'approve' ? 'selected' : ''}`}>
            <input
              type="radio"
              checked={method === 'approve'}
              onChange={() => switchMethod('approve')}
            />
            <span>
              <strong>Simulate successful payment</strong>
              <small>Deliver the demonstration keys</small>
            </span>
            <span>✓</span>
          </label>
          <label className={`payment-option ${method === 'decline' ? 'selected' : ''}`}>
            <input
              type="radio"
              checked={method === 'decline'}
              onChange={() => switchMethod('decline')}
            />
            <span>
              <strong>Simulate declined payment</strong>
              <small>Record a failed order without delivering keys</small>
            </span>
            <span>✕</span>
          </label>
          <div className="notice">
            All product keys are fictional and cannot be redeemed. Stock and prices are rechecked on
            the server.
          </div>
        </div>
        <aside className="order-summary">
          <h3>Checkout summary</h3>
          <p>
            {itemCount(items)} keys from {items.length} offers
          </p>
          {items.map((item) => (
            <div className="summary-line" key={item.offer_id}>
              <span>
                {item.title} × {item.quantity}
              </span>
              <strong>{euros(item.unit_price_cents * item.quantity)}</strong>
            </div>
          ))}
          <div className="summary-total">
            <span>Estimated total</span>
            <strong>{euros(cartTotal(items))}</strong>
          </div>
          {error && (
            <div role="alert" className="message error">
              {error}{' '}
              <small>For a network retry, this checkout reuses the same request identifier.</small>
            </div>
          )}
          <button disabled={busy} onClick={checkout} className="button button-yellow full">
            {busy ? 'Processing…' : `Complete demo payment →`}
          </button>
          <small>Final pricing is calculated by the backend.</small>
        </aside>
      </div>
    </div>
  )
}
