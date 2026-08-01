import { describe, expect, it, vi } from "vitest";

import getFolioCode from "../../app/utils/asset-classes";

vi.setConfig({ testTimeout: 10_000 });

type FolioTestCase = Readonly<{
  expected: string;
  name: string;
  value: string;
}>;

const folioTestCases: readonly FolioTestCase[] = [
  {
    expected: "EQ",
    name: "single word",
    value: "Equities",
  },
  {
    expected: "FI",
    name: "two words",
    value: "Fixed income",
  },
  {
    expected: "RE",
    name: "three words uses first two initials",
    value: "Real estate investment trusts",
  },
  {
    expected: "GO",
    name: "leading and trailing whitespace",
    value: "  Gold  ",
  },
  {
    expected: "--",
    name: "empty string",
    value: "",
  },
  {
    expected: "--",
    name: "whitespace only",
    value: "   ",
  },
  {
    expected: "OI",
    name: "two letter word",
    value: "Oil",
  },
  {
    expected: "A",
    name: "single character word",
    value: "a",
  },
];

describe("getFolioCode", () => {
  it.each(folioTestCases)(
    "returns $expected for $name",
    (args: FolioTestCase) => {
      expect.hasAssertions();
      expect(getFolioCode(args.value)).toBe(args.expected);
    },
  );
});
