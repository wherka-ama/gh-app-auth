import { test } from '@playwright/test';
import { slugOf } from './pages';

const outDir = process.env.DOCS_SCREENSHOT_DIR;
const targets = [
  '',
  'start-here/first-clone/',
  'start-here/guided-tour/',
  'understand/credential-routing/',
  'understand/architecture/',
  'reference/configuration/',
  'troubleshoot/',
  'contribute/testing/',
  'contribute/e2e-testing/',
  'guides/ci/gitlab-ci/',
  'guides/credential-lifecycle/',
];
const viewports = [
  { name: 'desktop', width: 1440, height: 1000 },
  { name: 'mobile', width: 390, height: 844 },
];

test.skip(!outDir, 'set DOCS_SCREENSHOT_DIR to capture screenshots');

for (const viewport of viewports) {
  for (const theme of ['light', 'dark']) {
    for (const path of targets) {
      test(`${viewport.name} ${theme} ${slugOf(path)}`, async ({ page }) => {
        await page.setViewportSize({ width: viewport.width, height: viewport.height });
        await page.goto(path);
        await page.evaluate((value) => localStorage.setItem('starlight-theme', value), theme);
        await page.reload();
        const base = `${outDir}/${viewport.name}-${theme}-${slugOf(path)}`;
        await page.screenshot({ path: `${base}-viewport.png` });
        await page.screenshot({ path: `${base}-full.png`, fullPage: true });
      });
    }
  }
}
