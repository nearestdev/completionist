import api from "./api";
import { UserStats, HeatmapDay } from "@/types/stats";

const statsService = {
  getUserStats: async (userId: number): Promise<UserStats> => {
    const response = await api.get<UserStats>(`/users/${userId}/stats`);
    return response.data;
  },
  getUserHeatmap: async (userId: number): Promise<HeatmapDay[]> => {
    const response = await api.get<HeatmapDay[]>(`/users/${userId}/stats/heatmap`);
    return response.data;
  },
};
export default statsService;
