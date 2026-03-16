import api from "./api";
import { User } from "@/types/user";

export const adminService = {
  listUsers: async (limit = 50, offset = 0): Promise<User[]> => {
    const response = await api.get<User[]>("/admin/users", {
      params: { limit, offset },
    });
    return response.data;
  },
};

export default adminService;
