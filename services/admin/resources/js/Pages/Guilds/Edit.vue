<script setup lang="ts">
import { useForm, Link } from "@inertiajs/vue3";
import AppLayout from "@/Layouts/AppLayout.vue";

const props = defineProps<{
  guild: {
    id: string;
    slug: string;
    name: string;
    description: string;
    icon: string;
    color: string;
    sort_order: number;
    current_exp: number;
    guild_level: number;
  };
}>();

const form = useForm({
  name: props.guild.name,
  description: props.guild.description,
  icon: props.guild.icon,
  color: props.guild.color,
  sort_order: props.guild.sort_order,
});

function submit() {
  form.put(`/guilds/${props.guild.id}`);
}
</script>

<template>
  <AppLayout>
    <div class="flex items-center gap-4 mb-6">
      <Link href="/guilds" class="text-gray-500 hover:text-gray-700">← 一覧へ</Link>
      <h1 class="text-2xl font-bold">ギルド編集</h1>
    </div>

    <div class="bg-gray-50 border rounded p-4 mb-6 text-sm text-gray-600 max-w-lg">
      <p><span class="font-medium">ID:</span> {{ guild.id }}</p>
      <p><span class="font-medium">スラッグ:</span> {{ guild.slug }} <span class="text-gray-400">（変更不可）</span></p>
    </div>

    <form @submit.prevent="submit" class="bg-white rounded shadow p-6 space-y-4 max-w-lg">
      <div>
        <label class="block text-sm font-medium mb-1">名前 <span class="text-red-500">*</span></label>
        <input v-model="form.name" type="text" class="w-full border rounded px-3 py-2" />
        <p v-if="form.errors.name" class="text-red-500 text-xs mt-1">{{ form.errors.name }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">説明 <span class="text-red-500">*</span></label>
        <textarea v-model="form.description" rows="3" class="w-full border rounded px-3 py-2" />
        <p v-if="form.errors.description" class="text-red-500 text-xs mt-1">{{ form.errors.description }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">アイコン（絵文字） <span class="text-red-500">*</span></label>
        <input v-model="form.icon" type="text" class="w-full border rounded px-3 py-2" />
        <p v-if="form.errors.icon" class="text-red-500 text-xs mt-1">{{ form.errors.icon }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">カラー <span class="text-red-500">*</span></label>
        <div class="flex gap-2 items-center">
          <input v-model="form.color" type="color" class="w-12 h-10 border rounded cursor-pointer" />
          <input v-model="form.color" type="text" class="flex-1 border rounded px-3 py-2" />
        </div>
        <p v-if="form.errors.color" class="text-red-500 text-xs mt-1">{{ form.errors.color }}</p>
      </div>

      <div>
        <label class="block text-sm font-medium mb-1">表示順序 <span class="text-red-500">*</span></label>
        <input v-model.number="form.sort_order" type="number" min="0" class="w-full border rounded px-3 py-2" />
        <p v-if="form.errors.sort_order" class="text-red-500 text-xs mt-1">{{ form.errors.sort_order }}</p>
      </div>

      <div class="flex gap-3 pt-2">
        <button type="submit" :disabled="form.processing" class="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50">
          更新
        </button>
        <Link href="/guilds" class="px-6 py-2 rounded border hover:bg-gray-50">キャンセル</Link>
      </div>
    </form>
  </AppLayout>
</template>
