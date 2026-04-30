<template>
  <div v-if="auth.user" class="card">
    <h2>Профиль</h2>
    <div v-if="okMsg" class="success">{{ okMsg }}</div>
    <div v-if="error" class="error">{{ error }}</div>
    <form @submit.prevent="saveProfile">
      <label>Email (только чтение)</label>
      <input :value="auth.user.email" type="email" disabled />
      <label>Имя пользователя (только чтение)</label>
      <input :value="auth.user.username" type="text" disabled />
      <label>Отображаемое имя</label>
      <input v-model="displayName" type="text" />
      <label>Ссылка на аватар</label>
      <input v-model="avatarUrl" type="text" />

      <h3>Сенсорные предпочтения (1 — низкий порог, 5 — высокий)</h3>
      <div class="sensory-grid">
        <label>Порог шума</label>
        <input type="range" v-model.number="noisePref" min="1" max="5" />
        <span>{{ noisePref }}</span>
        <label>Порог света</label>
        <input type="range" v-model.number="lightPref" min="1" max="5" />
        <span>{{ lightPref }}</span>
        <label>Порог скоплений</label>
        <input type="range" v-model.number="crowdPref" min="1" max="5" />
        <span>{{ crowdPref }}</span>
        <label>Порог запахов</label>
        <input type="range" v-model.number="smellPref" min="1" max="5" />
        <span>{{ smellPref }}</span>
        <label>Визуальный порог</label>
        <input type="range" v-model.number="visualPref" min="1" max="5" />
        <span>{{ visualPref }}</span>
      </div>

      <button class="button" type="submit" :disabled="auth.loading">Сохранить</button>
    </form>
  </div>

  <div class="card">
    <h2>Смена пароля</h2>
    <div v-if="passOk" class="success">{{ passOk }}</div>
    <div v-if="passErr" class="error">{{ passErr }}</div>
    <form @submit.prevent="changePass">
      <label>Текущий пароль</label>
      <input v-model="oldPw" type="password" required />
      <label>Новый пароль</label>
      <input v-model="newPw" type="password" required minlength="6" />
      <button class="button" type="submit">Сменить</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useAuthStore } from "../stores/auth";

const auth = useAuthStore();

const displayName = ref("");
const avatarUrl = ref("");
const noisePref = ref(3);
const lightPref = ref(3);
const crowdPref = ref(3);
const smellPref = ref(3);
const visualPref = ref(3);
const okMsg = ref("");
const error = ref("");

watch(
  () => auth.user,
  (u) => {
    if (u) {
      displayName.value = u.display_name || "";
      avatarUrl.value = u.avatar_url || "";
      noisePref.value = u.noise_pref || 3;
      lightPref.value = u.light_pref || 3;
      crowdPref.value = u.crowd_pref || 3;
      smellPref.value = u.smell_pref || 3;
      visualPref.value = u.visual_pref || 3;
    }
  },
  { immediate: true },
);

async function saveProfile() {
  okMsg.value = "";
  error.value = "";
  try {
    await auth.updateProfile({
      display_name: displayName.value,
      avatar_url: avatarUrl.value,
      noise_pref: noisePref.value,
      light_pref: lightPref.value,
      crowd_pref: crowdPref.value,
      smell_pref: smellPref.value,
      visual_pref: visualPref.value,
    });
    okMsg.value = "Профиль обновлён";
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Ошибка сохранения";
  }
}

const oldPw = ref("");
const newPw = ref("");
const passOk = ref("");
const passErr = ref("");

async function changePass() {
  passOk.value = "";
  passErr.value = "";
  try {
    await auth.changePassword(oldPw.value, newPw.value);
    oldPw.value = "";
    newPw.value = "";
    passOk.value = "Пароль изменён";
  } catch (e) {
    passErr.value = e instanceof Error ? e.message : "Ошибка";
  }
}
</script>