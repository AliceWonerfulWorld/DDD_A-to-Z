<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

defineProps<{
  pointTypes: {
    code: string;
    language: string;
    label: string;
  }[];
}>();
</script>

<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">ポイントタイプ一覧</h1>
      <Link href="/point-types/create" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
        + 新規作成
      </Link>
    </div>

    <div class="bg-white rounded shadow overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="text-left px-4 py-3">コード</th>
            <th class="text-left px-4 py-3">言語</th>
            <th class="text-left px-4 py-3">ラベル</th>
            <th class="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pt in pointTypes" :key="`${pt.code}-${pt.language}`" class="border-b hover:bg-gray-50">
            <td class="px-4 py-3 font-mono">{{ pt.code }}</td>
            <td class="px-4 py-3">{{ pt.language || '（なし）' }}</td>
            <td class="px-4 py-3">{{ pt.label }}</td>
            <td class="px-4 py-3">
              <Link
                :href="`/point-types/edit?code=${encodeURIComponent(pt.code)}&language=${encodeURIComponent(pt.language)}`"
                class="text-blue-600 hover:underline"
              >
                編集
              </Link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </AppLayout>
</template>
