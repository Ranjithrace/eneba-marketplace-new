import { createContext, useContext, useEffect, useState } from 'react'
import type { CartItem, Offer, ProductDetail } from '../../shared/types'
import { addItem, changeQuantity, loadCart, STORAGE_KEY } from './cartLogic'
interface CartContextValue {
  items: CartItem[]
  add: (p: ProductDetail, o: Offer) => void
  setQuantity: (id: number, qty: number) => void
  clear: () => void
}
const CartContext = createContext<CartContextValue | null>(null)
export function CartProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<CartItem[]>(loadCart)
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(items))
    } catch {
      /* storage unavailable */
    }
  }, [items])
  function add(p: ProductDetail, o: Offer) {
    setItems((current) =>
      addItem(current, {
        offer_id: o.id,
        slug: p.slug,
        title: p.title,
        seller: o.seller_name,
        platform: p.platform,
        region: p.region,
        cover_image_url: p.cover_image_url,
        unit_price_cents: o.price_cents,
        quantity: 1,
        stock: o.stock,
      }),
    )
  }
  return (
    <CartContext.Provider
      value={{
        items,
        add,
        setQuantity: (id, qty) => setItems((current) => changeQuantity(current, id, qty)),
        clear: () => setItems([]),
      }}
    >
      {children}
    </CartContext.Provider>
  )
}
export function useCart() {
  const ctx = useContext(CartContext)
  if (!ctx) throw new Error('CartProvider missing')
  return ctx
}
