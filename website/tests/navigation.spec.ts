import { expect, test } from '@playwright/test';
import fs from 'node:fs';
import { pages, primaryRoutes } from './pages';

const groups = ['Start here', 'Guides', 'Reference', 'Understand', 'Troubleshoot', 'Contribute'];
const redirects: Record<string, string> = JSON.parse(
  fs.readFileSync(new URL('../redirects.json', import.meta.url), 'utf8'),
);

test.describe('sidebar and landmarks', () => {
  test('sidebar shows every group and visible links under the open groups', async ({ page }) => {
    await page.goto('start-here/first-clone/');
    const sidebar = page.locator('#starlight__sidebar');
    for (const group of groups) {
      await expect(sidebar.locator('summary', { hasText: group }).first()).toBeVisible();
    }
    const links = sidebar.locator('a[href^="/gh-app-auth/"]:visible');
    expect(await links.count()).toBeGreaterThan(20);
    await expect(sidebar.locator('a[aria-current="page"]')).toHaveCount(1);
    await expect(sidebar.locator('a[aria-current="page"]')).toHaveText('Make your first authenticated clone');
  });

  test('every primary route has a sidebar link', async ({ page }) => {
    await page.goto('start-here/overview/');
    const sidebar = page.locator('#starlight__sidebar');
    for (const route of primaryRoutes) {
      await expect(sidebar.locator(`a[href$="/${route}"]`), route).toHaveCount(1);
    }
  });

  for (const [group, child] of [
    ['CLI commands', 'reference/cli/setup/'],
    ['Continuous integration', 'guides/ci/jenkins/'],
  ] as const) {
    test(`nested group "${group}" expands and its pages open`, async ({ page }) => {
      await page.goto('start-here/overview/');
      const sidebar = page.locator('#starlight__sidebar');
      const link = sidebar.locator(`a[href$="/${child}"]`);
      await expect(link).toBeHidden();
      await sidebar.locator('summary', { hasText: group }).click();
      await expect(link).toBeVisible();
      await link.click();
      await expect(page).toHaveURL(new RegExp(`${child}$`));
      await expect(page.locator('h1')).toBeVisible();
      await expect(sidebar.locator(`a[aria-current="page"][href$="/${child}"]`)).toBeVisible();
    });
  }

  for (const route of ['understand/credential-routing/', 'troubleshoot/known-limitations/', 'contribute/releasing/']) {
    test(`clicking the sidebar link for ${route} navigates`, async ({ page }) => {
      await page.goto('start-here/first-clone/');
      await page.locator(`#starlight__sidebar a[href$="/${route}"]`).click();
      await expect(page).toHaveURL(new RegExp(`${route}$`));
      await expect(page.locator('h1')).toBeVisible();
    });
  }

  test('home page exposes a named route-choice navigation landmark', async ({ page }) => {
    await page.goto('');
    const nav = page.getByRole('navigation', { name: 'Choose where to start' });
    await expect(nav.getByRole('link')).toHaveCount(3);
  });

  for (const path of pages) {
    test(`landmarks and single h1 on ${path || 'home'}`, async ({ page }) => {
      await page.goto(path);
      await expect(page.locator('h1')).toHaveCount(1);
      await expect(page.locator('main')).toHaveCount(1);
      await expect(page.locator('header').first()).toBeVisible();
      await expect(page.getByRole('navigation').first()).toBeAttached();
    });
  }

  test('favicon loads from the site base', async ({ page, request }) => {
    await page.goto('');
    const href = await page.locator('link[rel="shortcut icon"]').first().getAttribute('href');
    expect(href).toContain('/gh-app-auth/');
    const response = await request.get(new URL(href!, page.url()).toString());
    expect(response.status()).toBe(200);
  });

  test('unknown pages render the custom 404', async ({ page }) => {
    const response = await page.goto('does-not-exist/');
    expect(response?.status()).toBe(404);
    await expect(page.locator('h1')).toHaveText('Page not found');
    await page.getByRole('link', { name: 'Your first authenticated clone' }).click();
    await expect(page).toHaveURL(/start-here\/first-clone\/$/);
  });

  test('all internal links and anchors resolve', async ({ page, request }) => {
    const targets = new Map<string, Set<string>>();
    for (const path of pages) {
      await page.goto(path);
      const hrefs = await page.locator('a[href]').evaluateAll((as) => as.map((a) => (a as HTMLAnchorElement).href));
      for (const href of hrefs) {
        const url = new URL(href);
        if (url.origin !== new URL(page.url()).origin) continue;
        const key = url.origin + url.pathname;
        if (!targets.has(key)) targets.set(key, new Set());
        if (url.hash.length > 1) targets.get(key)!.add(decodeURIComponent(url.hash.slice(1)));
      }
    }
    expect(targets.size).toBeGreaterThan(pages.length / 2);
    for (const [url, anchors] of targets) {
      expect(url, url).toContain('/gh-app-auth/');
      const response = await request.get(url);
      expect(response.status(), url).toBe(200);
      const html = await response.text();
      for (const anchor of anchors) {
        expect(html.includes(`id="${anchor}"`), `${url}#${anchor}`).toBe(true);
      }
    }
  });

  test('every published page is linked from somewhere reachable', async ({ page }) => {
    const reached = new Set<string>();
    await page.goto('start-here/overview/');
    const sidebarHrefs = await page
      .locator('#starlight__sidebar a[href^="/gh-app-auth/"]')
      .evaluateAll((as) => as.map((a) => (a as HTMLAnchorElement).getAttribute('href')!));
    sidebarHrefs.forEach((h) => reached.add(h));
    for (const path of pages.filter((p) => p !== '')) {
      expect(reached.has(`/gh-app-auth/${path}`), path).toBe(true);
    }
  });

  test('heading anchors link to their headings', async ({ page }) => {
    for (const path of ['start-here/first-clone/', 'understand/architecture/']) {
      await page.goto(path);
      const headings = page.locator('.sl-markdown-content h2[id]');
      expect(await headings.count()).toBeGreaterThan(3);
      const first = headings.nth(1);
      const id = await first.getAttribute('id');
      const anchor = first.locator('xpath=..').locator(`a.sl-anchor-link[href="#${id}"]`);
      await expect(anchor).toHaveCount(1);
      await anchor.click({ force: true });
      await expect(page).toHaveURL(new RegExp(`#${id}$`));
    }
  });

  for (const [name, target] of Object.entries(redirects)) {
    for (const suffix of ['', '.md']) {
      test(`legacy route docs/${name}${suffix} redirects`, async ({ page }) => {
        await page.goto(`docs/${name}${suffix}/`);
        await page.waitForURL(new RegExp(`/gh-app-auth/${target}$`));
        await expect(page.locator('h1')).toHaveCount(1);
      });
    }
  }
});

test.describe('mobile menu', () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test('opens, navigates and closes', async ({ page }) => {
    await page.goto('start-here/first-clone/');
    const toggle = page.locator('button[popovertarget="starlight__sidebar"]');
    const sidebar = page.locator('#starlight__sidebar');
    await expect(sidebar).toBeHidden();
    await toggle.click();
    await expect(sidebar).toBeVisible();
    await sidebar.locator('a[href$="/understand/credential-routing/"]').click();
    await expect(page).toHaveURL(/understand\/credential-routing\/$/);
    await expect(page.locator('h1')).toHaveText('How credential routing works');
    await expect(sidebar).toBeHidden();
    await expect(toggle).toBeVisible();
  });
});
