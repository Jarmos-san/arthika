<script setup lang="ts">
  import { useHead } from "nuxt/app";

  import useAssetClasses from "~/composables/useAssetClasses";

  // Add a title for the page
  useHead({
    title: "Dashboard",
  });

  // Fetch the asset class data from the server
  const { data: assetClasses, refresh } = useAssetClasses();
</script>

<template>
  <div class="mx-auto w-full max-w-4xl px-6 py-10">
    <!-- Header section containing the button to add new assets among other features -->
    <header class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight text-stone-800">
          Asset Classes
        </h1>
        <p class="mt-1 text-sm text-stone-500">
          {{ assetClasses?.length }}
          {{ assetClasses?.length === 1 ? "asset class" : "asset classes" }}
          tracked
        </p>
      </div>
      <AssetAddForm @saved="() => refresh()" />
    </header>

    <!-- Table containing the list of asset classes currently tracked -->
    <main class="mt-8">
      <AssetDescriptionTable
        :asset-classes="assetClasses"
        @deleted="() => refresh()"
      />
    </main>
  </div>
</template>
