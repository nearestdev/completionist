export type BadgeCriteriaType = "item_completion_count" | "item_type_count" | "streak_threshold" | "franchise_completion";

export interface Badge {
  id: number;
  code: string;
  name: string;
  description?: string;
  icon?: string;
  criteriaType: BadgeCriteriaType;
  criteriaJson: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
}

export interface BadgeWithEarned extends Badge {
  earnedAt?: string;
}
