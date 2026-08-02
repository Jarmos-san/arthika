import { describe, expect, it, vi } from "vitest";

import { validateAssetClass, validateLogin } from "../../app/utils/validators";

vi.setConfig({ testTimeout: 10_000 });

const emailRequired = { email: "Email is required", password: undefined };
const passwordRequired = { email: undefined, password: "Password is required" };
const bothRequired = {
  email: "Email is required",
  password: "Password is required",
};

type TestCase = Readonly<{
  email: string | undefined;
  expected: Readonly<{
    email: string | undefined;
    password: string | undefined;
  }>;
  name: string;
  password: string | undefined;
}>;

const testCases: readonly TestCase[] = [
  {
    email: "user@example.com",
    expected: { email: undefined, password: undefined },
    name: "valid email and password",
    password: "password123",
  },
  {
    email: undefined,
    expected: emailRequired,
    name: "undefined email",
    password: "password123",
  },
  {
    email: "not-an-email",
    expected: {
      email: "Please enter a valid email address",
      password: undefined,
    },
    name: "invalid email format",
    password: "password123",
  },
  {
    email: "user@example.com",
    expected: passwordRequired,
    name: "undefined password",
    password: undefined,
  },
  {
    email: "user@example.com",
    expected: {
      email: undefined,
      password: "Password must be at least 8 characters",
    },
    name: "password shorter than 8 characters",
    password: "short",
  },
  {
    email: undefined,
    expected: bothRequired,
    name: "both fields undefined",
    password: undefined,
  },
];

describe("validateLogin", () => {
  it.each(testCases)("returns correct errors for $name", (args: TestCase) => {
    expect.hasAssertions();
    expect(validateLogin(args.email, args.password)).toStrictEqual(
      args.expected,
    );
  });
});

type AssetClassErrorTestCase = Readonly<{
  expected: string | undefined;
  name: string;
  value: string | undefined;
}>;

const assetClassErrorTestCases: readonly AssetClassErrorTestCase[] = [
  {
    expected: "Name is required",
    name: "undefined name",
    value: undefined,
  },
  {
    expected: "Name is required",
    name: "empty name",
    value: "",
  },
  {
    expected: "Name is required",
    name: "whitespace only name",
    value: "   ",
  },
  {
    expected: undefined,
    name: "valid name",
    value: "Equities",
  },
];

describe("validateAssetClass", () => {
  it.each(assetClassErrorTestCases)(
    "returns correct errors for $name",
    (args: AssetClassErrorTestCase) => {
      expect.hasAssertions();
      expect(validateAssetClass(args.value)).toStrictEqual({
        name: args.expected,
      });
    },
  );
});
