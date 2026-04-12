import api from "./api";
import { Collection, CollectionWithItems, NewCollection } from "@/types/collection";

const collectionService = {
  create: async (data: NewCollection): Promise<Collection> => {
    const response = await api.post<Collection>("/collections", data);
    return response.data;
  },
  getMyCollections: async (): Promise<Collection[]> => {
    const response = await api.get<Collection[]>("/collections");
    return response.data;
  },
  getById: async (id: number): Promise<CollectionWithItems> => {
    const response = await api.get<CollectionWithItems>(`/collections/${id}`);
    return response.data;
  },
  update: async (id: number, data: Partial<NewCollection>): Promise<Collection> => {
    const response = await api.put<Collection>(`/collections/${id}`, data);
    return response.data;
  },
  delete: async (id: number): Promise<void> => {
    await api.delete(`/collections/${id}`);
  },
  assignItem: async (itemId: number, collectionId: number): Promise<void> => {
    await api.put(`/lists/${itemId}/collection`, { collectionId });
  },
  unassignItem: async (itemId: number): Promise<void> => {
    await api.delete(`/lists/${itemId}/collection`);
  },
};
export default collectionService;
