import api from "./api";
import { AdCampaign } from "@/types/ad";

const adService = {
  getActiveAd: async (placement: string = "feed"): Promise<AdCampaign | null> => {
    const response = await api.get<AdCampaign | null>(`/ads?placement=${placement}`);
    return response.data;
  },
  recordImpression: async (adId: number, type: string = "view"): Promise<void> => {
    await api.post(`/ads/${adId}/impression`, { type });
  },
};
export default adService;
