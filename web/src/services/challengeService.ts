import api from "./api";
import {
  Challenge,
  ChallengeFrequency,
  LeaderboardEntry,
  Season,
  UserChallenge,
} from "@/types/challenge";

const challengeService = {
  getActiveSeason: async (): Promise<Season> => {
    const response = await api.get<Season>("/seasons/current");
    return response.data;
  },

  getChallenges: async (frequency?: ChallengeFrequency): Promise<Challenge[]> => {
    const response = await api.get<Challenge[]>("/challenges", {
      params: { frequency },
    });
    return response.data || [];
  },

  getMyChallenges: async (): Promise<UserChallenge[]> => {
    const response = await api.get<UserChallenge[]>("/me/challenges");
    return response.data || [];
  },

  getLeaderboard: async (limit: number = 20, offset: number = 0): Promise<LeaderboardEntry[]> => {
    const response = await api.get<LeaderboardEntry[]>("/seasons/leaderboard", {
      params: { limit, offset },
    });
    return response.data || [];
  },
};

export default challengeService;