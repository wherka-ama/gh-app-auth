import fs from 'node:fs';

export function loadRedirects(file, base) {
  const legacy = JSON.parse(fs.readFileSync(file, 'utf8'));
  const redirects = {};
  for (const [name, target] of Object.entries(legacy)) {
    for (const source of [`/docs/${name}/`, `/docs/${name}.md/`]) {
      redirects[source] = `${base}${target}`;
    }
  }
  return redirects;
}
