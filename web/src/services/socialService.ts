import api from "./api";
import { EnhancedUserProfile, UserFollowResponse } from "@/types/social";

const socialService = {
  followUser: async (userId: number): Promise<void> => {
    await api.post(`/users/${userId}/follow`);
  },
  unfollowUser: async (userId: number): Promise<void> => {
    await api.delete(`/users/${userId}/unfollow`);
  },
  getUserFollowers: async (userId: number): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>(`/users/${userId}/followers`);
    return response.data;
  },
  getUserFollowing: async (userId: number): Promise<UserFollowResponse[]> => {
    const response = await api.get<UserFollowResponse[]>(`/users/${userId}/following`);
    return response.data;
  },
  getUserProfile: async (userId: number): Promise<EnhancedUserProfile> => {
    const response = await api.get<EnhancedUserProfile>(`/users/${userId}/profile`);
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