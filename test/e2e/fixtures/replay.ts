import type { Page } from "@playwright/test";

export async function mockReplayRoutes(page: Page) {
  await page.route("**/v1/incidents/*/replays", async (route) => {
    if (route.request().method() === "POST") {
      await route.fulfill({
        status: 201,
        body: JSON.stringify({ id: "replay-mock-1", status: "pending" }),
      });
      return;
    }
    await route.continue();
  });
  await page.route("**/v1/replays/*", async (route) => {
    await route.fulfill({
      status: 200,
      body: JSON.stringify({
        replay: { id: "replay-mock-1", status: "succeeded" },
      }),
    });
  });
}
