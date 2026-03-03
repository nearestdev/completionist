export interface UserFollowResponse {
  userId: number;
  username: string;
  followedAt: string;
}

export interface UserStatsResponse {
  totalCompleted: number;
  totalPlanning: number;
  totalCurrent: number;
  totalPaused: number;
  totalDropped: number;
  followersCount: number;
  followingCount: number;
}

export interface EnhancedUserProfile {
  id: number;
  username: string;
  bio?: string;
  favoriteGenres?: string[];
  createdAt: string;
  updatedAt: string;
  stats: UserStatsResponse;
  isFollowing?: boolean;
  isFollowedBy?: boolean;
}