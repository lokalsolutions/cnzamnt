import { expect, test } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('loads the mobile app shell from the live API', async ({ page }) => {
  const failedResponses: string[] = [];
  page.on('response', (response) => {
    if (response.url().includes('/api/') && !response.ok()) {
      failedResponses.push(`${response.status()} ${response.url()}`);
    }
  });

  await page.goto('/');

  await expect(page.getByRole('heading', { name: 'CnzAMnt' })).toBeVisible();
  await expect(page.getByText('5,000 CNZ')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Post artwork' })).toBeVisible();
  await expect(page.getByText('No artwork yet.')).toBeVisible();
  expect(failedResponses).toEqual([]);
});

test('creates artwork and renders it in the vertical feed', async ({ page }) => {
  await page.goto('/');

  await page.getByRole('button', { name: 'Post artwork' }).click();
  await page.getByLabel('Title').fill('E2E Window Study');
  await page.getByLabel('Image URL').fill('data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw==');
  await page.getByLabel('Caption').fill('A compact test piece for the first art feed.');
  await page.getByRole('button', { name: 'Post', exact: true }).click();

  await expect(page.getByRole('heading', { name: 'E2E Window Study' })).toBeVisible();
  await expect(page.getByText('@demo_artist')).toBeVisible();
  await expect(page.getByText('A compact test piece for the first art feed.')).toBeVisible();
  await expect(page.getByText('0').first()).toBeVisible();
  await expect(page.getByText('comments').first()).toBeVisible();
});

test('shows validation without posting incomplete artwork', async ({ page }) => {
  await page.goto('/');

  await page.getByRole('button', { name: 'Post artwork' }).click();
  await page.getByLabel('Title').fill('Missing image');
  await page.getByRole('button', { name: 'Post', exact: true }).click();

  await expect(page.getByText('Title and image URL are required.')).toBeVisible();
});
