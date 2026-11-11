import { useEffect, useState } from 'react'
import { api } from '../../shared/api'
import type { CatalogPage } from '../../shared/types'
import ProductCard from './ProductCard'
import Icon from '../../shared/Icon'
const platforms = ['', 'PC', 'Nintendo Switch', 'Xbox', 'PlayStation 5']
const regions = ['', 'Global', 'Europe']
export default function StorePage() {
  const query = new URLSearchParams(window.location.search)
  const [q, setQ] = useState(query.get('q') || '')
  const [platform, setPlatform] = useState(query.get('platform') || '')
  const [region, setRegion] = useState(query.get('region') || '')
  const [min, setMin] = useState(
    query.get('min_price') ? String(Number(query.get('min_price')) / 100) : '',
  )
  const [max, setMax] = useState(
    query.get('max_price') ? String(Number(query.get('max_price')) / 100) : '',
  )
  const [sort, setSort] = useState(query.get('sort') || 'popular')
  const [page, setPage] = useState(Math.max(1, Number(query.get('page')) || 1))
  const [data, setData] = useState<CatalogPage | null>(null)
  const [error, setError] = useState('')
  function apply(pageNumber = 1) {
    const params = new URLSearchParams()
    if (q.trim()) params.set('q', q.trim())
    if (platform) params.set('platform', platform)
    if (region) params.set('region', region)
    if (min && Number(min) >= 0) params.set('min_price', String(Math.round(Number(min) * 100)))
    if (max && Number(max) >= 0) params.set('max_price', String(Math.round(Number(max) * 100)))
    params.set('sort', sort)
    params.set('page', String(pageNumber))
    window.location.assign(`/store?${params.toString()}`)
  }
  useEffect(() => {
    const c = new AbortController()
    setError('')
    api<CatalogPage>(`/products${window.location.search || '?page=1'}`, { signal: c.signal })
      .then(setData)
      .catch((e) => {
        if (!c.signal.aborted) setError(e.message)
      })
    return () => c.abort()
  }, [])
  return (
    <div className="page-container store-page">
      <div className="page-heading compact-heading">
        <span className="eyebrow">THE GAME STORE</span>
        <h1>
          Find your next favourite<span className="yellow-dot">.</span>
        </h1>
        <p>Discover games, compare sellers, and choose your offer.</p>
      </div>
      <div className="store-layout">
        <aside className="filters">
          <div className="filters-head">
            <h3>
              <Icon name="filters" /> Filters
            </h3>
            <a href="/store">Clear all</a>
          </div>
          <form
            onSubmit={(e) => {
              e.preventDefault()
              apply()
            }}
          >
            <label>
              Search games
              <input
                value={q}
                onChange={(e) => setQ(e.target.value)}
                placeholder="Game name"
                maxLength={100}
              />
            </label>
            <label>
              Platform
              <select value={platform} onChange={(e) => setPlatform(e.target.value)}>
                {platforms.map((p) => (
                  <option key={p} value={p}>
                    {p || 'All platforms'}
                  </option>
                ))}
              </select>
            </label>
            <label>
              Region
              <select value={region} onChange={(e) => setRegion(e.target.value)}>
                {regions.map((p) => (
                  <option key={p} value={p}>
                    {p || 'All regions'}
                  </option>
                ))}
              </select>
            </label>
            <div className="price-field">
              <label>
                Min price (€)
                <input
                  type="number"
                  min="0"
                  step="0.01"
                  value={min}
                  onChange={(e) => setMin(e.target.value)}
                  placeholder="0"
                />
              </label>
              <label>
                Max price (€)
                <input
                  type="number"
                  min="0"
                  step="0.01"
                  value={max}
                  onChange={(e) => setMax(e.target.value)}
                  placeholder="Any"
                />
              </label>
            </div>
            <button className="button button-yellow full" type="submit">
              Apply filters
            </button>
          </form>
        </aside>
        <section className="store-results">
          <div className="results-toolbar">
            <div>
              <strong>{data?.total ?? '…'} games found</strong>
              <span> {q && `matching “${q}”`}</span>
            </div>
            <label className="inline-sort">
              Sort by{' '}
              <select
                value={sort}
                onChange={(e) => {
                  setSort(e.target.value)
                  const u = new URLSearchParams(window.location.search)
                  u.set('sort', e.target.value)
                  u.set('page', '1')
                  window.location.assign(`/store?${u}`)
                }}
              >
                <option value="popular">Popular</option>
                <option value="price_asc">Price: low to high</option>
                <option value="price_desc">Price: high to low</option>
                <option value="title">Title A–Z</option>
                <option value="newest">Recently added</option>
              </select>
            </label>
          </div>
          {error && <div className="message error">{error}</div>}
          {!data && !error && <p className="loading">Finding games…</p>}
          {data && data.items.length === 0 && (
            <div className="empty">
              <h3>No games found</h3>
              <p>Try another search or clear your filters.</p>
              <a className="button button-outline" href="/store">
                Clear filters
              </a>
            </div>
          )}
          {data && (
            <>
              <div className="product-grid">
                {data.items.map((p) => (
                  <ProductCard key={p.id} product={p} />
                ))}
              </div>
              <div className="pagination">
                <button disabled={page <= 1} onClick={() => apply(page - 1)}>
                  ← Previous
                </button>
                <span>
                  Page {data.page} of {Math.max(1, data.pages)}
                </span>
                <button disabled={page >= data.pages} onClick={() => apply(page + 1)}>
                  Next →
                </button>
              </div>
            </>
          )}
        </section>
      </div>
    </div>
  )
}
