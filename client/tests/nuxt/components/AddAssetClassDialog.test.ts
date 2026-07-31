import { mountSuspended } from "@nuxt/test-utils/runtime";
import { afterEach, describe, expect, it, vi } from "vitest";
import { nextTick } from "vue";

// @ts-ignore — oxlint's type-check can't resolve .vue modules
import AddAssetClassDialog from "~/components/AddAssetClassDialog.vue";
import useToast from "~/composables/useToast";

vi.setConfig({ testTimeout: 10_000 });

const NAME_SELECTOR = "#name";
const DESCRIPTION_SELECTOR = "#description";
const SUBMIT_SELECTOR = 'form button[type="submit"]';
const NO_TOASTS = 0;

const queryElement = <El extends Element>(
  selector: string,
  type: new () => El,
): El => {
  const element = document.querySelector(selector);
  if (element instanceof type) {
    return element;
  }

  throw new Error(`Expected selector "${selector}" to match a ${type.name}`);
};

describe("the AddAssetClassDialog component", () => {
  afterEach(() => {
    document.body.innerHTML = "";
    const { messages } = useToast();
    messages.value = [];
  });

  it("shows a validation error for an empty name and does not publish a toast", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AddAssetClassDialog);

    await wrapper.get("button").trigger("click");
    await nextTick();

    const submitButton = queryElement(
      SUBMIT_SELECTOR,
      globalThis.HTMLButtonElement,
    );
    submitButton.click();
    await nextTick();

    expect(document.body.textContent).toContain("Name is required");
    expect(useToast().messages.value).toHaveLength(NO_TOASTS);
  });

  it("publishes a toast with the added asset class name", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AddAssetClassDialog);

    await wrapper.get("button").trigger("click");
    await nextTick();

    const nameInput = queryElement(NAME_SELECTOR, globalThis.HTMLInputElement);
    nameInput.value = "Equities";
    nameInput.dispatchEvent(new Event("input"));
    await nextTick();

    const submitButton = queryElement(
      SUBMIT_SELECTOR,
      globalThis.HTMLButtonElement,
    );
    submitButton.click();
    await nextTick();

    expect(useToast().messages.value).toStrictEqual([
      'Asset class "Equities" created',
    ]);
  });

  it("resets the form fields after a successful submit", async () => {
    expect.hasAssertions();
    const wrapper = await mountSuspended(AddAssetClassDialog);

    await wrapper.get("button").trigger("click");
    await nextTick();

    const nameInput = queryElement(NAME_SELECTOR, globalThis.HTMLInputElement);
    nameInput.value = "Equities";
    nameInput.dispatchEvent(new Event("input"));
    const descriptionInput = queryElement(
      DESCRIPTION_SELECTOR,
      globalThis.HTMLTextAreaElement,
    );
    descriptionInput.value = "Stocks and equity securities";
    descriptionInput.dispatchEvent(new Event("input"));
    await nextTick();

    const submitButton = queryElement(
      SUBMIT_SELECTOR,
      globalThis.HTMLButtonElement,
    );
    submitButton.click();
    await nextTick();

    await wrapper.get("button").trigger("click");
    await nextTick();

    expect(queryElement(NAME_SELECTOR, globalThis.HTMLInputElement).value).toBe(
      "",
    );
    expect(
      queryElement(DESCRIPTION_SELECTOR, globalThis.HTMLTextAreaElement).value,
    ).toBe("");
  });
});
