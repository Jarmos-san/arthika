<script setup lang="ts">
  import {
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogOverlay,
    DialogPortal,
    DialogRoot,
    DialogTitle,
    DialogTrigger,
    Label,
  } from "reka-ui";
  import { ref } from "vue";

  import useToast from "~/composables/useToast";
  import type { AssetClassFormErrors } from "~/types/utils/validators";
  import { validateAssetClass } from "~/utils/validators";

  // Initialise the state of the form modal to add an asset class
  const dialogOpen = ref(false);
  const name = ref<string | undefined>(undefined);
  const description = ref<string | undefined>(undefined);
  const errors = ref<AssetClassFormErrors>({ name: undefined });

  // Create an instance of the "toast" component from Reka UI
  const toast = useToast();

  // Close the dialog and reset the form after a successful submission
  const submitAssetClass = (): void => {
    const addedName = name.value;
    dialogOpen.value = false;
    name.value = undefined;
    description.value = undefined;
    errors.value = { name: undefined };
    toast.publish(`Asset class "${addedName}" created`);
  };

  // Event handler for the submit button
  const onSubmit = (): void => {
    errors.value = validateAssetClass(name.value);

    if (errors.value.name) {
      return;
    }

    submitAssetClass();
  };
</script>

<template>
  <DialogRoot v-model:open="dialogOpen">
    <DialogTrigger
      class="flex cursor-pointer items-center gap-1.5 rounded-lg bg-amber-600 px-3.5 py-2.5 text-sm leading-5 font-medium text-stone-950 transition-colors hover:bg-amber-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-stone-950 active:scale-[0.98] sm:px-4"
    >
      + Add Asset Class
    </DialogTrigger>
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-40 bg-stone-900/90" />
      <DialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl border border-stone-200 bg-white p-6 shadow-xl"
      >
        <DialogTitle class="text-lg font-semibold text-stone-800">
          Add Asset Class
        </DialogTitle>
        <DialogDescription class="mt-1 text-sm text-stone-500">
          Give your new asset class a name and description.
        </DialogDescription>
        <form class="mt-5 flex flex-col gap-5" @submit.prevent="onSubmit">
          <div>
            <Input
              id="name"
              v-model="name"
              aria-describedby="name-error"
              autofocus
              label="Name"
              placeholder="e.g. Equities"
            />
            <p
              v-if="errors.name"
              id="name-error"
              class="mt-1.5 text-xs text-red-600"
            >
              {{ errors.name }}
            </p>
          </div>
          <div>
            <Label for="description" class="text-sm font-medium text-stone-600">
              Description
            </Label>
            <textarea
              id="description"
              v-model="description"
              rows="3"
              placeholder="e.g. Stocks and equity securities"
              class="mt-1.5 w-full rounded-lg border border-stone-300 bg-stone-50 px-3.5 py-2.5 text-sm text-stone-800 placeholder:text-stone-400"
            />
          </div>
          <div class="flex justify-end gap-3">
            <DialogClose
              class="cursor-pointer rounded-lg border border-stone-300 px-4 py-2.5 text-sm leading-5 font-medium text-stone-700 transition-colors hover:border-stone-400 hover:bg-stone-100 hover:text-stone-900 active:scale-[0.98]"
            >
              Cancel
            </DialogClose>
            <button class="btn-primary" type="submit"> Add Asset Class </button>
          </div>
        </form>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
