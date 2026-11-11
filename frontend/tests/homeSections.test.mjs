import test from 'node:test';
import assert from 'node:assert/strict';
import { loadHomeBlock } from '../src/features/home/loadHomeBlock.ts';

test('one failed homepage request does not hide successful sections', async () => {
  const view = {};
  await Promise.all([
    loadHomeBlock(async () => ({ heading: 'Banner OK' }), x => { view.banner = x.heading; }, e => { view.bannerError = e; }),
    loadHomeBlock(async () => { throw new Error('Best sellers unavailable'); }, x => { view.best = x.heading; }, e => { view.bestError = e; }),
    loadHomeBlock(async () => ({ heading: 'Deals OK' }), x => { view.deals = x.heading; }, e => { view.dealsError = e; }),
  ]);
  assert.deepEqual(view, {banner: 'Banner OK', bestError: 'Best sellers unavailable', deals: 'Deals OK'});
});

test('late result is ignored after navigation away', async () => {
  let updated = false;
  await loadHomeBlock(async () => 1, () => { updated = true; }, () => { updated = true; }, () => true);
  assert.equal(updated, false);
});
