<script setup lang="ts">
  import { useHead } from "nuxt/app";

  import AssetClassTable from "~/components/asset/AssetClassTable.vue";
  import type { AssetClass } from "~/openapi";

  // Add a title for the page
  useHead({
    title: "Dashboard",
  });

  // Sample (hardcoded temporarily) set of asset classes
  const assetClasses: AssetClass[] = [
    {
      description:
        "Ownership shares in public or private companies, offering long-term capital appreciation and dividend income.",
      id: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
      name: "Equity",
    },
    {
      description:
        "Bonds and debt securities issued by governments or corporations that pay regular interest and return principal at maturity.",
      id: "b28f89e1-638e-4a94-82a9-02fb232a2f8c",
      name: "Fixed Income",
    },
    {
      description:
        "Low-risk, highly liquid short-term investments like Treasury bills, money market funds, and certificates of deposit.",
      id: "c8e9b012-7e04-4b53-911e-812e96d9f821",
      name: "Cash & Cash Equivalents",
    },
    {
      description:
        "Physical real estate assets, REITs, and property funds generating income through rental yields and long-term appreciation.",
      id: "d41f38e3-05c2-4a87-95bd-44a631bf345d",
      name: "Real Estate",
    },
    {
      description:
        "Raw materials and primary products including precious metals, energy resources, and agricultural goods used often as inflation hedges.",
      id: "e50b1827-023a-4f46-a5df-192561cd62a9",
      name: "Commodities",
    },
    {
      description:
        "Blockchain-native digital tokens and currencies characterized by higher volatility and alternative yield mechanisms.",
      id: "f72c3d10-8b1e-4503-9df2-38d1720a23b1",
      name: "Cryptocurrency & Digital Assets",
    },
    {
      description:
        "Direct investments or venture capital holdings in private companies not listed on public stock exchanges.",
      id: "01f8d4e9-1234-45a8-bc91-561238ef9012",
      name: "Private Equity",
    },
    {
      description:
        "Investments in essential public services and physical networks such as energy grids, toll roads, utilities, and telecommunications.",
      id: "19a2b8e3-4c5d-4f12-8e9a-0123456789bc",
      name: "Infrastructure",
    },
  ];
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
          {{ assetClasses.length }}
          {{ assetClasses.length === 1 ? "asset class" : "asset classes" }}
          tracked
        </p>
      </div>
      <button class="btn-positive" type="button">Add asset class</button>
    </header>

    <!-- Table containing the list of asset classes currently tracked -->
    <main class="mt-8">
      <AssetClassTable :asset-classes="assetClasses" />
    </main>
  </div>
</template>
