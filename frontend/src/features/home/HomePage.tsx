import { useEffect, useState } from 'react';
import { api } from '../../shared/api';
import type { Banner, Collection } from '../../shared/types';
import ProductCard from '../catalog/ProductCard';
import { loadHomeBlock } from './loadHomeBlock';

/** Every section has independent loading/error state and remains backend-driven. */
export default function HomePage() {
  const [banner, setBanner] = useState<Banner | null>(null);
  const [best, setBest] = useState<Collection | null>(null);
  const [deals, setDeals] = useState<Collection | null>(null);
  const [bannerError, setBannerError] = useState('');
  const [bestError, setBestError] = useState('');
  const [dealsError, setDealsError] = useState('');

  useEffect(() => {
    const c = new AbortController();
    const canceled = () => c.signal.aborted;
    void loadHomeBlock(() => api<Banner>('/home/banner', {signal: c.signal}), setBanner, setBannerError, canceled);
    void loadHomeBlock(() => api<Collection>('/home/bestsellers', {signal: c.signal}), setBest, setBestError, canceled);
    void loadHomeBlock(() => api<Collection>('/home/deals', {signal: c.signal}), setDeals, setDealsError, canceled);
    return () => c.abort();
  }, []);

  useEffect(() => {
    const anchor = window.location.hash.slice(1);
    if ((anchor === 'bestsellers' && best) || (anchor === 'deals' && deals)) {
      document.getElementById(anchor)?.scrollIntoView({ block: 'start' });
    }
  }, [best, deals]);

  function collection(id: string, section: Collection | null, error: string, eyebrow: string) {
    return <section className="collection" id={id}>
      {error && <div className="message error" role="status">{id === 'bestsellers' ? 'Best sellers' : 'Deals'} could not load: {error}. Other homepage sections are still available.</div>}
      {!section && !error && <p className="loading">Loading {id === 'bestsellers' ? 'best sellers' : 'deals'}…</p>}
      {section && <><div className="section-header"><div><span className="eyebrow">{eyebrow}</span>
        <h2>{section.heading}</h2><p>{section.subtitle}</p></div>
        <a href="/store" className="text-link">View all products →</a></div>
        <div className="product-grid">{section.items.map(p => <ProductCard key={p.id} product={p}/>)}</div>
      </>}
    </section>;
  }

  return <>
    {banner && <section className="hero-banner"><div className="hero-copy"><span className="eyebrow">{banner.eyebrow}</span>
      <h1>{banner.heading}</h1><p>{banner.subtitle}</p><a className="button button-yellow" href={banner.href}>{banner.action} <span>↗</span></a>
      <div className="hero-trust">Demo checkout · Seller offer comparison · Non-redeemable keys</div></div>
      <div className="hero-decoration" aria-hidden="true"><div className="hero-halo"/><span className="hero-orbit">✦</span><span className="hero-controller">✧</span><div className="hero-sticker">DEMO<br/>STORE</div></div>
    </section>}
    {!banner && !bannerError && <p className="loading">Loading homepage banner…</p>}
    {bannerError && <div className="page-container"><div className="message error" role="status">Homepage banner could not load: {bannerError}</div></div>}
    {collection('bestsellers', best, bestError, 'DISCOVER')}
    {collection('deals', deals, dealsError, 'YOUR NEXT DEAL')}
    <section className="promise-row"><div><strong>⚡ Instant delivery</strong><span>Demonstration keys after a successful test purchase</span></div>
      <div><strong>◇ Multiple sellers</strong><span>Compare prices and available stock</span></div>
      <div><strong>♢ No real payment</strong><span>A safe internship demonstration marketplace</span></div></section>
  </>;
}
