import { beforeEach, describe, expect, it, vi } from "vitest";

import getRedirectPath from "../../app/utils/routing";

vi.setConfig({ testTimeout: 10_000 });

type TestCase = Readonly<{
  expected: string | undefined;
  isAuthenticated: boolean;
  targetPath: string;
}>;

const testCases: readonly TestCase[] = [
  { expected: "/dashboard", isAuthenticated: true, targetPath: "/login" },
  { expected: undefined, isAuthenticated: true, targetPath: "/dashboard" },
  { expected: "/login", isAuthenticated: false, targetPath: "/dashboard" },
  { expected: undefined, isAuthenticated: false, targetPath: "/login" },
  { expected: undefined, isAuthenticated: false, targetPath: "/register" },
  { expected: undefined, isAuthenticated: true, targetPath: "/register" },
  { expected: "/dashboard", isAuthenticated: true, targetPath: "/" },
  { expected: "/dashboard", isAuthenticated: false, targetPath: "/" },
];

describe("getRedirectPath", () => {
  beforeEach(() => {
    vi.stubEnv("NODE_ENV", "production");
  });

  it.each(testCases)(
    "returns $expected for $targetPath when isAuthenticated is $isAuthenticated",
    ({ expected, isAuthenticated, targetPath }: TestCase) => {
      expect.hasAssertions();
      expect(getRedirectPath(isAuthenticated, targetPath)).toBe(expected);
    },
  );
});
