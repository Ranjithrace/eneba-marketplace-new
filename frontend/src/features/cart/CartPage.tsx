import { useState } from 'react'
import { useCart } from './CartContext'
import { cartTotal, itemCount } from './cartLogic'
import { euros } from '../../shared/money'
import { artworkUrl } from '../../shared/artworkUrl'

function CartArtwork({ title, cover }: { title: string; cover?: string }) {
  const [failed, setFailed] = useState(false)
  const src = artworkUrl(cover)
  return (
    <div className="cart-symbol">
      {src && !failed ? (
        <img src={src} alt={`${title} cover`} loading="lazy" onError={() => setFailed(true)} />
      ) : (
        <span aria-label={title}>{title.charAt(0)}</span>
      )}
    </div>
  )
}

export default function CartPage() {
  const { items, setQuantity } = useCart()
  return (
    <div className="page-container narrow">
      <div className="page-heading">
        <span className="eyebrow">ONE STEP CLOSER</span>
        <h1>
          Your cart<span className="yellow-dot">.</span>
        </h1>
        <p>{itemCount(items)} digital keys in your cart</p>
      </div>
      {items.length === 0 ? (
        <div className="empty">
          <h3>Your cart is empty</h3>
          <p>Explore the marketplace to find your next game.</p>
          <a className="button button-yellow" href="/store">
            Explore games
          </a>
        </div>
      ) : (
        <div className="cart-layout">
          <div className="cart-items">
            {items.map((x) => (
              <div className="cart-row" key={x.offer_id}>
                <CartArtwork title={x.title} cover={x.cover_image_url} />
                <div className="cart-name">
                  <a href={`/product/${x.slug}`}>
                    <strong>{x.title}</strong>
                  </a>
                  <small>
                    {x.platform} · {x.region} · {x.seller}
                  </small>
                  <button className="text-button" onClick={() => setQuantity(x.offer_id, 0)}>
                    Remove
                  </button>
                </div>
                <div className="cart-qty">
                  <button
                    aria-label={`Decrease ${x.title}`}
                    onClick={() => setQuantity(x.offer_id, x.quantity - 1)}
                  >
                    −
                  </button>
                  <span>{x.quantity}</span>
                  <button
                    aria-label={`Increase ${x.title}`}
                    disabled={x.quantity >= x.stock || x.quantity >= 20}
                    onClick={() => setQuantity(x.offer_id, x.quantity + 1)}
                  >
                    +
                  </button>
                </div>
                <strong className="cart-price">{euros(x.unit_price_cents * x.quantity)}</strong>
              </div>
            ))}
          </div>
          <aside className="order-summary">
            <h3>Order summary</h3>
            <div className="summary-line">
              <span>Subtotal</span>
              <strong>{euros(cartTotal(items))}</strong>
            </div>
            <div className="summary-line">
              <span>Delivery</span>
              <strong>Free</strong>
            </div>
            <div className="summary-total">
              <span>Total</span>
              <strong>{euros(cartTotal(items))}</strong>
            </div>
            <a className="button button-yellow full" href="/checkout">
              Continue to checkout →
            </a>
            <small>Demo payments only. Prices and stock will be checked again at checkout.</small>
          </aside>
        </div>
      )}
    </div>
  )
}
