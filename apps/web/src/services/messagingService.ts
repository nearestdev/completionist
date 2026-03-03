import api from './api';
import type { 
  DirectMessage, 
  Conversation, 
  SendDirectMessageRequest 
} from '@/types/messaging';

export const messagingService = {
  async sendDirectMessage(receiverId: number, content: string, replyToId?: number): Promise<DirectMessage> {
    const response = await api.post<DirectMessage>('/messages/send', {
      receiverId,
      content,
      replyToId,
    } as SendDirectMessageRequest);
    return response.data;
  },

  async getConversations(): Promise<Conversation[]> {
    const response = await api.get<any[]>('/messages/conversations');
    return response.data.map(c => ({
      userId: c.otherUserId,
      username: c.otherUsername,
      lastMessage: c.lastMessage,
      lastMessageAt: c.lastMessageTime,
      unreadCount: c.unreadCount,
    }));
  },

  async getConversation(userId: number, limit?: number, cursor?: number): Promise<DirectMessage[]> {
    const params = new URLSearchParams();
    if (limit) params.append('limit', limit.toString());
    if (cursor) params.append('cursor', cursor.toString());
    
    const response = await api.get<DirectMessage[]>(
      `/messages/conversation/${userId}?${params.toString()}`
    );
    return response.data.reverse();
  },

  async reactToMessage(messageId: number, reaction: string): Promise<void> {
    await api.post(`/messages/${messageId}/react`, { reaction });
  },

  async pinMessage(messageId: number): Promise<void> {
    await api.post(`/messages/${messageId}/pin`);
  },

  async uploadFile(file: File): Promise<string> {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post<{ url: string }>('/attachments/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data.url;
  },
};
