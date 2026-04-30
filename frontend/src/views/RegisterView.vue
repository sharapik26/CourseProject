<template>
  <div class="card" style="max-width: 480px; margin: 2rem auto">
    <h2>Регистрация</h2>
    <div v-if="error" class="error">{{ error }}</div>
    <form @submit.prevent="onSubmit">
      <label>Email</label>
      <input
        v-model="email"
        type="email"
        required
        autocomplete="email"
        @blur="checkEmail"
      />
      <div v-if="!emailValid" class="error" style="margin-top:4px">формат email некорректен</div>

      <label>Имя пользователя</label>
      <input v-model="username" type="text" required minlength="3" maxlength="64" />

      <label>Отображаемое имя (необязательно)</label>
      <input v-model="displayName" type="text" />

      <label>Пароль</label>
      <div class="password-wrap">
        <input
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          required
          minlength="6"
          autocomplete="new-password"
        />
        <button
          type="button"
          class="toggle-pass"
          @click="showPassword = !showPassword"
          :title="showPassword ? 'Скрыть пароль' : 'Показать пароль'"
        >
          {{ showPassword ? '🙈' : '👁' }}
        </button>
      </div>

      <label>Повторите пароль</label>
      <div class="password-wrap">
        <input
          v-model="passwordConfirm"
          :type="showConfirm ? 'text' : 'password'"
          required
          minlength="6"
          autocomplete="new-password"
        />
        <button
          type="button"
          class="toggle-pass"
          @click="showConfirm = !showConfirm"
          :title="showConfirm ? 'Скрыть пароль' : 'Показать пароль'"
        >
          {{ showConfirm ? '🙈' : '👁' }}
        </button>
      </div>

      <div style="margin-top: 1rem; display: flex; gap: 0.5rem">
        <button class="button" :disabled="loading">
          {{ loading ? 'Отправка кода…' : 'Отправить код на email' }}
        </button>
        <router-link to="/login" class="button secondary">Уже есть аккаунт</router-link>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import * as api from "../services/api";

const router = useRouter();

const email = ref("");
const username = ref("");
const displayName = ref("");
const password = ref("");
const passwordConfirm = ref("");
const showPassword = ref(false);
const showConfirm = ref(false);
const error = ref("");
const loading = ref(false);
const emailValid = ref(true);

function checkEmail() {
  emailValid.value = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value);
}

async function onSubmit() {
  error.value = "";
  checkEmail();
  if (!emailValid.value) {
    error.value = "введите корректный email";
    return;
  }
  if (password.value !== passwordConfirm.value) {
    error.value = "Пароли не совпадают";
    return;
  }
  loading.value = true;
  try {
    await api.auth.requestRegister({
      email: email.value,
      username: username.value,
      password: password.value,
      display_name: displayName.value,
    });
    router.push({ path: "/verify-email", query: { email: email.value } });
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка регистрации";
  } finally {
    loading.value = false;
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