import assert from 'node:assert/strict';
import test from 'node:test';
import { extractLinks, frontmatterOf, headingAnchors } from '../scripts/lib/markdown.mjs';

const sample = `---
title: T
description: D
---

# Top \`code\` heading

See [a](./a.md#x) and \`[no](./skip.md)\`.

\`\`\`md
[ignored](./nope.md)
# not a heading
\`\`\`

## Top \`code\` heading

[ref]: ./ref.md
`;

test('extractLinks ignores code and reports reference definitions', () => {
  assert.deepEqual(
    extractLinks(sample).map((link) => link.href),
    ['./a.md#x', './ref.md'],
  );
});

test('headingAnchors uses github slugs with duplicate suffixes and skips fences', () => {
  assert.deepEqual([...headingAnchors(sample)], ['top-code-heading', 'top-code-heading-1']);
});

test('frontmatterOf reads simple fields', () => {
  assert.deepEqual(frontmatterOf(sample), { title: 'T', description: 'D' });
  assert.deepEqual(frontmatterOf('# no frontmatter'), {});
});
