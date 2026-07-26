import { describe, expect, it, vi } from "vitest";

import getRedirectPath from "../../app/utils/routing";

vi.setConfig({ testTimeout: 10_000 });

describe("getRedirectPath", () => {
  it("redirects authenticated user from /login to /dashboard", () => {
    expect.hasAssertions();
    expect(getRedirectPath(true, "/login")).toBe("/dashboard");
  });

  it("returns undefined for authenticated user on non-login path", () => {
    expect.hasAssertions();
    expect(getRedirectPath(true, "/dashboard")).toBeUndefined();
  });

  it("redirects unauthenticated user from protected path to /login", () => {
    expect.hasAssertions();
    expect(getRedirectPath(false, "/dashboard")).toBe("/login");
  });

  it("returns undefined for unauthenticated user on /login", () => {
    expect.hasAssertions();
    expect(getRedirectPath(false, "/login")).toBeUndefined();
  });

  it("returns undefined for unauthenticated user on /register", () => {
    expect.hasAssertions();
    expect(getRedirectPath(false, "/register")).toBeUndefined();
  });

  it("returns undefined for authenticated user on /register", () => {
    expect.hasAssertions();
    expect(getRedirectPath(true, "/register")).toBeUndefined();
  });
});
