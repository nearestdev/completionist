import api from './api';
import type { 
  Room, 
  RoomMessage, 
  RoomState, 
  CreateRoomRequest,
  UpdateRoomStateRequest,
  RoomType,
  RoomMember
} from '@/types/messaging';

export const roomService = {
  async createRoom(data: CreateRoomRequest): Promise<Room> {
    const response = await api.post<Room>('/rooms', data);
    return response.data;
  },

  async getRooms(type?: RoomType, limit?: number): Promise<Room[]> {
    const params = new URLSearchParams();
    if (type) params.append('type', type);
    if (limit) params.append('limit', limit.toString());
    
    const response = await api.get<Room[]>(`/rooms?${params.toString()}`);
    return response.data;
  },

  async getRoom(roomId: number): Promise<Room> {
    const response = await api.get<Room>(`/rooms/${roomId}`);
    return response.data;
  },

  async joinRoom(roomId: number): Promise<void> {
    await api.post(`/rooms/${roomId}/join`);
  },

  async leaveRoom(roomId: number): Promise<void> {
    await api.delete(`/rooms/${roomId}/leave`);
  },

  async getRoomMessages(roomId: number, limit?: number): Promise<RoomMessage[]> {
    const params = new URLSearchParams();
    if (limit) params.append('limit', limit.toString());
    
    const response = await api.get<RoomMessage[]>(
      `/rooms/${roomId}/messages?${params.toString()}`
    );
    return response.data;
  },

  async getRoomState(roomId: number): Promise<RoomState> {
    const response = await api.get<RoomState>(`/rooms/${roomId}/state`);
    return response.data;
  },

  async updateRoomState(roomId: number, state: UpdateRoomStateRequest): Promise<RoomState> {
    const response = await api.patch<RoomState>(`/rooms/${roomId}/state`, state);
    return response.data;
  },

  async getRoomMembers(roomId: number): Promise<RoomMember[]> {
    const response = await api.get<RoomMember[]>(`/rooms/${roomId}/members`);
    return response.data;
  },

  async reactToMessage(roomId: number, messageId: number, reaction: string): Promise<void> {
    await api.post(`/rooms/${roomId}/messages/${messageId}/react`, { reaction });
  },

  async pinMessage(roomId: number, messageId: number): Promise<void> {
    await api.post(`/rooms/${roomId}/messages/${messageId}/pin`);
  },

  async getPinnedMessages(roomId: number): Promise<RoomMessage[]> {
    const response = await api.get<RoomMessage[]>(`/rooms/${roomId}/pinned`);
    return response.data;
  },
};
