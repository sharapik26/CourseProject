<template>
  <div class="layout">
    <header class="topbar">
      <div class="brand">Сенсорный навигатор · Пользователи и отзывы</div>
      <nav>
        <router-link v-if="auth.isAuthenticated" to="/profile">Профиль</router-link>
        <router-link v-if="auth.isAuthenticated" to="/reviews">Мои отзывы</router-link>
        <router-link v-if="auth.isAuthenticated" to="/favorites">Избранное</router-link>
        <router-link v-if="auth.isAuthenticated" to="/new-review">Добавить отзыв</router-link>
      </nav>
      <div class="user">
        <template v-if="auth.isAuthenticated">
          {{ auth.user?.display_name || auth.user?.username }}
          <button class="button secondary" style="margin-left: 0.5rem" @click="onLogout">Выйти</button>
        </template>
        <template v-else>
          <router-link to="/login">Вход</router-link> ·
          <router-link to="/register">Регистрация</router-link>
        </template>
      </div>
    </header>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "./stores/auth";

const auth = useAuthStore();
const router = useRouter();

onMounted(async () => {
  if (auth.isAuthenticated) {
    await auth.fetchMe();
  }
});

function onLogout() {
  auth.logout();
  router.push("/login");
}
</script>