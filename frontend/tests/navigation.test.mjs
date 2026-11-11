import test from 'node:test'
import assert from 'node:assert/strict'
import { safeReturn, storeSearch } from '../src/shared/navigation.ts'
test('safe post-login navigation refuses external redirects',()=>{
 assert.equal(safeReturn('//evil.example'),'/')
 assert.equal(safeReturn('https://evil.example'),'/')
 assert.equal(safeReturn('/checkout'),'/checkout')
})
test('global search encodes a partial-name query',()=>{
 assert.equal(storeSearch('  mine craft  '),'/store?q=mine%20craft')
 assert.equal(storeSearch(' '),'/store')
})
