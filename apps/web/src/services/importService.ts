import api from "./api";
import { ConnectedAccount, ImportJob } from "@/types/connected_account";

const importService = {
  getConnections: async (): Promise<ConnectedAccount[]> => {
    const response = await api.get<ConnectedAccount[]>("/connections");
    return response.data;
  },
  connectProvider: async (provider: string): Promise<{ url: string }> => {
    const response = await api.post<{ url: string }>(`/connections/${provider}/connect`);
    return response.data;
  },
  disconnect: async (provider: string): Promise<void> => {
    await api.delete(`/connections/${provider}`);
  },
  triggerSync: async (provider: string): Promise<ImportJob> => {
    const response = await api.post<ImportJob>(`/connections/${provider}/sync`);
    return response.data;
  },
  getImportJobs: async (): Promise<ImportJob[]> => {
    const response = await api.get<ImportJob[]>("/imports");
    return response.data;
  },
};
export default importService;
