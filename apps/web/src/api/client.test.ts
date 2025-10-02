import { describe, expect, it } from "vitest";
import { ApiError } from "./client";

describe("ApiError", () => {
  it("carries status and code", () => {
    const err = new ApiError(404, "not_found", "missing");
    expect(err.status).toBe(404);
    expect(err.code).toBe("not_found");
  });
});
