import { useState } from 'react';
import type { Product } from '../../shared/types';
import { euros } from '../../shared/money';
import { artworkUrl } from '../../shared/artworkUrl';

/** The image is optional: a matching real cover or a neutral title card, never fake art. */
export function ProductArt({ product, large = false }: {
  product: Pick<Product, 'title' | 'platform' | 'cover_image_url' | 'slug'>;
  large?: boolean;
}) {
  const [failed, setFailed] = useState(false);
  const src = artworkUrl(product.cover_image_url);
  return <div className={`product-art${large ? ' large' : ''}`}>
    {src && !failed ? <img src={src} alt={`${product.title} cover art`} loading={large ? 'eager' : 'lazy'} onError={() => setFailed(true)}/> :
      <div className="art-unavailable"><span className="art-title">{product.title}</span><small>Cover not available</small></div>}
    <span className="cover-platform">{product.platform}</span>
  </div>;
}

export default function ProductCard({ product }: { product: Product }) {
  return <a href={`/product/${encodeURIComponent(product.slug)}`} className="product-card">
    <ProductArt product={product}/>
    <div className="product-card-content">
      <h3 className="product-title">{product.title}</h3>
      <div className="product-meta">{product.platform} <span aria-hidden="true">·</span> {product.region}</div>
      <div className="product-card-bottom">
        <div className="price-block"><small>FROM</small><strong>{euros(product.lowest_price_cents)}</strong></div>
        <span className="seller-tag">{product.seller_count} {product.seller_count === 1 ? 'seller' : 'sellers'}</span>
      </div>
    </div>
  </a>;
}
