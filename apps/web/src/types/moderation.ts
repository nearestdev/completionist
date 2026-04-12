export type ModerationEntityType = "post" | "comment" | "room_message" | "direct_message";

export interface NewReport {
  entityType: ModerationEntityType;
  entityId: number;
  reason: string;
}
