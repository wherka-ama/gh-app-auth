export const relatedLinks = {
  completion: [['Installation', '../../start-here/installation.md']],
  config: [
    ['Configuration reference', '../configuration.md'],
    ['Environment variables', '../environment-variables.md'],
  ],
  exec: [
    ['Run your first command with exec', '../../start-here/first-exec.md'],
    ['Security model', '../../understand/security-model.md'],
  ],
  gitconfig: [['Git credential helper guide', '../../guides/git-credential-helper.md']],
  list: [['Credential lifecycle', '../../guides/credential-lifecycle.md']],
  migrate: [['Credential lifecycle', '../../guides/credential-lifecycle.md']],
  remove: [['Credential lifecycle', '../../guides/credential-lifecycle.md']],
  scope: [['Credential routing', '../../understand/credential-routing.mdx']],
  setup: [
    ['Your first authenticated clone', '../../start-here/first-clone.md'],
    ['Personal access tokens', '../../guides/personal-access-tokens.md'],
    ['Known limitations', '../../troubleshoot/known-limitations.md'],
  ],
  test: [['Troubleshooting', '../../troubleshoot/index.md']],
  token: [['Security model', '../../understand/security-model.md']],
};

export function parseCommands(helpText) {
  const lines = helpText.split('\n');
  const start = lines.findIndex((line) => line.trim() === 'Available Commands:');
  if (start === -1) return [];
  const commands = [];
  for (const line of lines.slice(start + 1)) {
    if (line.trim() === '') break;
    const match = line.match(/^ {2}(\S+)\s+(.*)$/);
    if (match && match[1] !== 'help') commands.push({ name: match[1], summary: match[2].trim() });
  }
  return commands;
}

export function normalizeHelp(text) {
  return text
    .replace(/\r\n/g, '\n')
    .replace(/\t/g, '  ')
    .split('\n')
    .map((line) => line.replace(/\s+$/, ''))
    .join('\n')
    .replace(/\n{3,}/g, '\n\n')
    .trim();
}

export function renderPage({ name, summary, sections }) {
  const related = relatedLinks[name] ?? [];
  const body = sections
    .map(({ title, help }, index) => {
      const heading = index === 0 ? '## Help output' : `## ${title}`;
      return `${heading}\n\n\`\`\`text\n${help}\n\`\`\``;
    })
    .join('\n\n');
  const relatedBlock = related.length
    ? `\n\n## Related\n\n${related.map(([label, target]) => `- [${label}](${target})`).join('\n')}`
    : '';
  return `---
title: gh app-auth ${name}
description: ${summary.replace(/:/g, ' -')}
---

This page is generated from the output of \`gh-app-auth ${name} --help\` in the checked-out source. To refresh it, run \`npm run reference:generate\` in \`website/\`.

${body}${relatedBlock}
`;
}
