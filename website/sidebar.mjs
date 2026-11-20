import fs from 'node:fs';
import path from 'node:path';

const slugs = (...names) => names.map((slug) => ({ slug }));

const cliDirectory = new URL('../docs/reference/cli/', import.meta.url);
const cliCommands = fs
  .readdirSync(cliDirectory)
  .filter((name) => name.endsWith('.md') && name !== 'index.md')
  .map((name) => path.basename(name, '.md'))
  .sort();

export const sidebar = [
  {
    label: 'Start here',
    items: slugs(
      'start-here/overview',
      'start-here/installation',
      'start-here/first-clone',
      'start-here/first-exec',
      'start-here/guided-tour',
    ),
  },
  {
    label: 'Guides',
    items: [
      ...slugs(
        'guides',
        'guides/create-github-app',
        'guides/personal-access-tokens',
        'guides/git-credential-helper',
        'guides/multiple-organizations',
        'guides/github-enterprise-server',
        'guides/credential-lifecycle',
      ),
      {
        label: 'Continuous integration',
        collapsed: true,
        items: slugs('guides/ci', 'guides/ci/github-actions', 'guides/ci/jenkins', 'guides/ci/gitlab-ci'),
      },
    ],
  },
  {
    label: 'Reference',
    items: [
      { slug: 'reference' },
      {
        label: 'CLI commands',
        collapsed: true,
        items: [{ slug: 'reference/cli' }, ...cliCommands.map((name) => ({ slug: `reference/cli/${name}` }))],
      },
      ...slugs('reference/configuration', 'reference/environment-variables', 'reference/github-actions'),
    ],
  },
  {
    label: 'Understand',
    items: slugs(
      'understand/identity-and-credentials',
      'understand/credential-routing',
      'understand/tokens-and-caching',
      'understand/security-model',
      'understand/architecture',
      'understand/origins',
    ),
  },
  {
    label: 'Troubleshoot',
    items: slugs('troubleshoot', 'troubleshoot/diagnostics', 'troubleshoot/known-limitations'),
  },
  {
    label: 'Contribute',
    items: slugs(
      'contribute',
      'contribute/development',
      'contribute/testing',
      'contribute/e2e-testing',
      'contribute/e2e-infrastructure',
      'contribute/releasing',
      'contribute/documentation',
      'contribute/architecture-decisions',
    ),
  },
];
