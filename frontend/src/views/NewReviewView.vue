<template>
  <div class="card">
    <h2>Новый отзыв</h2>
    <p class="meta">
      Введите id места (в индивидуальном модуле доступны демо-id 1, 2, 3),
      выставите сенсорные оценки и при желании напишите текст.
    </p>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="okMsg" class="success">{{ okMsg }}</div>

    <form @submit.prevent="submit">
      <label>ID места</label>
      <input v-model.number="placeId" type="number" min="1" required />

      <label>Текст отзыва</label>
      <textarea v-model="text" maxlength="2000" />

      <h3>Сенсорные оценки (1 — комфортно, 5 — сильное воздействие)</h3>
      <div class="sensory-grid">
        <label>Шум</label>
        <input type="range" v-model.number="noise" min="1" max="5" />
        <span>{{ noise }}</span>
        <label>Освещённость</label>
        <input type="range" v-model.number="light" min="1" max="5" />
        <span>{{ light }}</span>
        <label>Скопления людей</label>
        <input type="range" v-model.number="crowd" min="1" max="5" />
        <span>{{ crowd }}</span>
        <label>Запахи</label>
        <input type="range" v-model.number="smell" min="1" max="5" />
        <span>{{ smell }}</span>
        <label>Визуальная нагрузка</label>
        <input type="range" v-model.number="visual" min="1" max="5" />
        <span>{{ visual }}</span>
      </div>

      <button class="button" type="submit">Опубликовать</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import * as api from "../services/api";

const router = useRouter();
const placeId = ref(1);
const text = ref("");
const noise = ref(3);
const light = ref(3);
const crowd = ref(3);
const smell = ref(3);
const visual = ref(3);
const error = ref("");
const okMsg = ref("");

async function submit() {
  error.value = "";
  okMsg.value = "";
  try {
    await api.reviews.create(placeId.value, {
      text: text.value,
      noise: noise.value,
      light: light.value,
      crowd: crowd.value,
      smell: smell.value,
      visual: visual.value,
    });
    okMsg.value = "Отзыв опубликован";
    setTimeout(() => router.push("/reviews"), 700);
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка отправки";
  }
}
</script>