<template>
  <div class="card" style="max-width: 460px; margin: 2rem auto">
    <h2>Подтверждение email</h2>
    <div v-if="error" class="error">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <p class="meta">
      Мы отправили 6-значный код на адрес <strong>{{ email }}</strong>.
      Введите его ниже, чтобы завершить регистрацию.
    </p>

    <form @submit.prevent="onSubmit">
      <label>Код подтверждения</label>
      <input
        v-model="code"
        type="text"
        inputmode="numeric"
        pattern="\d{6}"
        maxlength="6"
        minlength="6"
        required
        autocomplete="one-time-code"
        placeholder="123456"
        class="code-input"
      />

      <div style="margin-top: 1rem; display: flex; gap: 0.5rem">
        <button class="button" :disabled="loading">
          {{ loading ? 'Проверка…' : 'Подтвердить' }}
        </button>
        <button
          type="button"
          class="button secondary"
          @click="onResend"
          :disabled="resending"
        >
          {{ resending ? 'Отправка…' : 'Отправить ещё раз' }}
        </button>
      </div>
    </form>

    <div style="margin-top: 1rem">
      <router-link to="/register">← к регистрации</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import * as api from "../services/api";
import { useAuthStore } from "../stores/auth";

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

const email = ref<string>((route.query.email as string) || "");
const code = ref("");
const error = ref("");
const success = ref("");
const loading = ref(false);
const resending = ref(false);

onMounted(() => {
  if (!email.value) router.replace("/register");
});

async function onSubmit() {
  error.value = "";
  success.value = "";
  if (!/^\d{6}$/.test(code.value)) {
    error.value = "введите 6-значный код";
    return;
  }
  loading.value = true;
  try {
    const res = await api.auth.confirmRegister({
      email: email.value,
      code: code.value,
    });
    auth.applyToken(res.token, res.user);
    router.push("/profile");
  } catch (e) {
    error.value = e instanceof Error ? e.message : "ошибка проверки кода";
  } finally {
    loading.value = false;
  }
}

async function onResend() {
  error.value = "";
  success.value = "";
  resending.value = true;
  try {
    await api.auth.resendCode({ email: email.value });
    success.value = "новый код отправлен";
  } catch (e) {
    error.value = e instanceof Error ? e.message : "ошибка отправки";
  } finally {
    resending.value = false;
  }
}
</script>

<style scoped>
.code-input {
  font-size: 22px;
  letter-spacing: 8px;
  text-align: center;
  font-family: ui-monospace, "JetBrains Mono", Consolas, monospace;
}
.success {
  color: #1b5e20;
  background: #e8f5e9;
  border: 1px solid #a5d6a7;
  padding: 6px 10px;
  border-radius: 6px;
  margin-bottom: 8px;
}
</style>