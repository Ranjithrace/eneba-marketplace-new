import test from 'node:test';
import assert from 'node:assert/strict';
import { artworkUrl } from '../src/shared/artworkUrl.ts';

test('artwork is supplied by the API, not derived from a slug', () => {
  assert.equal(artworkUrl(''), null);
  assert.equal(artworkUrl(null), null);
  assert.equal(artworkUrl('https://cdn.example.org/cover.jpg'), 'https://cdn.example.org/cover.jpg');
});

test('unsafe and malformed image links are not rendered', () => {
  assert.equal(artworkUrl('javascript:alert(1)'), null);
  assert.equal(artworkUrl('http://example.org/image.jpg'), null);
  assert.equal(artworkUrl('not-a-url'), null);
});
