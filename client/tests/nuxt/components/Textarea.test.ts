import { mount } from "@vue/test-utils";
import { describe, expect, it, vi } from "vitest";

// @ts-ignore Oxlint cannot identify the Nuxt components
import Textarea from "~/components/Textarea.vue";

vi.setConfig({ testTimeout: 10_000 });

/** @description Number of rows forwarded to the textarea in the fallthrough test. */
const ROWS_ATTRIBUTE = 3;

/** @description Id used to associate the label with the textarea. */
const TEXTAREA_ID = "asset-description";

/** @description Label text shown above the textarea. */
const TEXTAREA_LABEL = "Description";

describe("textarea", () => {
  it("renders the label text", () => {
    expect.hasAssertions();

    const wrapper = mount(Textarea, {
      attrs: { id: TEXTAREA_ID },
      props: { label: TEXTAREA_LABEL },
    });

    expect(wrapper.get("label").text()).toBe(TEXTAREA_LABEL);
  });

  it("associates the label with the textarea via the id attribute", () => {
    expect.hasAssertions();

    const wrapper = mount(Textarea, {
      attrs: { id: TEXTAREA_ID },
      props: { label: TEXTAREA_LABEL },
    });

    expect(wrapper.get("label").attributes("for")).toBe(TEXTAREA_ID);
    expect(wrapper.get("textarea").attributes("id")).toBe(TEXTAREA_ID);
  });

  it("forwards placeholder and rows attributes to the native textarea", () => {
    expect.hasAssertions();

    const wrapper = mount(Textarea, {
      attrs: {
        id: TEXTAREA_ID,
        placeholder: "What belongs here?",
        rows: ROWS_ATTRIBUTE,
      },
      props: { label: TEXTAREA_LABEL },
    });

    const textarea = wrapper.get("textarea");
    expect(textarea.attributes("placeholder")).toBe("What belongs here?");
    expect(textarea.attributes("rows")).toBe(`${ROWS_ATTRIBUTE}`);
  });
});
