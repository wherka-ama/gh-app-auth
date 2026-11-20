import GithubSlugger from 'github-slugger';

function withoutFences(text) {
  const lines = text.split('\n');
  const kept = [];
  let fence = null;
  lines.forEach((line, index) => {
    const marker = line.match(/^\s*(`{3,}|~{3,})/);
    if (marker) {
      if (fence === null) fence = marker[1][0];
      else if (marker[1][0] === fence) fence = null;
      kept.push({ line: '', number: index + 1 });
      return;
    }
    kept.push({ line: fence === null ? line : '', number: index + 1 });
  });
  return kept;
}

function stripFrontmatter(text) {
  if (!text.startsWith('---\n')) return text;
  const end = text.indexOf('\n---', 4);
  if (end === -1) return text;
  const removed = text.slice(0, end + 4).split('\n').length - 1;
  return '\n'.repeat(removed) + text.slice(end + 4);
}

export function extractLinks(text) {
  const links = [];
  for (const { line, number } of withoutFences(stripFrontmatter(text))) {
    const clean = line.replace(/`[^`]*`/g, '');
    for (const match of clean.matchAll(/!?\[[^\]]*\]\(\s*<?([^)\s>]+)>?(?:\s+"[^"]*")?\s*\)/g)) {
      links.push({ href: match[1], line: number });
    }
    const definition = clean.match(/^\[[^\]]+\]:\s*(\S+)/);
    if (definition) links.push({ href: definition[1], line: number });
  }
  return links;
}

function headingText(raw) {
  return raw
    .replace(/\s+#+\s*$/, '')
    .replace(/!?\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/[`*_]/g, '')
    .replace(/<[^>]+>/g, '')
    .trim();
}

export function headingAnchors(text) {
  const slugger = new GithubSlugger();
  const anchors = new Set();
  for (const { line } of withoutFences(stripFrontmatter(text))) {
    const heading = line.match(/^#{1,6}\s+(.*)$/);
    if (heading) anchors.add(slugger.slug(headingText(heading[1])));
    for (const explicit of line.matchAll(/<a\s+(?:name|id)="([^"]+)"/g)) anchors.add(explicit[1]);
  }
  return anchors;
}

export function frontmatterOf(text) {
  if (!text.startsWith('---\n')) return {};
  const end = text.indexOf('\n---', 4);
  if (end === -1) return {};
  const fields = {};
  for (const line of text.slice(4, end).split('\n')) {
    const match = line.match(/^([A-Za-z]+):\s*(.*)$/);
    if (match) fields[match[1]] = match[2];
  }
  return fields;
}
