export interface User {
  id: number;
  email: string;
  username: string;
  display_name: string;
  avatar_url: string;
  noise_pref: number;
  light_pref: number;
  crowd_pref: number;
  smell_pref: number;
  visual_pref: number;
}

export interface AuthResponse {
  token: string;
  expires_at: string;
  user: User;
}

export interface Review {
  id: number;
  place_id: number;
  user_id: number;
  user?: User;
  text: string;
  noise: number;
  light: number;
  crowd: number;
  smell: number;
  visual: number;
  created_at: string;
  updated_at: string;
}

export interface Aggregate {
  place_id: number;
  avg_noise: number;
  avg_light: number;
  avg_crowd: number;
  avg_smell: number;
  avg_visual: number;
  overall_avg: number;
  reviews_count: number;
}

export interface PlaceRef {
  id: number;
  name: string;
}