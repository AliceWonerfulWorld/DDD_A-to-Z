<script setup lang="ts">
import { Link, usePage } from "@inertiajs/vue3";

const page = usePage<{ flash: { success?: string; error?: string } }>();
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-gray-900 text-white px-6 py-4 flex items-center gap-6">
      <span class="font-bold text-lg">Lang War Admin</span>
      <Link href="/guilds" class="hover:text-gray-300">ギルド</Link>
      <Link href="/point-types" class="hover:text-gray-300">ポイントタイプ</Link>
    </nav>

    <div v-if="page.props.flash.success" class="bg-green-100 border border-green-400 text-green-800 px-6 py-3 m-4 rounded">
      {{ page.props.flash.success }}
    </div>
    <div v-if="page.props.flash.error" class="bg-red-100 border border-red-400 text-red-800 px-6 py-3 m-4 rounded">
      {{ page.props.flash.error }}
    </div>

    <main class="max-w-5xl mx-auto py-8 px-4">
      <slot />
    </main>
  </div>
</template>
