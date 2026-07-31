import { mountSuspended } from "@nuxt/test-utils/runtime";
import { afterEach, describe, expect, it, vi } from "vitest";
import { nextTick } from "vue";

// @ts-ignore — oxlint's type-check can't resolve .vue modules
import AppNavbar from "~/components/AppNavbar.vue";

vi.setConfig({ testTimeout: 10_000 });

describe("the AppNavbar component", () => {
  afterEach(() => {
    document.body.innerHTML = "";
  });

  it("renders the Overview link to the dashboard as the current page", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AppNavbar);
    const link = wrapper.get("a");

    expect(link.text()).toBe("Overview");
    expect(link.attributes("href")).toBe("/dashboard");
    expect(link.attributes("aria-current")).toBe("page");
  });

  it("renders the Add Asset Class trigger button", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AppNavbar);

    expect(wrapper.get("button").text()).toContain("Add Asset Class");
  });

  it("opens the asset class dialog when the trigger is clicked", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AppNavbar);

    await wrapper.get("button").trigger("click");
    await nextTick();

    const dialog = document.querySelector('[role="dialog"]');
    expect(dialog).not.toBeNull();
    expect(document.body.textContent).toContain(
      "Give your new asset class a name and description.",
    );
  });
});
