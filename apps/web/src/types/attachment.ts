export interface Attachment {
  id: string;
  kind: string;
  value: string;
  storageProvider?: string;
  contentType?: string;
  sizeBytes?: number;
  createdAt: string;
  updatedAt: string;
}

export interface createAttachmentBody {
  kind: string;
  value: string;
  storageProvider?: string;
  contentType?: string;
  sizeBytes?: number;
}

export interface linkBody {
  entityTable: string;
  entityPk: string;
  attachmentId: string;
}