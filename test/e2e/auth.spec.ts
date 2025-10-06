import { test, expect } from "@playwright/test";
import { createTestUser } from "./fixtures/auth";

test("register login logout", async ({ page }) => {
  await createTestUser(page);
  await expect(page).toHaveURL("/");
  await page.click("text=Logout");
  await page.goto("/login");
  await expect(page.locator("h1")).toContainText("Login");
});
