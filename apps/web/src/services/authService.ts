import api from "./api";
import { User, RegisterRequest, LoginRequest, LoginResponse } from "@/types/user";

const authService = {
  register: async (data: RegisterRequest): Promise<User> => {
    const response = await api.post<User>("/register", data);
    return response.data;
  },

  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>("/login", data);
    return response.data;
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<User>("/me");
    return response.data;
  },
};

export default authService;