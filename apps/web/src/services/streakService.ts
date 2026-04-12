import api from "./api";
import { UserStreak } from "@/types/streak";

const streakService = {
  getUserStreak: async (userId: number): Promise<UserStreak> => {
    const response = await api.get<UserStreak>(`/users/${userId}/streak`);
    return response.data;
  },
  getLeaderboard: async (limit: number = 20): Promise<UserStreak[]> => {
    const response = await api.get<UserStreak[]>(`/streaks/leaderboard?limit=${limit}`);
    return response.data;
  },
};
export default streakService;
