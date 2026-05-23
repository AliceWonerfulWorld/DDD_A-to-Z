<script setup lang="ts">
import { Link } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

defineProps<{
  guilds: {
    id: string;
    slug: string;
    name: string;
    description: string;
    icon: string;
    color: string;
    sort_order: number;
    current_exp: number;
    guild_level: number;
  }[];
}>();
</script>

<template>
  <AppLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">ギルド一覧</h1>
      <Link href="/guilds/create" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
        + 新規作成
      </Link>
    </div>

    <div class="bg-white rounded shadow overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="text-left px-4 py-3">順序</th>
            <th class="text-left px-4 py-3">アイコン</th>
            <th class="text-left px-4 py-3">名前</th>
            <th class="text-left px-4 py-3">スラッグ</th>
            <th class="text-left px-4 py-3">カラー</th>
            <th class="text-left px-4 py-3">レベル</th>
            <th class="text-left px-4 py-3">EXP</th>
            <th class="px-4 py-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="guild in guilds" :key="guild.id" class="border-b hover:bg-gray-50">
            <td class="px-4 py-3">{{ guild.sort_order }}</td>
            <td class="px-4 py-3 text-xl">{{ guild.icon }}</td>
            <td class="px-4 py-3 font-medium">{{ guild.name }}</td>
            <td class="px-4 py-3 text-gray-500">{{ guild.slug }}</td>
            <td class="px-4 py-3">
              <span class="inline-flex items-center gap-2">
                <span class="w-4 h-4 rounded-full border" :style="{ backgroundColor: guild.color }" />
                {{ guild.color }}
              </span>
            </td>
            <td class="px-4 py-3">{{ guild.guild_level }}</td>
            <td class="px-4 py-3">{{ guild.current_exp }}</td>
            <td class="px-4 py-3">
              <Link :href="`/guilds/${guild.id}/edit`" class="text-blue-600 hover:underline">編集</Link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </AppLayout>
</template>
