import api from "./api";
import { Badge, BadgeWithEarned } from "@/types/badge";

const badgeService = {
  getAllBadges: async (): Promise<Badge[]> => {
    const response = await api.get<Badge[]>("/badges");
    return response.data;
  },
  getUserBadges: async (userId: number): Promise<BadgeWithEarned[]> => {
    const response = await api.get<BadgeWithEarned[]>(`/users/${userId}/badges`);
    return response.data;
  },
};
export default badgeService;
