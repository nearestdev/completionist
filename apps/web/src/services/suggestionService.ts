import api from "./api";
import { MediaItem } from "@/types/list";

const suggestionService = {
  getSuggestions: async (limit: number = 10): Promise<MediaItem[]> => {
    const response = await api.get<MediaItem[]>(`/suggestions?limit=${limit}`);
    return response.data;
  },
};
export default suggestionService;
