import type { CartItem } from '../../shared/types'
export const STORAGE_KEY = 'marketplace-demo-cart-v1'
export function addItem(current: CartItem[], item: CartItem): CartItem[] {
  if (item.stock < 1) return current
  const existing = current.find((row) => row.offer_id === item.offer_id)
  if (!existing) return [...current, { ...item, quantity: 1 }]
  return current.map((row) =>
    row.offer_id === item.offer_id
      ? {
          ...row,
          stock: item.stock,
          cover_image_url: item.cover_image_url ?? row.cover_image_url,
          unit_price_cents: item.unit_price_cents,
          quantity: Math.min(row.quantity + 1, item.stock, 20),
        }
      : row,
  )
}
export function changeQuantity(current: CartItem[], id: number, quantity: number): CartItem[] {
  if (quantity <= 0) return current.filter((item) => item.offer_id !== id)
  return current.map((item) =>
    item.offer_id === id
      ? { ...item, quantity: Math.max(1, Math.min(Math.floor(quantity), item.stock, 20)) }
      : item,
  )
}
export function cartTotal(items: CartItem[]): number {
  return items.reduce((sum, item) => sum + item.unit_price_cents * item.quantity, 0)
}
export function itemCount(items: CartItem[]): number {
  return items.reduce((sum, item) => sum + item.quantity, 0)
}
export function loadCart(): CartItem[] {
  try {
    const json = localStorage.getItem(STORAGE_KEY)
    if (!json) return []
    const parsed: unknown = JSON.parse(json)
    if (!Array.isArray(parsed)) return []
    return parsed
      .filter(
        (x): x is CartItem =>
          x !== null &&
          typeof x === 'object' &&
          Number.isSafeInteger(x.offer_id) &&
          x.offer_id > 0 &&
          typeof x.title === 'string' &&
          typeof x.slug === 'string' &&
          typeof x.seller === 'string' &&
          typeof x.unit_price_cents === 'number' &&
          x.unit_price_cents > 0 &&
          Number.isInteger(x.quantity) &&
          x.quantity > 0 &&
          x.quantity <= 20 &&
          Number.isInteger(x.stock) &&
          x.stock > 0,
      )
      .slice(0, 20)
  } catch {
    return []
  }
}
