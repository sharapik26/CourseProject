import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";
import { useAuthStore } from "../stores/auth";

const routes: RouteRecordRaw[] = [
  { path: "/", redirect: "/profile" },
  { path: "/login", component: () => import("../views/LoginView.vue") },
  { path: "/register", component: () => import("../views/RegisterView.vue") },
  { path: "/verify-email", component: () => import("../views/VerifyEmailView.vue") },
  {
    path: "/profile",
    component: () => import("../views/ProfileView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/reviews",
    component: () => import("../views/MyReviewsView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/favorites",
    component: () => import("../views/FavoritesView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/new-review",
    component: () => import("../views/NewReviewView.vue"),
    meta: { requiresAuth: true },
  },
];

const router = createRouter({ history: createWebHistory(), routes });

router.beforeEach((to) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { path: "/login" };
  }
  return true;
});

export default router;