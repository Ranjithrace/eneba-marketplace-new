import test from 'node:test'
import assert from 'node:assert/strict'
import { addItem, changeQuantity, cartTotal, itemCount } from '../src/features/cart/cartLogic.ts'
const offer = {offer_id:12,slug:'demo',title:'Demo game',seller:'Seller A',platform:'PC',region:'Global',unit_price_cents:714,quantity:1,stock:2}
test('cart adds one item, clamps to available stock, totals cents accurately',()=>{
  let cart = addItem([],offer)
  cart=addItem(cart,offer)
  cart=addItem(cart,offer)
  assert.equal(cart.length,1)
  assert.equal(itemCount(cart),2)
  assert.equal(cartTotal(cart),1428)
  cart=changeQuantity(cart,12,50)
  assert.equal(cart[0].quantity,2)
  cart=changeQuantity(cart,12,0)
  assert.equal(cart.length,0)
})
test('sold out offer cannot enter cart',()=>{
 assert.deepEqual(addItem([],{...offer,stock:0}),[])
})

test('saved cart image URL is optional and adding a fresh API offer updates it',()=>{
 let cart=addItem([],{...offer,cover_image_url:'https://example.test/a.jpg'});
 assert.equal(cart[0].cover_image_url,'https://example.test/a.jpg');
 cart=addItem(cart,{...offer,cover_image_url:'https://example.test/b.jpg'});
 assert.equal(cart[0].cover_image_url,'https://example.test/b.jpg');
});
