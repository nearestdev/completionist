import api from "./api";
import {
  Attachment,
  createAttachmentBody,
  linkBody,
} from "@/types/attachment";

const attachmentService = {
  createAttachment: async (data: createAttachmentBody): Promise<Attachment> => {
    const response = await api.post<Attachment>("/attachments", data);
    return response.data;
  },

  linkAttachment: async (data: linkBody): Promise<any> => {
    const response = await api.post<any>("/attachments/link", data);
    return response.data;
  },

  unlinkAttachment: async (data: linkBody): Promise<void> => {
    await api.delete("/attachments/unlink", { data });
  },

  listAttachmentsByEntity: async (table: string, pk: string): Promise<Attachment[]> => {
    const response = await api.get<Attachment[]>(`/attachments/by-entity?table=${table}&pk=${pk}`);
    return response.data;
  },
};

export default attachmentService;