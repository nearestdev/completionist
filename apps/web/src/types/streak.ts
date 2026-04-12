export type StreakTier = "none" | "spark" | "fire" | "inferno";

export interface UserStreak {
  userId: number;
  currentStreak: number;
  longestStreak: number;
  lastActivityDate?: string;
  tier: StreakTier;
  updatedAt: string;
}
