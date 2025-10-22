import { test, expect } from "@playwright/test";
import { createTestUser } from "./fixtures/auth";
import { seedTopics } from "./fixtures/incidents";

test("create incident await ready", async ({ page }) => {
  await createTestUser(page);
  await seedTopics(page);
  await page.goto("/incidents");
  await page.click("text=Create incident");
  await expect(page.locator("dialog")).toBeVisible();
});
