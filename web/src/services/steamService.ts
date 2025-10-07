import api from "./api";
import {
  AttachSteamRequest,
  SteamAccount,
  SteamAchievements,
  SteamGameSchema,
  SteamOwnedGames,
} from "@/types/steam";

const steamService = {
  steamLogin: (): void => {
    const token = localStorage.getItem("token");
    if (token) {
      window.location.href = `${
        api.defaults.baseURL
      }/auth/steam/login?auth_token=${encodeURIComponent(token)}`;
    } else {
      console.error("No auth token found for Steam login");
    }
  },

  getMySteamAccount: async (): Promise<SteamAccount> => {
    const response = await api.get<SteamAccount>("/me/steam");
    return response.data;
  },

  getUserSteamAccount: async (userId: number): Promise<SteamAccount> => {
    const response = await api.get<SteamAccount>(`/users/${userId}/steam`);
    return response.data;
  },

  attachSteamToUser: async (data: AttachSteamRequest): Promise<SteamAccount> => {
    const response = await api.post<SteamAccount>("/me/steam/attach", data);
    return response.data;
  },

  getMySteamOwnedGames: async (): Promise<SteamOwnedGames> => {
    const response = await api.get<SteamOwnedGames>("/me/steam/owned");
    return response.data;
  },

  getMySteamAchievementsForApp: async (appId: number): Promise<SteamAchievements> => {
    const response = await api.get<SteamAchievements>(`/steam/achievements/${appId}`);
    return response.data;
  },

  getSteamGameSchema: async (appId: number): Promise<SteamGameSchema> => {
    const response = await api.get<SteamGameSchema>(`/steam/schema/${appId}`);
    return response.data;
  },
};

export default steamService;