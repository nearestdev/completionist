import api from "./api";
import { LastFMAccount, RecentTracksResponse } from "@/types/lastfm";

const lastfmService = {
  lastfmAuth: (): void => {
    const token = localStorage.getItem("token");
    if (token) {
      window.location.href = `${
        api.defaults.baseURL
      }/auth/lastfm?auth_token=${encodeURIComponent(token)}`;
    } else {
      console.error("No auth token found for Last.fm login");
    }
  },

  getMyLastFMAccount: async (): Promise<LastFMAccount> => {
    const response = await api.get<LastFMAccount>("/me/lastfm");
    return response.data;
  },

  getMyRecentTracks: async (limit: number = 20): Promise<RecentTracksResponse> => {
    const response = await api.get<RecentTracksResponse>("/me/lastfm/recent", {
      params: { limit },
    });
    return response.data;
  },
};

export default lastfmService;