<script setup lang="ts">
  import DeleteIcon from "@iconify-vue/material-symbols/delete";
  import EditOutlineRoundedIcon from "@iconify-vue/material-symbols/edit-outline-rounded";
  import { ref } from "vue";

  import type { AssetClass } from "~/openapi";

  interface Props {
    /** @description Asset classes to render as ledger rows. */
    assetClasses: AssetClass[] | undefined;
  }

  const props = defineProps<Props>();

  /** @description The list of header elements to render on top of the table. */
  const headers = ref(["CLASS", "DESCRIPTION", "ACTIONS"]);
</script>

<template>
  <div class="rounded-xl border border-stone-200/60 bg-white shadow-lg">
    <div class="overflow-x-auto">
      <table
        v-if="props.assetClasses ? props.assetClasses.length > 0 : undefined"
        class="w-full text-sm"
      >
        <!-- Table header -->
        <thead>
          <tr class="border-b border-stone-200 text-center">
            <th
              v-for="(header, index) in headers"
              :key="index"
              class="px-6 py-3.5 font-mono text-xs font-medium tracking-wider text-stone-400"
            >
              {{ header }}
            </th>
          </tr>
        </thead>

        <!-- Table body -->
        <tbody class="divide-y divide-stone-200">
          <tr
            v-for="assetClass in props.assetClasses"
            :key="assetClass.id"
            class="transition-colors duration-150 hover:bg-stone-50/60"
          >
            <!-- Name of the class -->
            <td class="px-6 py-4 font-medium text-stone-800">
              {{ assetClass.name }}
            </td>

            <!-- Description of the class -->
            <td class="px-6 py-4 text-stone-500">
              {{ assetClass.description }}
            </td>

            <!-- Action buttons -->
            <td class="p-4">
              <div class="flex justify-center gap-1">
                <!-- Edit button -->
                <button class="btn-ghost" type="button">
                  <EditOutlineRoundedIcon height="1rem" />
                </button>

                <!-- Delete button -->
                <button
                  class="btn-ghost text-red-600 hover:text-red-700"
                  type="button"
                >
                  <DeleteIcon height="1rem" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div
        v-else
        class="m-6 flex flex-col items-center rounded-xl border-2 border-dashed border-stone-200 p-10 text-center"
      >
        <h2 class="text-base font-semibold tracking-tight text-stone-800">
          No asset classes tracked yet
        </h2>
        <p class="mt-1.5 text-sm text-stone-500">
          Add your first one to start tracking.
        </p>
        <button class="btn-positive mt-5" type="button">
          Add asset class
        </button>
      </div>
    </div>
  </div>
</template>
