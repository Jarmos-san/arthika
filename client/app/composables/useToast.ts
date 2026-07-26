import { ref } from "vue";

/** @description Shared list of active toast messages. */
const messages = ref<string[]>([]);

/** @description Number of items to remove when dismissing a toast. */
const DISMISS_COUNT = 1;

/** @description Return type for the `useToast` composable. */
interface UseToast {
  /** @description Removes the toast at the given index from the list. */
  dismiss: (index: number) => void;
  /** @description Reactive list of current toast messages. */
  messages: typeof messages;
  /** @description Adds a new toast message to the viewport. */
  publish: (message: string) => void;
}

/**
 * @description Composable for triggering toast notifications from any component. Uses a
 * shared message list so toasts triggered anywhere in the app are rendered by
 * the `ToastProvider` in `app.vue`.
 *
 * @returns {UseToast} An object with `publish`, `dismiss`, and the reactive
 *   `messages` list.
 */
const useToast = (): UseToast => {
  const dismiss = (index: number): void => {
    messages.value.splice(index, DISMISS_COUNT);
  };

  const publish = (message: string): void => {
    messages.value.push(message);
  };

  return { dismiss, messages, publish };
};

export default useToast;
