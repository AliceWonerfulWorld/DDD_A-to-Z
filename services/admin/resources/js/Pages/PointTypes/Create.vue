<script setup lang="ts">
import { useForm, Link } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

const form = useForm({
  code: "",
  language: "",
  label: "",
});

function submit() {
  form.post("/point-types");
}
</script>

<template>
  <AppLayout>
    <div class="flex items-center gap-4 mb-6">
      <Link href="/point-types" class="text-gray-500 hover:text-gray-700">← 一覧へ</Link>
      <h1 class="text-2xl font-bold">ポイントタイプ作成</h1>
    </div>

    <form @submit.prevent="submit" class="bg-white rounded shadow p-6 space-y-4 max-w-lg">
      <div>
        <label class="block text-sm font-medium mb-1">コード <span class="text-red-500">*</span></label>
        <input v-model="form.code" type="text" class="w-full border rounded px-3 py-2 font-mono" placeholder="SP" />
        <p class="text-gray-500 text-xs mt-1">作成後は変更できません</p>
        <p v-if="form.errors.code" class="text-red-500 text-xs mt-1">{{ form.errors.code }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">言語</label>
        <input v-model="form.language" type="text" class="w-full border rounded px-3 py-2" placeholder="Go（空白可）" />
        <p class="text-gray-500 text-xs mt-1">作成後は変更できません</p>
        <p v-if="form.errors.language" class="text-red-500 text-xs mt-1">{{ form.errors.language }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">ラベル <span class="text-red-500">*</span></label>
        <input v-model="form.label" type="text" class="w-full border rounded px-3 py-2" placeholder="Go Skill Point" />
        <p v-if="form.errors.label" class="text-red-500 text-xs mt-1">{{ form.errors.label }}</p>
      </div>

      <div class="flex gap-3 pt-2">
        <button type="submit" :disabled="form.processing" class="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50">
          作成
        </button>
        <Link href="/point-types" class="px-6 py-2 rounded border hover:bg-gray-50">キャンセル</Link>
      </div>
    </form>
  </AppLayout>
</template>
