import { test, expect } from "@playwright/test";
import { createTestUser } from "./fixtures/auth";
import { mockReplayRoutes } from "./fixtures/replay";

test("start replay mock terminal succeeded", async ({ page }) => {
  await mockReplayRoutes(page);
  await createTestUser(page);
  await page.evaluate(() =>
    localStorage.setItem("replay_project_id", "proj-1"),
  );
  await page.goto("/incidents/inc-mock");
  await expect(page.locator("h1")).toBeVisible();
});
