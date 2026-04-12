export type ImportJobStatus = "pending" | "running" | "completed" | "failed" | "cancelled";

export interface ConnectedAccount {
  id: number;
  userId: number;
  provider: string;
  providerUserId?: string;
  providerUsername?: string;
  lastSyncedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ImportJob {
  id: number;
  userId: number;
  connectedAccountId?: number;
  provider: string;
  status: ImportJobStatus;
  totalItems: number;
  importedItems: number;
  skippedItems: number;
  failedItems: number;
  startedAt?: string;
  completedAt?: string;
  createdAt: string;
}
