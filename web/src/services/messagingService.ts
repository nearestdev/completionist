import api from './api';
import type { 
  DirectMessage, 
  Conversation, 
  SendDirectMessageRequest 
} from '@/types/messaging';

export const messagingService = {
  async sendDirectMessage(receiverId: number, content: string): Promise<DirectMessage> {
    const response = await api.post<DirectMessage>('/messages/send', {
      receiverId,
      content,
    } as SendDirectMessageRequest);
    return response.data;
  },

  async getConversations(): Promise<Conversation[]> {
    const response = await api.get<Conversation[]>('/messages/conversations');
    return response.data;
  },

  async getConversation(userId: number, limit?: number): Promise<DirectMessage[]> {
    const params = new URLSearchParams();
    if (limit) params.append('limit', limit.toString());
    
    const response = await api.get<DirectMessage[]>(
      `/messages/conversation/${userId}?${params.toString()}`
    );
    return response.data;
  },
};
