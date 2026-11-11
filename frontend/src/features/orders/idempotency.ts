import type { CartItem } from '../../shared/types'
const STORAGE = 'keymarket-checkout-intent-v1'
export function checkoutSignature(items: CartItem[], payment: 'approve' | 'decline'): string {
  return JSON.stringify({
    payment,
    lines: items.map((i) => ({ id: i.offer_id, quantity: i.quantity })).sort((a, b) => a.id - b.id),
  })
}
export function checkoutKey(
  items: CartItem[],
  payment: 'approve' | 'decline',
  store: Pick<Storage, 'getItem' | 'setItem'> | null = typeof sessionStorage === 'undefined'
    ? null
    : sessionStorage,
): string {
  const signature = checkoutSignature(items, payment)
  try {
    const old = store?.getItem(STORAGE)
    if (old) {
      const saved = JSON.parse(old) as {
        signature?: string
        key?: string
      }
      if (saved.signature === signature && typeof saved.key === 'string' && saved.key.length >= 12)
        return saved.key
    }
  } catch {
    /* session storage can be disabled */
  }
  const key = crypto.randomUUID()
  try {
    store?.setItem(STORAGE, JSON.stringify({ signature, key }))
  } catch {
    /* storage unavailable */
  }
  return key
}
export function clearCheckoutKey(): void {
  try {
    sessionStorage.removeItem(STORAGE)
  } catch {
    /* storage unavailable */
  }
}
