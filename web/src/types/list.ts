export type ItemType = "manga" | "game" | "book" | "movie" | "series" | "music";
export type ItemStatus = "planning" | "current" | "completed" | "paused" | "dropped";
export type PriorityLevel = "low" | "medium" | "high";

export interface MediaItem {
  id: string;
  itemType: ItemType;
  source?: string;
  externalId?: string;
  title: string;
  description?: string;
  coverImageUrl?: string;
  releaseDate?: string;
  genres?: string[];
  criticRatingValue?: number;
  criticRatingCount?: number;
  userRatingExternal?: number;
  metadata?: any;
  createdAt: string;
  updatedAt: string;
}

export interface UserListItem {
  id: number;
  userId: number;
  mediaItemId: string;
  status: ItemStatus;
  progress?: string;
  rating?: number;
  createdAt: string;
  updatedAt: string;
}

export interface WishlistItem {
  id: number;
  userId: number;
  mediaItemId: string;
  priority: PriorityLevel;
  priceCents?: number;
  createdAt: string;
  updatedAt: string;
}

export interface NewMediaItem {
  itemType: ItemType;
  source?: string;
  externalId?: string;
  title: string;
  description?: string;
  coverImageUrl?: string;
  releaseDate?: string;
  genres?: string[];
  metadata?: any;
}

export interface NewUserListItemBody {
  status: ItemStatus;
  progress?: string;
  rating?: number;
}

export interface NewUserListItem {
  mediaItemId: string;
  status: ItemStatus;
  progress?: string;
  rating?: number;
}


export interface UpdateUserListItem {
  status: ItemStatus;
  progress?: string;
  rating?: number;
}

export interface NewWishlistItem {
  mediaItemId: string;
  priority: PriorityLevel;
  priceCents?: number;
}

export interface CreateListItemPayload {
  mediaData: NewMediaItem;
  listData: NewUserListItemBody;
}