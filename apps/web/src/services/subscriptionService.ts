import api from "./api";
import { Subscription } from "@/types/subscription";

const subscriptionService = {
  createCheckout: async (): Promise<{ url: string }> => {
    const response = await api.post<{ url: string }>("/subscriptions/checkout");
    return response.data;
  },
  getMySubscription: async (): Promise<Subscription | { active: false }> => {
    const response = await api.get<Subscription | { active: false }>("/subscriptions/me");
    return response.data;
  },
  createPortalSession: async (): Promise<{ url: string }> => {
    const response = await api.post<{ url: string }>("/subscriptions/portal");
    return response.data;
  },
};
export default subscriptionService;
