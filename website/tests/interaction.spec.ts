import { expect, test } from '@playwright/test';

test.describe('search', () => {
  test('keyboard search finds routing and the result opens, with back navigation', async ({ page }) => {
    await page.goto('');
    await page.keyboard.press('Control+k');
    const dialog = page.locator('dialog[open]');
    await expect(dialog).toBeVisible();
    const input = dialog.locator('.pagefind-ui__search-input');
    await expect(input).toBeAttached();
    await page.keyboard.press('Escape');
    await expect(dialog).toHaveCount(0);
    await page.keyboard.press('Control+k');
    await expect(input).toBeFocused();
    await input.fill('routing');
    const result = dialog.locator('.pagefind-ui__result-link').first();
    await expect(result).toBeVisible();
    await result.click();
    await expect(page).toHaveURL(/credential-routing\//);
    await expect(page.locator('dialog[open]')).toHaveCount(0);
    await page.goBack();
    await expect(page).toHaveURL(/\/gh-app-auth\/$/);
  });

  test('Escape closes search and returns focus to the search button', async ({ page }) => {
    await page.goto('');
    const button = page.locator('button[data-open-modal]');
    await expect(button).toBeEnabled();
    await button.click();
    await expect(page.locator('dialog[open]')).toBeVisible();
    await page.keyboard.press('Escape');
    await expect(page.locator('dialog[open]')).toHaveCount(0);
    await expect(button).toBeFocused();
  });

  test('search button has a name that includes its visible label', async ({ page }) => {
    await page.goto('');
    const button = page.getByRole('button', { name: 'Search' });
    await expect(button).toHaveCount(1);
    await expect(button).toHaveAttribute('aria-keyshortcuts', /Control\+K|Meta\+K/);
  });
});

test.describe('search on mobile', () => {
  test.use({ viewport: { width: 390, height: 844 } });

  test('icon button keeps its accessible name', async ({ page }) => {
    await page.goto('');
    await expect(page.getByRole('button', { name: 'Search' })).toBeVisible();
  });
});

test('copy button copies the code block', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await page.goto('start-here/first-clone/');
  const block = page.locator('.expressive-code').first();
  await block.hover();
  await block.locator('.copy button').click();
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe(
    'gh extension install AmadeusITGroup/gh-app-auth\n'.trimEnd(),
  );
});

test('theme choice persists across navigation', async ({ page }) => {
  await page.goto('');
  const select = page.locator('starlight-theme-select select:visible');
  await select.selectOption('dark');
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
  await page.getByRole('link', { name: 'Make your first clone' }).click();
  await expect(page).toHaveURL(/start-here\/first-clone\/$/);
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
  await page.locator('#starlight__sidebar a[href$="/understand/credential-routing/"]').click();
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark');
  await select.selectOption('light');
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light');
});
