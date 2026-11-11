import { useEffect, useState } from 'react'
import { api } from '../../shared/api'
import { euros } from '../../shared/money'
import type { ProductDetail } from '../../shared/types'
import { ProductArt } from '../catalog/ProductCard'
import { useCart } from '../cart/CartContext'
export default function ProductPage({ slug }: { slug: string }) {
  const [product, setProduct] = useState<ProductDetail | null>(null)
  const [error, setError] = useState('')
  const [added, setAdded] = useState(false)
  const { add } = useCart()
  useEffect(() => {
    const c = new AbortController()
    api<ProductDetail>(`/products/${encodeURIComponent(slug)}`, { signal: c.signal })
      .then(setProduct)
      .catch((e) => {
        if (!c.signal.aborted) setError(e.message)
      })
    return () => c.abort()
  }, [slug])
  if (error)
    return (
      <div className="page-container">
        <div className="message error">{error}</div>
        <a href="/store">← Back to catalog</a>
      </div>
    )
  if (!product) return <div className="page-container loading">Loading game details…</div>
  return (
    <div className="page-container product-page">
      <div className="breadcrumbs">
        <a href="/store">Store</a> / {product.platform} / {product.title}
      </div>
      <div className="product-main">
        <div className="product-visual">
          <ProductArt product={product} large />
        </div>
        <div className="product-info">
          <span className="eyebrow">DIGITAL GAME KEY</span>
          <h1>{product.title}</h1>
          <p className="subtitle">{product.description}</p>
          <div className="specs">
            <div>
              <small>Platform</small>
              <strong>{product.platform}</strong>
            </div>
            <div>
              <small>Activation region</small>
              <strong>{product.region}</strong>
            </div>
            <div>
              <small>Category</small>
              <strong>{product.genre}</strong>
            </div>
            <div>
              <small>Delivery</small>
              <strong>After successful demo payment</strong>
            </div>
          </div>
          <div className="price-feature">
            <span>Starting from</span>
            <strong>{euros(product.lowest_price_cents)}</strong>
            <small>
              {product.seller_count} seller offers · {product.stock_count} keys in stock
            </small>
          </div>
        </div>
      </div>
      <section className="offers-section">
        <div className="section-header">
          <div>
            <span className="eyebrow">CHOOSE YOUR OFFER</span>
            <h2>Available sellers</h2>
            <p>Prices are shown in EUR. All products and sellers are demonstration data.</p>
          </div>
        </div>
        <div className="offers-list">
          {product.offers.map((o, i) => (
            <div className="offer-row" key={o.id}>
              <div className="seller-icon">{o.seller_name.slice(0, 1)}</div>
              <div className="seller-name">
                <strong>{o.seller_name}</strong>
                <small>
                  {o.stock > 0 ? `${o.stock} keys available` : 'Out of stock'}{' '}
                  {i === 0 && o.stock > 0 ? ' · Best price' : ''}
                </small>
              </div>
              <div className="offer-price">{euros(o.price_cents)}</div>
              <button
                disabled={o.stock === 0}
                className="button button-yellow"
                onClick={() => {
                  add(product, o)
                  setAdded(true)
                }}
              >
                {o.stock === 0 ? 'Sold out' : 'Add to cart +'}
              </button>
            </div>
          ))}
        </div>
        {added && (
          <div className="message success">
            Added to your cart. <a href="/cart">View cart →</a>
          </div>
        )}
      </section>
    </div>
  )
}
