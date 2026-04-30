<template>
  <div class="card">
    <h2>Избранные места</h2>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="loading">Загрузка...</div>
    <div v-else-if="!items.length" class="meta">
      Вы ещё не добавили места в избранное.
      В индивидуальном модуле можно использовать демо-id 1, 2, 3.
    </div>

    <div v-for="p in items" :key="p.id" class="card" style="margin-top: 0.5rem">
      <h3>{{ p.name }}</h3>
      <div class="meta">id: {{ p.id }}</div>
      <button class="button danger" @click="onRemove(p.id)">Удалить из избранного</button>
    </div>

    <div class="card" style="margin-top: 1rem; background: #fafafa">
      <h3>Добавить место в избранное по id</h3>
      <div class="row">
        <input v-model.number="newPlaceId" type="number" min="1" placeholder="ID места" />
        <button class="button" @click="onAdd">Добавить</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import * as api from "../services/api";
import type { PlaceRef } from "../types/models";

const items = ref<PlaceRef[]>([]);
const loading = ref(true);
const error = ref("");
const newPlaceId = ref(1);

async function load() {
  loading.value = true;
  error.value = "";
  try {
    const res = await api.favorites.list();
    items.value = res.items;
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка";
  } finally {
    loading.value = false;
  }
}

async function onRemove(id: number) {
  try {
    await api.favorites.remove(id);
    await load();
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка удаления";
  }
}

async function onAdd() {
  try {
    await api.favorites.add(newPlaceId.value);
    await load();
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка добавления";
  }
}

onMounted(load);
</script>