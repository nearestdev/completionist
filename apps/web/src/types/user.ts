export type UserRole = "user" | "admin";

export interface User {
  id: number;
  username: string;
  email: string;
  role: UserRole;
  xp: number;
  level: number;
  createdAt: string;
  updatedAt: string;
}
export interface RegisterRequest {
  username: string;
  email: string;
  password?: string;
}
export interface LoginRequest {
  email: string;
  password?: string;
}
export interface LoginResponse {
  token: string;
}
