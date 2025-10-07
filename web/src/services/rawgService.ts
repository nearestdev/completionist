import api from "./api";
import { RAWGGameAchievements, RAWGGameSearch } from "@/types/rawg";

const rawgService = {
  searchGames: async (query: string, page: number = 1): Promise<RAWGGameSearch> => {
    const response = await api.get<RAWGGameSearch>("/games/rawg/search", {
      params: { q: query, page },
    });
    return response.data;
  },

  getGameAchievements: async (rawgId: number, page: number = 1): Promise<RAWGGameAchievements> => {
    const response = await api.get<RAWGGameAchievements>(`/games/rawg/${rawgId}/achievements`, {
      params: { page },
    });
    return response.data;
  },
};

export default rawgService;