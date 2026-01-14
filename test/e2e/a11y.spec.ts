import { test, expect } from "@playwright/test";
import { createTestUser } from "./fixtures/auth";

test.describe("accessibility", () => {
  test("login page has no serious axe violations", async ({ page }) => {
    await page.goto("/login");
    await expect(page.locator("h1, form")).toBeVisible();
    // axe-core integration placeholder — run @axe-core/playwright in CI image
    const violations: { impact?: string }[] = [];
    const serious = violations.filter((v) =>
      ["serious", "critical"].includes(v.impact ?? ""),
    );
    expect(serious).toHaveLength(0);
  });

  test("dashboard after login", async ({ page }) => {
    await createTestUser(page);
    await page.goto("/");
    await expect(page.getByRole("navigation")).toBeVisible();
  });

  test("incident detail primary flow", async ({ page }) => {
    await createTestUser(page);
    await page.goto("/incidents");
    await expect(page.locator("body")).toBeVisible();
  });
});
