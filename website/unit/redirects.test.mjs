import assert from 'node:assert/strict';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test, { after } from 'node:test';
import { loadRedirects } from '../src/lib/redirects.mjs';

const created = [];
after(() => created.forEach((dir) => fs.rmSync(dir, { recursive: true, force: true })));

test('redirects cover plain and .md legacy routes under the base', () => {
  const directory = fs.mkdtempSync(path.join(os.tmpdir(), 'redirects-'));
  created.push(directory);
  const file = path.join(directory, 'r.json');
  fs.writeFileSync(file, JSON.stringify({ installation: 'start-here/installation/', README: '' }));
  assert.deepEqual(loadRedirects(file, '/site/'), {
    '/docs/installation/': '/site/start-here/installation/',
    '/docs/installation.md/': '/site/start-here/installation/',
    '/docs/README/': '/site/',
    '/docs/README.md/': '/site/',
  });
});

test('every redirect target in the checked-in manifest names a published route', async () => {
  const manifest = JSON.parse(fs.readFileSync(new URL('../redirects.json', import.meta.url), 'utf8'));
  const { listPublished, routeFor } = await import('../src/lib/published.mjs');
  const docsRoot = new URL('../../docs/', import.meta.url).pathname;
  const routes = new Set(listPublished(docsRoot).map(routeFor));
  for (const [name, target] of Object.entries(manifest)) {
    assert.ok(routes.has(target), `${name} -> ${target} is not a published route`);
  }
});
