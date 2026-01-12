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
};
