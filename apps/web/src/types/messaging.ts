export type RoomType = 'chat' | 'music' | 'watch';

export type MessageType = 
  | 'dm'
  | 'room_chat'
  | 'room_join'
  | 'room_leave'
  | 'room_state'
  | 'room_users'
  | 'typing'
  | 'ping'
  | 'pong'
  | 'reaction'
  | 'pin'
  | 'error';

export interface RoomUsersPayload {
  roomId: number;
  users: number[];
}

export interface Attachment {
  id: number;
  filePath: string;
  fileType: string;
  fileName: string;
  fileSize: number;
  createdAt: string;
}

export interface Reaction {
  id: number;
  messageId: number;
  userId: number;
  username: string;
  reaction: string;
}

export interface DirectMessage {
  id: number;
  senderId: number;
  senderUsername: string;
  receiverId: number;
  receiverUsername: string;
  content: string;
  isRead: boolean;
  replyToId?: number;
  isPinned: boolean;
  reactions?: Reaction[];
  attachments?: Attachment[];
  createdAt: string;
}

export interface Conversation {
  userId: number;
  username: string;
  lastMessage: string;
  lastMessageAt: string;
  unreadCount: number;
}

export interface Room {
  id: number;
  name: string;
  description?: string;
  roomType: RoomType;
  creatorId: number;
  isPublic: boolean;
  maxMembers: number;
  memberCount?: number;
  isMember?: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface RoomMember {
  id: number;
  roomId: number;
  userId: number;
  username: string;
  isModerator: boolean;
  joinedAt: string;
}

export interface RoomMessage {
  id: number;
  roomId: number;
  userId: number;
  username: string;
  content: string;
  replyToId?: number;
  isPinned: boolean;
  reactions?: Reaction[];
  attachments?: Attachment[];
  createdAt: string;
}

export interface RoomState {
  roomId: number;
  currentMediaUrl?: string;
  currentMediaTitle?: string;
  currentPositionMs: number;
  isPlaying: boolean;
  updatedBy: number;
  updatedByUsername: string;
  updatedAt: string;
}

export interface WSMessage<T = unknown> {
  type: MessageType;
  payload: T;
}

export interface DMPayload {
  id: number;
  senderId: number;
  senderUsername: string;
  receiverId: number;
  receiverUsername: string;
  content: string;
  isRead: boolean;
  replyToId?: number;
  isPinned: boolean;
  reactions?: Reaction[]; // Optional in payload, usually fetched via REST or separate event
  attachments?: Attachment[]; // Same
  createdAt: string;
}

export interface RoomChatPayload {
  id: number;
  roomId: number;
  userId: number;
  username: string;
  content: string;
  replyToId?: number;
  isPinned: boolean;
  createdAt: string;
}

export interface RoomJoinPayload {
  roomId: number;
  userId: number;
  username: string;
}

export interface RoomLeavePayload {
  roomId: number;
  userId: number;
  username: string;
}

export interface RoomStatePayload {
  roomId: number;
  currentMediaUrl?: string;
  currentMediaTitle?: string;
  currentPositionMs: number;
  isPlaying: boolean;
  updatedBy: number;
  updatedByUsername: string;
  updatedAt: string;
}

export interface TypingPayload {
  roomId?: number;
  userId: number;
  username: string;
}

export interface ErrorPayload {
  message: string;
}

export interface CreateRoomRequest {
  name: string;
  description?: string;
  roomType: RoomType;
  isPublic: boolean;
  maxMembers?: number;
}

export interface SendDirectMessageRequest {
  receiverId: number;
  content: string;
  replyToId?: number;
}

export interface UpdateRoomStateRequest {
  currentMediaUrl?: string;
  currentMediaTitle?: string;
  currentPositionMs?: number;
  isPlaying?: boolean;
}
