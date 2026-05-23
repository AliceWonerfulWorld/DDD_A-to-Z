<script setup lang="ts">
import AppLayout from "@/Layouts/AppLayout.vue";

interface OperationLog {
  id: number;
  action: string;
  target_type: string;
  target_id: string;
  payload: Record<string, unknown>;
  created_at: string;
}

defineProps<{
  operationLogs: OperationLog[];
  appLogLines: string[];
}>();

const actionLabel = (action: string) =>
  action === "created" ? "作成" : "更新";

const actionClass = (action: string) =>
  action === "created"
    ? "bg-green-100 text-green-800"
    : "bg-blue-100 text-blue-800";
</script>

<template>
  <AppLayout>
    <h1 class="text-2xl font-bold mb-6">ログ</h1>

    <section class="mb-10">
      <h2 class="text-lg font-semibold mb-3">操作ログ</h2>
      <div class="bg-white rounded shadow overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="bg-gray-50 border-b">
            <tr>
              <th class="px-4 py-3 text-left">日時</th>
              <th class="px-4 py-3 text-left">操作</th>
              <th class="px-4 py-3 text-left">種別</th>
              <th class="px-4 py-3 text-left">対象</th>
              <th class="px-4 py-3 text-left">内容</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="operationLogs.length === 0">
              <td colspan="5" class="px-4 py-6 text-center text-gray-400">
                まだ操作ログはありません
              </td>
            </tr>
            <tr
              v-for="log in operationLogs"
              :key="log.id"
              class="border-b last:border-0 hover:bg-gray-50"
            >
              <td class="px-4 py-3 whitespace-nowrap text-gray-500">
                {{ new Date(log.created_at).toLocaleString("ja-JP") }}
              </td>
              <td class="px-4 py-3">
                <span
                  :class="actionClass(log.action)"
                  class="px-2 py-0.5 rounded text-xs font-medium"
                >
                  {{ actionLabel(log.action) }}
                </span>
              </td>
              <td class="px-4 py-3 text-gray-600">{{ log.target_type }}</td>
              <td class="px-4 py-3 font-mono text-gray-800">{{ log.target_id }}</td>
              <td class="px-4 py-3 text-gray-500 font-mono text-xs">
                {{ JSON.stringify(log.payload) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section>
      <h2 class="text-lg font-semibold mb-3">アプリログ（直近 200 行）</h2>
      <pre class="bg-gray-900 text-gray-100 text-xs rounded shadow p-4 overflow-x-auto max-h-[600px] overflow-y-auto whitespace-pre-wrap">{{ appLogLines.join("\n") || "ログファイルが空です" }}</pre>
    </section>
  </AppLayout>
</template>
