import api from "./api";

export interface QueueItem {
  id: number;
  mediaItemId: string;
  queuePosition: number;
  title: string;
  coverImageUrl?: string;
  itemType: string;
}

const queueService = {
  getMyQueue: async (): Promise<QueueItem[]> => {
    const response = await api.get<QueueItem[]>("/queue");
    return response.data;
  },
  addToQueue: async (itemId: number): Promise<void> => {
    await api.post(`/queue/${itemId}`);
  },
  removeFromQueue: async (itemId: number): Promise<void> => {
    await api.delete(`/queue/${itemId}`);
  },
  reorder: async (items: { id: number; position: number }[]): Promise<void> => {
    await api.put("/queue/reorder", { items });
  },
};
export default queueService;
