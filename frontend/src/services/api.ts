import type {
  AuthResponse,
  Review,
  Aggregate,
  User,
  PlaceRef,
} from "../types/models";

const API_BASE =
  (import.meta.env.VITE_API_BASE as string) || "http://localhost:8081";

class ApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = "ApiError";
  }
}

let currentToken: string | null = localStorage.getItem("snu_token");

export function setToken(token: string | null): void {
  currentToken = token;
  if (token) localStorage.setItem("snu_token", token);
  else localStorage.removeItem("snu_token");
}

export function getToken(): string | null {
  return currentToken;
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  authRequired = false,
): Promise<T> {
  const headers: Record<string, string> = {
    Accept: "application/json",
    ...((options.headers as Record<string, string>) || {}),
  };
  if (options.body && !(options.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }
  if (authRequired || currentToken) {
    if (!currentToken) {
      throw new ApiError("требуется авторизация", 401);
    }
    headers["Authorization"] = `Bearer ${currentToken}`;
  }
  const res = await fetch(API_BASE + path, { ...options, headers });
  if (res.status === 204) {
    return undefined as unknown as T;
  }
  const text = await res.text();
  const data = text ? JSON.parse(text) : ({} as Record<string, unknown>);
  if (!res.ok) {
    const msg = (data as { error?: string }).error || res.statusText;
    throw new ApiError(msg, res.status);
  }
  return data as T;
}

export const auth = {
  register: (input: {
    email: string;
    username: string;
    password: string;
    display_name?: string;
  }) =>
    request<AuthResponse>("/api/auth/register", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  // Шаг 1 регистрации: проверка email + отправка 6-значного кода.
  requestRegister: (input: {
    email: string;
    username: string;
    password: string;
    display_name?: string;
  }) =>
    request<{ status: string; expires_at: string }>("/api/auth/register-request", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  // Шаг 2: подтверждение кода → создание учётной записи и выдача JWT.
  confirmRegister: (input: { email: string; code: string }) =>
    request<AuthResponse>("/api/auth/register-confirm", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  resendCode: (input: { email: string }) =>
    request<{ status: string; expires_at: string }>("/api/auth/resend-code", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  login: (input: { email: string; password: string }) =>
    request<AuthResponse>("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(input),
    }),
  me: () => request<User>("/api/users/me", {}, true),
  updateMe: (input: Partial<User>) =>
    request<User>("/api/users/me", {
      method: "PUT",
      body: JSON.stringify(input),
    }, true),
  changePassword: (old_password: string, new_password: string) =>
    request<void>("/api/users/me/password", {
      method: "PUT",
      body: JSON.stringify({ old_password, new_password }),
    }, true),
};

export const reviews = {
  byPlace: (placeId: number) =>
    request<{ items: Review[] }>(`/api/places/${placeId}/reviews`),
  aggregate: (placeId: number) =>
    request<Aggregate>(`/api/places/${placeId}/aggregate`),
  create: (
    placeId: number,
    input: {
      text?: string;
      noise: number;
      light: number;
      crowd: number;
      smell: number;
      visual: number;
    },
  ) =>
    request<Review>(`/api/places/${placeId}/reviews`, {
      method: "POST",
      body: JSON.stringify(input),
    }, true),
  update: (
    id: number,
    input: {
      text?: string;
      noise: number;
      light: number;
      crowd: number;
      smell: number;
      visual: number;
    },
  ) =>
    request<Review>(`/api/reviews/${id}`, {
      method: "PUT",
      body: JSON.stringify(input),
    }, true),
  remove: (id: number) =>
    request<void>(`/api/reviews/${id}`, { method: "DELETE" }, true),
  myReviews: () =>
    request<{ items: Review[] }>(`/api/users/me/reviews`, {}, true),
};

export const favorites = {
  list: () =>
    request<{ items: PlaceRef[] }>("/api/users/me/favorites", {}, true),
  add: (placeId: number) =>
    request<void>(`/api/favorites/${placeId}`, { method: "POST" }, true),
  remove: (placeId: number) =>
    request<void>(`/api/favorites/${placeId}`, { method: "DELETE" }, true),
};

export { ApiError };