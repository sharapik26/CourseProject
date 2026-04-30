<template>
  <div class="card">
    <h2>Мои отзывы</h2>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="loading">Загрузка...</div>
    <div v-else-if="!items.length" class="meta">У вас пока нет отзывов.</div>

    <div v-for="r in items" :key="r.id" class="card" style="margin-top: 0.5rem">
      <div class="review-meta">
        <span>Место #{{ r.place_id }}</span>
        <span>{{ formatDate(r.created_at) }}</span>
      </div>
      <p>{{ r.text || "(без текста)" }}</p>
      <div class="row">
        <span>Шум: {{ r.noise }}/5</span>
        <span>Свет: {{ r.light }}/5</span>
        <span>Скопл.: {{ r.crowd }}/5</span>
        <span>Запах: {{ r.smell }}/5</span>
        <span>Визуал: {{ r.visual }}/5</span>
      </div>
      <div style="margin-top: 0.5rem">
        <button class="button danger" @click="onDelete(r.id)">Удалить</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import * as api from "../services/api";
import type { Review } from "../types/models";

const items = ref<Review[]>([]);
const loading = ref(true);
const error = ref("");

async function load() {
  loading.value = true;
  error.value = "";
  try {
    const res = await api.reviews.myReviews();
    items.value = res.items;
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка загрузки";
  } finally {
    loading.value = false;
  }
}

async function onDelete(id: number) {
  if (!confirm("Удалить отзыв?")) return;
  try {
    await api.reviews.remove(id);
    await load();
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Не удалось удалить";
  }
}

function formatDate(s: string): string {
  return new Date(s).toLocaleString("ru-RU");
}

onMounted(load);
</script>