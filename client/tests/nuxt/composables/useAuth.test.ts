import { registerEndpoint } from "@nuxt/test-utils/runtime";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ref } from "vue";

import { useState } from "#app";
import useAuth, {
  STATUS_CONFLICT,
  STATUS_NETWORK_ERROR,
  STATUS_OK,
  STATUS_UNAUTHORIZED,
  STATUS_UNPROCESSABLE,
} from "~/composables/useAuth";

vi.setConfig({ testTimeout: 10_000 });

vi.mock(import("#app"), async (importOriginal) => {
  const actual = await importOriginal();

  // oxlint-disable-next-line eslint/prefer-object-spread
  return Object.assign({}, actual, {
    useAsyncData: (): { data: unknown; pending: boolean } => ({
      data: ref<unknown>(undefined),
      pending: false,
    }),
  });
});

describe("useAuth", () => {
  let endpointCleanups: (() => void)[];

  beforeEach(() => {
    endpointCleanups = [];
  });

  const rememberEndpoint = (
    url: string,
    // oxlint-disable-next-line typescript/prefer-readonly-parameter-types, eslint/no-magic-numbers
    options: Parameters<typeof registerEndpoint>[1],
  ): void => {
    endpointCleanups.push(registerEndpoint(url, options));
  };

  afterEach(() => {
    for (const cleanup of endpointCleanups) {
      cleanup();
    }

    endpointCleanups = [];
    useState("auth-user").value = undefined;
  });

  it("exports status constants", () => {
    expect.hasAssertions();

    const STATUS_OK_CODE = 200;
    const STATUS_UNAUTHORIZED_CODE = 401;
    const STATUS_CONFLICT_CODE = 409;
    const STATUS_UNPROCESSABLE_CODE = 422;
    const STATUS_NETWORK_ERROR_CODE = 0;

    expect(STATUS_OK).toBe(STATUS_OK_CODE);
    expect(STATUS_UNAUTHORIZED).toBe(STATUS_UNAUTHORIZED_CODE);
    expect(STATUS_CONFLICT).toBe(STATUS_CONFLICT_CODE);
    expect(STATUS_UNPROCESSABLE).toBe(STATUS_UNPROCESSABLE_CODE);
    expect(STATUS_NETWORK_ERROR).toBe(STATUS_NETWORK_ERROR_CODE);
  });

  it("is not authenticated by default", () => {
    expect.hasAssertions();
    const { isAuthenticated } = useAuth();

    expect(isAuthenticated.value).toBeFalsy();
  });

  it("is authenticated when state is pre-populated", () => {
    expect.hasAssertions();
    useState("auth-user").value = {
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    };

    const { isAuthenticated } = useAuth();

    expect(isAuthenticated.value).toBeTruthy();
  });

  it("login succeeds and updates state", async () => {
    expect.hasAssertions();
    rememberEndpoint("/api/users/login", {
      handler: () => ({
        email: "user@example.com",
        id: "550e8400-e29b-41d4-a716-446655440000",
      }),
      method: "POST",
    });

    const { login, isAuthenticated } = useAuth();
    const result = await login("user@example.com", "password");

    expect(result).toStrictEqual({ ok: true, status: STATUS_OK });
    expect(isAuthenticated.value).toBeTruthy();
  });

  it("login failure returns 401 with error body", async () => {
    expect.hasAssertions();
    rememberEndpoint("/api/users/login", {
      handler: (): Response =>
        Response.json(
          { message: "Invalid email or password" },
          { status: 401 },
        ),
      method: "POST",
    });

    const { login, isAuthenticated } = useAuth();
    const result = await login("user@example.com", "wrong");

    expect(result).toStrictEqual({
      body: { message: "Invalid email or password" },
      ok: false,
      status: STATUS_UNAUTHORIZED,
    });
    expect(isAuthenticated.value).toBeFalsy();
  });

  it("register succeeds and updates state", async () => {
    expect.hasAssertions();
    rememberEndpoint("/api/users/register", {
      handler: () => ({
        email: "newuser@example.com",
        id: "550e8400-e29b-41d4-a716-446655440000",
      }),
      method: "POST",
    });

    const { register, isAuthenticated } = useAuth();
    const result = await register("newuser@example.com", "password");

    expect(result).toStrictEqual({ ok: true, status: 201 });
    expect(isAuthenticated.value).toBeTruthy();
  });

  it("register failure returns error body", async () => {
    expect.hasAssertions();
    rememberEndpoint("/api/users/register", {
      handler: (): Response =>
        Response.json({ message: "Email already taken" }, { status: 409 }),
      method: "POST",
    });

    const { register, isAuthenticated } = useAuth();
    const result = await register("existing@example.com", "password");

    expect(result).toStrictEqual({
      body: { message: "Email already taken" },
      ok: false,
      status: STATUS_CONFLICT,
    });
    expect(isAuthenticated.value).toBeFalsy();
  });
});
