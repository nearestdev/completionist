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

export interface ProgressData {
  current?: number;
  total?: number;
  unit?: "pages" | "chapters" | "volumes" | "episodes" | "tracks" | "albums";
  hoursPlayed?: number;
  achievementsUnlocked?: number;
  achievementsTotal?: number;
  watched?: boolean;
  season?: number;
  tracksListened?: number;
  totalTracks?: number;
  albumsListened?: number;
  totalAlbums?: number;
}

export interface UserListItem {
  id: number;
  userId: number;
  mediaItemId: string;
  status: ItemStatus;
  progressData?: ProgressData;
  rating?: number;
  itemType: ItemType;
  createdAt: string;
  updatedAt: string;
  media?: MediaItem;
}

export interface WishlistItem {
  id: number;
  userId: number;
  mediaItemId: string;
  priority: PriorityLevel;
  priceCents?: number;
  createdAt: string;
  updatedAt: string;
  media?: MediaItem;
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
  progressData?: ProgressData;
  rating?: number;
}

export interface NewUserListItem {
  mediaItemId: string;
  status: ItemStatus;
  progressData?: ProgressData;
  rating?: number;
}


export interface UpdateUserListItem {
  status: ItemStatus;
  progressData?: ProgressData;
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