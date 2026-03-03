import api from "./api";
import { EnhancedUserProfile, UserFollowResponse } from "@/types/social";
const socialService = {
  followUser: async (username: string): Promise<void> => {
    await api.post(`/users/${username}/follow`);
  },
  unfollowUser: async (username: string): Promise<void> => {
    await api.delete(`/users/${username}/unfollow`);
  },
  getUserFollowers: async (username: string): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>(`/users/${username}/followers`);
    return response.data;
  },
  getUserFollowing: async (username: string): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>(`/users/${username}/following`);
    return response.data;
  },
  getUserProfile: async (username: string): Promise<EnhancedUserProfile> => {
    const response = await api.get<EnhancedUserProfile>(`/users/${username}/profile`);
    return response.data;
  },
  getFollowSuggestions: async (): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>("/users/suggestions");
    return response.data;
  },
  searchUsers: async (query: string): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>("/users/search", {
      params: { q: query },
    });
    return response.data;
  },
};
export default socialService;