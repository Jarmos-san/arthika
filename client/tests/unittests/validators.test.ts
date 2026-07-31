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

const nameRequired = { name: "Name is required" };
const nameValid = { name: undefined };

const LONG_NAME =
  "Asset class name that exceeds typical length limits for testing purposes";

type AssetClassTestCase = Readonly<{
  expected: Readonly<{ name: string | undefined }>;
  name: string;
  value: string | undefined;
}>;

const assetClassTestCases: readonly AssetClassTestCase[] = [
  {
    expected: nameValid,
    name: "valid name",
    value: "Equities",
  },
  {
    expected: nameRequired,
    name: "undefined name",
    value: undefined,
  },
  {
    expected: nameRequired,
    name: "empty name",
    value: "",
  },
  {
    expected: nameRequired,
    name: "whitespace-only name",
    value: "   ",
  },
  {
    expected: nameValid,
    name: "name with surrounding whitespace",
    value: "  Equities  ",
  },
  {
    expected: nameValid,
    name: "long name",
    value: LONG_NAME,
  },
];

describe("validateAssetClass", () => {
  it.each(assetClassTestCases)(
    "returns correct errors for $name",
    (args: AssetClassTestCase) => {
      expect.hasAssertions();
      expect(validateAssetClass(args.value)).toStrictEqual(args.expected);
    },
  );
});
