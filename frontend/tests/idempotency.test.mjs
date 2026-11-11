import test from 'node:test'
import assert from 'node:assert/strict'
import { checkoutKey, checkoutSignature } from '../src/features/orders/idempotency.ts'
const items = [{ offer_id: 3, quantity: 1, slug: 'x', title: 'x', seller: 'a', platform: 'PC', region: 'Global', unit_price_cents: 100, stock: 2 }]
test('same cart and payment reuse checkout key even across component reload', () => {
  const memory = new Map()
  const store = { getItem: k => memory.get(k) ?? null, setItem: (k, v) => memory.set(k, v) }
  const a = checkoutKey(items, 'approve', store)
  const b = checkoutKey(items, 'approve', store)
  assert.equal(a, b)
  const c = checkoutKey(items, 'decline', store)
  assert.notEqual(c, a)
})
test('fingerprint ignores cart display-only price', () => {
  assert.equal(checkoutSignature(items, 'approve'), checkoutSignature([{ ...items[0], unit_price_cents: 999 }], 'approve'))
})
