import type { Page } from "@playwright/test";

export async function createTestUser(page: Page) {
  const email = `user-${Date.now()}@example.com`;
  await page.goto("/register");
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', "password-12chars");
  await page.fill('input:not([type="email"]):not([type="password"])', "Test Org");
  await page.click('button[type="submit"]');
  return { email };
}
