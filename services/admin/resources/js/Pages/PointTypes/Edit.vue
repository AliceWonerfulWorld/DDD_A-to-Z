<script setup lang="ts">
import { useForm, Link } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

const props = defineProps<{
  pointType: {
    code: string;
    language: string;
    label: string;
  };
}>();

const form = useForm({
  code: props.pointType.code,
  language: props.pointType.language,
  label: props.pointType.label,
});

function submit() {
  form.put("/point-types");
}
</script>

<template>
  <AppLayout>
    <div class="flex items-center gap-4 mb-6">
      <Link href="/point-types" class="text-gray-500 hover:text-gray-700">← 一覧へ</Link>
      <h1 class="text-2xl font-bold">ポイントタイプ編集</h1>
    </div>

    <form @submit.prevent="submit" class="bg-white rounded shadow p-6 space-y-4 max-w-lg">
      <div>
        <label class="block text-sm font-medium mb-1">コード</label>
        <input :value="pointType.code" type="text" disabled class="w-full border rounded px-3 py-2 font-mono bg-gray-50 text-gray-500 cursor-not-allowed" />
        <p class="text-gray-400 text-xs mt-1">変更不可</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">言語</label>
        <input :value="pointType.language || '（なし）'" type="text" disabled class="w-full border rounded px-3 py-2 bg-gray-50 text-gray-500 cursor-not-allowed" />
        <p class="text-gray-400 text-xs mt-1">変更不可</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">ラベル <span class="text-red-500">*</span></label>
        <input v-model="form.label" type="text" class="w-full border rounded px-3 py-2" />
        <p v-if="form.errors.label" class="text-red-500 text-xs mt-1">{{ form.errors.label }}</p>
      </div>

      <div class="flex gap-3 pt-2">
        <button type="submit" :disabled="form.processing" class="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50">
          更新
        </button>
        <Link href="/point-types" class="px-6 py-2 rounded border hover:bg-gray-50">キャンセル</Link>
      </div>
    </form>
  </AppLayout>
</template>
