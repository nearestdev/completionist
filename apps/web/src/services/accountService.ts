import api from "./api";

const accountService = {
  requestDeletion: async (reason?: string): Promise<void> => {
    await api.post("/me/delete", { reason });
  },
  cancelDeletion: async (): Promise<void> => {
    await api.delete("/me/delete");
  },
  exportData: async (): Promise<Blob> => {
    const response = await api.get("/me/export", { responseType: "blob" });
    return response.data;
  },
  submitAppeal: async (appealText: string): Promise<void> => {
    await api.post("/me/appeal", { appealText });
  },
};
export default accountService;
