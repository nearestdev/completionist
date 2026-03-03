import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {  SignOutIcon, SidebarIcon } from '@phosphor-icons/react/dist/ssr';
import ChatBox from '@/components/messaging/ChatBox';
import MembersList from '@/components/rooms/MembersList';
import { roomService } from '@/services/roomService';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useAuth } from '@/hooks/useAuth';
import type { Room, RoomMessage, RoomChatPayload, RoomJoinPayload, RoomLeavePayload, RoomMember, RoomUsersPayload } from '@/types/messaging';

interface ChatRoomProps {
  room: Room;
}

export default function ChatRoom({ room }: ChatRoomProps) {
  const router = useRouter();
  const { user } = useAuth();
  const [messages, setMessages] = useState<RoomMessage[]>([]);
  const [members, setMembers] = useState<RoomMember[]>([]);
  const [onlineUsers, setOnlineUsers] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(true);
  const [showSidebar, setShowSidebar] = useState(true);
  const { subscribe, sendMessage, isConnected } = useWebSocket();

  useEffect(() => {
    loadData();
  }, [room.id]);

  useEffect(() => {
    const unsubscribeChat = subscribe<RoomChatPayload>('room_chat', (payload) => {
          if (payload.roomId === room.id) {
        const newMessage: RoomMessage = {
          id: payload.id,
          roomId: payload.roomId,
          userId: payload.userId,
          username: payload.username,
          content: payload.content,
          replyToId: payload.replyToId,
          isPinned: payload.isPinned,
          createdAt: payload.createdAt,
        };
        setMessages(prev => [...prev, newMessage]);
      }
    });

    const unsubscribeJoin = subscribe<RoomJoinPayload>('room_join', (payload) => {
      if (payload.roomId === room.id) {
        setOnlineUsers(prev => {
            const next = new Set(prev);
            next.add(payload.userId);
            return next;
        });
        setMembers(prev => {
            if (prev.some(m => m.userId === payload.userId)) return prev;
            return [...prev, {
                id: 0,
                roomId: payload.roomId,
                userId: payload.userId,
                username: payload.username,
                isModerator: false,
                joinedAt: new Date().toISOString()
            }];
        });
      }
    });

    const unsubscribeLeave = subscribe<RoomLeavePayload>('room_leave', (payload) => {
      if (payload.roomId === room.id) {
        setMembers(prev => prev.filter(m => m.userId !== payload.userId));
        setOnlineUsers(prev => {
            const next = new Set(prev);
            next.delete(payload.userId);
            return next;
        });
      }
    });

    const unsubscribeUsers = subscribe<RoomUsersPayload>('room_users', (payload) => {
        if (payload.roomId === room.id) {
            setOnlineUsers(new Set(payload.users));
        }
    });

    return () => {
      unsubscribeChat();
      unsubscribeJoin();
      unsubscribeLeave();
      unsubscribeUsers();
    };
  }, [subscribe, room.id]);

  const loadData = async () => {
    try {
      setLoading(true);
      const [msgs, mems] = await Promise.all([
        roomService.getRoomMessages(room.id),
        roomService.getRoomMembers(room.id)
      ]);
      setMessages(msgs.reverse());
      setMembers(mems);
    } catch (error: any) {
      console.error('Failed to load room data:', error);
      if (error?.response?.status === 403) {
        console.error('Not a member/forbidden');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSendMessage = (content: string) => {
    sendMessage('room_chat', {
      roomId: room.id,
      content,
    });
  };

  const handleLeaveRoom = async () => {
    try {
      await roomService.leaveRoom(room.id);
      router.push('/rooms');
    } catch (error) {
      console.error('Failed to leave room:', error);
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="p-4 border-b border-border flex items-center justify-between bg-card/30">
        <div>
          <h2 className="text-xl font-bold text-foreground flex items-center gap-2">
            {room.name}
            {room.isPublic ? null : <span className="text-xs bg-yellow-500/10 text-yellow-500 px-2 py-0.5 rounded-full">Private</span>}
          </h2>
          {room.description && (
            <p className="text-sm text-muted line-clamp-1">{room.description}</p>
          )}
        </div>
        <div className="flex items-center gap-4">
          <button
            onClick={() => setShowSidebar(!showSidebar)}
            className={`p-2 rounded-lg transition-colors ${showSidebar ? 'bg-primary/10 text-primary' : 'hover:bg-accent text-muted-foreground'}`}
            title="Toggle Members List"
          >
            <SidebarIcon size={20} />
          </button>
          
          <div className="h-6 w-px bg-border hidden sm:block"></div>

          <button
            onClick={handleLeaveRoom}
            className="flex items-center gap-2 px-3 py-1.5 text-sm text-red-500 hover:bg-red-500/10 rounded-lg transition-all"
          >
            <SignOutIcon size={18} />
            <span className="hidden sm:inline">Leave</span>
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-hidden flex">
        <div className="flex-1 flex flex-col min-w-0">
          <ChatBox
            messages={messages}
            onSendMessage={handleSendMessage}
            loading={loading}
            disabled={!isConnected}
            placeholder="Send a message..."
          />
        </div>
        
        {/* Sidebar */}
        {showSidebar && (
          <div className="w-64 border-l border-border hidden md:block transition-all animate-in slide-in-from-right-4 duration-200">
            <MembersList members={members} onlineUsers={onlineUsers} />
          </div>
        )}
      </div>
    </div>
  );
}
