export interface Product {
  id: number
  slug: string
  title: string
  description: string
  platform: string
  region: string
  genre: string
  cover_image_url: string
  seller_count: number
  lowest_price_cents: number | null
  stock_count: number
  sold_count: number
}
export interface Offer {
  id: number
  seller_id: number
  seller_name: string
  price_cents: number
  stock: number
}
export interface ProductDetail extends Product {
  offers: Offer[]
}
export interface CatalogPage {
  items: Product[]
  count: number
  total: number
  page: number
  page_size: number
  pages: number
}
export interface Banner {
  heading: string
  subtitle: string
  eyebrow: string
  action: string
  href: string
}
export interface Collection {
  heading: string
  subtitle: string
  items: Product[]
}
export interface User {
  id: number
  email: string
}
export interface CartItem {
  /** Optional API-supplied image URL, retained for backward-compatible saved carts. */
  cover_image_url?: string
  offer_id: number
  slug: string
  title: string
  seller: string
  platform: string
  region: string
  unit_price_cents: number
  quantity: number
  stock: number
}
export interface OrderItem {
  id: number
  offer_id: number
  title: string
  slug: string
  seller: string
  price_cents: number
  key_code: string | null
}
export interface Order {
  id: number
  status: 'paid' | 'failed' | 'pending'
  total_cents: number
  created_at: string
  items: OrderItem[]
}
