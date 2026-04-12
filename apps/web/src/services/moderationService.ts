import api from "./api";
import { NewReport } from "@/types/moderation";

const moderationService = {
  reportContent: async (data: NewReport): Promise<void> => {
    await api.post("/moderation/report", data);
  },
};
export default moderationService;
