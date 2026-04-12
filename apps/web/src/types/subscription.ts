export type SubscriptionStatus = "active" | "past_due" | "cancelled" | "incomplete" | "trialing";

export interface Subscription {
  id: number;
  userId: number;
  status: SubscriptionStatus;
  currentPeriodStart?: string;
  currentPeriodEnd?: string;
  cancelAtPeriodEnd: boolean;
  cancelledAt?: string;
  createdAt: string;
  updatedAt: string;
}
