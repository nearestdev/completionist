import api from "./api";
import { CompletionReview, ReviewWithUser, NewReview } from "@/types/review";

const reviewService = {
  create: async (data: NewReview): Promise<CompletionReview> => {
    const response = await api.post<CompletionReview>("/reviews", data);
    return response.data;
  },
  getUserReviews: async (userId: number, limit: number = 20, offset: number = 0): Promise<ReviewWithUser[]> => {
    const response = await api.get<ReviewWithUser[]>(`/reviews/user/${userId}?limit=${limit}&offset=${offset}`);
    return response.data;
  },
  getMediaReviews: async (mediaId: string, limit: number = 20, offset: number = 0): Promise<ReviewWithUser[]> => {
    const response = await api.get<ReviewWithUser[]>(`/reviews/media/${mediaId}?limit=${limit}&offset=${offset}`);
    return response.data;
  },
  update: async (id: number, data: Partial<NewReview>): Promise<CompletionReview> => {
    const response = await api.put<CompletionReview>(`/reviews/${id}`, data);
    return response.data;
  },
  delete: async (id: number): Promise<void> => {
    await api.delete(`/reviews/${id}`);
  },
};
export default reviewService;
