import test from 'node:test'
import assert from 'node:assert/strict'
import { euros } from '../src/shared/money.ts'
test('format currency from integer cents',()=>{assert.equal(euros(714),'€7.14');assert.equal(euros(0),'€0.00')})
test('missing prices indicate unavailability',()=>{assert.equal(euros(null),'Unavailable')})
