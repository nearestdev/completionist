export type ChallengeFrequency = 'daily' | 'weekly' | 'monthly' | 'seasonal' | 'infinite';
export type ChallengeType = 'global' | 'personal';

export interface Season {
  id: number;
  name: string;
  startAt: string;
  endAt: string;
  isActive: boolean;
  createdAt: string;
}

export interface Challenge {
  id: number;
  title: string;
  description?: string;
  type: ChallengeType;
  frequency: ChallengeFrequency;
  xpReward: number;
  criteriaType: string;
}

export interface UserChallenge {
  id: number;
  userId: number;
  challengeId: number;
  seasonId?: number;
  currentProgress: number;
  targetProgress: number;
  isCompleted: boolean;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
  
  challengeTitle?: string;
  challengeDescription?: string;
  xpReward?: number;
}

export interface LeaderboardEntry {
  userId: number;
  username: string;
  totalXp: number;
  rank: number;
}