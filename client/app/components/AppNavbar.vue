<template>
  <nav
    class="sticky top-0 z-10 border-b border-amber-200 shadow-sm shadow-amber-900/5"
  >
    <div
      class="mx-auto flex h-14 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6"
    >
      <NuxtLink
        to="/dashboard"
        aria-current="page"
        class="text-md rounded-md px-3 py-1.5 font-medium text-stone-900"
      >
        Overview
      </NuxtLink>
      <AddAssetClassDialog />
    </div>
  </nav>
</template>
