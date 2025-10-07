import api from "./api";
import {
  CreateListItemPayload,
  NewWishlistItem,
  UpdateUserListItem,
  UserListItem,
  WishlistItem,
} from "@/types/list";

const listService = {
  createListItem: async (data: CreateListItemPayload): Promise<UserListItem> => {
    const response = await api.post<UserListItem>("/lists", data);
    return response.data;
  },

  getMyListItems: async (): Promise<UserListItem[]> => {
    const response = await api.get<UserListItem[]>("/lists");
    return response.data;
  },

  updateMyListItem: async (
    itemId: number,
    data: UpdateUserListItem
  ): Promise<UserListItem> => {
    const response = await api.patch<UserListItem>(`/lists/${itemId}`, data);
    return response.data;
  },

  deleteMyListItem: async (itemId: number): Promise<void> => {
    await api.delete(`/lists/${itemId}`);
  },

  addAnimeFromJikan: async (
    malId: number,
    data: { status: string; progress?: string; rating?: number }
  ): Promise<UserListItem> => {
    const response = await api.post<UserListItem>(`/lists/jikan/anime/${malId}`, data);
    return response.data;
  },

  addMangaFromJikan: async (
    malId: number,
    data: { status: string; progress?: string; rating?: number }
  ): Promise<UserListItem> => {
    const response = await api.post<UserListItem>(`/lists/jikan/manga/${malId}`, data);
    return response.data;
  },

  addMovieFromTMDB: async (
    tmdbId: number,
    data: { status: string; progress?: string; rating?: number }
  ): Promise<UserListItem> => {
    const response = await api.post<UserListItem>(`/lists/tmdb/movie/${tmdbId}`, data);
    return response.data;
  },

  addTVFromTMDB: async (
    tmdbId: number,
    data: { status: string; progress?: string; rating?: number }
  ): Promise<UserListItem> => {
    const response = await api.post<UserListItem>(`/lists/tmdb/tv/${tmdbId}`, data);
    return response.data;
  },

  addToWishlist: async (data: NewWishlistItem): Promise<WishlistItem> => {
    const response = await api.post<WishlistItem>("/wishlist", data);
    return response.data;
  },

  getMyWishlist: async (): Promise<WishlistItem[]> => {
    const response = await api.get<WishlistItem[]>("/wishlist");
    return response.data;
  },

  removeFromWishlist: async (itemId: number): Promise<void> => {
    await api.delete(`/wishlist/${itemId}`);
  },
};

export default listService;