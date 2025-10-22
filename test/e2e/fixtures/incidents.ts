export async function seedTopics(page: import("@playwright/test").Page) {
  await page.evaluate(() => {
    localStorage.setItem("replay_project_id", "00000000-0000-0000-0000-000000000001");
  });
}
