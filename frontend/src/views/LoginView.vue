<template>
  <div class="card" style="max-width: 420px; margin: 2rem auto">
    <h2>Вход</h2>
    <div v-if="error" class="error">{{ error }}</div>
    <form @submit.prevent="onSubmit">
      <label>Email</label>
      <input v-model="email" type="email" required autocomplete="email" />
      <label>Пароль</label>
      <div class="password-wrap">
        <input
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          required
          minlength="6"
          autocomplete="current-password"
        />
        <button
          type="button"
          class="toggle-pass"
          @click="showPassword = !showPassword"
          :title="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
          :aria-label="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
        >
          {{ showPassword ? '🙈' : '👁' }}
        </button>
      </div>
      <div style="margin-top: 1rem; display: flex; gap: 0.5rem">
        <button class="button" :disabled="auth.loading">Войти</button>
        <router-link to="/register" class="button secondary">Регистрация</router-link>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const auth = useAuthStore();
const router = useRouter();

const email = ref("");
const password = ref("");
const showPassword = ref(false);
const error = ref("");

async function onSubmit() {
  error.value = "";
  try {
    await auth.login(email.value, password.value);
    router.push("/profile");
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка входа";
  }
}
</script>

<style scoped>
.password-wrap {
  position: relative;
  display: flex;
  align-items: center;
}
.password-wrap input { width: 100%; padding-right: 44px; }
.toggle-pass {
  position: absolute; right: 6px; top: 50%; transform: translateY(-50%);
  background: transparent; border: 0; padding: 4px 8px;
  cursor: pointer; font-size: 18px; line-height: 1; color: #6b7280;
}
.toggle-pass:hover { color: var(--fg, #111); }
</style>