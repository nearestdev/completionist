"use client";

import { useState, useEffect } from 'react';
import { messagingService } from '@/services/messagingService';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useAuth } from '@/hooks/useAuth';
import type { Conversation, DMPayload } from '@/types/messaging';

interface ConversationListProps {
  onSelectConversation: (userId: number, username: string) => void;
  selectedUserId?: number;
}

export default function ConversationList({ onSelectConversation, selectedUserId }: ConversationListProps) {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const { subscribe } = useWebSocket();
  const { user } = useAuth();

  useEffect(() => {
    loadConversations();
  }, []);

  useEffect(() => {
    const unsubscribe = subscribe<DMPayload>('dm', (payload) => {
      if (!user) return;
      
      setConversations(prev => {
        const existing = prev.find(c => 
          c.userId === payload.senderId || c.userId === payload.receiverId
        );

        if (existing) {
          return prev.map(c => {
            if (c.userId === payload.senderId || c.userId === payload.receiverId) {
              const isFromOther = c.userId === payload.senderId;
              return {
                ...c,
                lastMessage: payload.content,
                lastMessageAt: payload.createdAt,
                unreadCount: isFromOther && selectedUserId !== payload.senderId 
                  ? c.unreadCount + 1 
                  : c.unreadCount,
              };
            }
            return c;
          }).sort((a, b) => 
            new Date(b.lastMessageAt).getTime() - new Date(a.lastMessageAt).getTime()
          );
        }

        const isMe = user.id === payload.senderId;
        const otherUserId = isMe ? payload.receiverId : payload.senderId;
        const otherUsername = isMe ? payload.receiverUsername : payload.senderUsername;

        const newConv: Conversation = {
          userId: otherUserId,
          username: otherUsername,
          lastMessage: payload.content,
          lastMessageAt: payload.createdAt,
          unreadCount: selectedUserId !== payload.senderId && !isMe ? 1 : 0,
        };
        return [newConv, ...prev];
      });
    });

    return unsubscribe;
  }, [subscribe, selectedUserId, user]);



  const loadConversations = async () => {
    try {
      const data = await messagingService.getConversations();
      setConversations(data);
    } catch (error) {
      console.error('Failed to load conversations:', error);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = (userId: number) => {
    setConversations(prev =>
      prev.map(c => c.userId === userId ? { ...c, unreadCount: 0 } : c)
    );
  };

  const handleSelectConversation = (userId: number, username: string) => {
    markAsRead(userId);
    onSelectConversation(userId, username);
  };

  if (loading) {
    return (
      <div className="p-4 space-y-3">
        {[1, 2, 3, 4].map(i => (
          <div key={i} className="flex gap-3 items-center">
            <div className="w-12 h-12 rounded-full bg-muted/10 animate-pulse" />
            <div className="flex-1 space-y-2">
              <div className="h-4 w-1/3 bg-muted/10 animate-pulse rounded" />
              <div className="h-3 w-3/4 bg-muted/10 animate-pulse rounded" />
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (conversations.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
        <div className="text-muted mb-2">No conversations yet</div>
        <p className="text-sm text-muted/70">
          Start a conversation by visiting a user&apos;s profile
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-y-auto p-2">
      {conversations.map((conv) => {
        const isSelected = selectedUserId === conv.userId;
        const hasUnread = conv.unreadCount > 0;
        
        return (
          <button
            key={conv.userId}
            onClick={() => handleSelectConversation(conv.userId, conv.username)}
            className={`
              group w-full flex items-center gap-3 p-3 rounded-xl transition-all mb-1
              ${isSelected 
                ? 'bg-primary/10 border-primary/20 shadow-sm' 
                : 'hover:bg-muted/10 border border-transparent'
              }
            `}
          >
            <div className={`
              w-12 h-12 rounded-full flex items-center justify-center text-lg font-bold shadow-sm transition-all
              ${isSelected ? 'bg-primary text-white scale-105' : 'bg-muted/20 text-muted-foreground group-hover:bg-muted/30'}
            `}>
              {conv.username.charAt(0).toUpperCase()}
            </div>
            
            <div className="flex-1 min-w-0 text-left">
              <div className="flex items-center justify-between gap-2 mb-0.5">
                <h4 className={`text-sm font-semibold truncate ${isSelected ? 'text-primary' : 'text-foreground'}`}>
                  {conv.username}
                </h4>
                <span className="text-[10px] text-muted flex-shrink-0">
                  {formatTime(conv.lastMessageAt)}
                </span>
              </div>
              <div className="flex items-center justify-between gap-2">
                <p className={`text-xs truncate ${hasUnread ? 'text-foreground font-semibold' : 'text-muted'}`}>
                  {(() => {
                      const imageMatch = conv.lastMessage.match(/^\[(.*?)\]\((.*?)\)$/);
                      const isImage = imageMatch && (/\.(jpg|jpeg|png|gif|webp)$/i.test(imageMatch[2]) || /\.(jpg|jpeg|png|gif|webp)\?/i.test(imageMatch[2]));
                      return isImage ? '📷 Sent an image' : conv.lastMessage;
                  })()}
                </p>
                {hasUnread && (
                  <div className="w-5 h-5 rounded-full bg-primary flex items-center justify-center text-[10px] text-white font-bold flex-shrink-0 shadow-sm animate-in scale-in duration-200">
                    {conv.unreadCount}
                  </div>
                )}
              </div>
            </div>
          </button>
        );
      })}
    </div>
  );
}

function formatTime(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const hours = diff / (1000 * 60 * 60);

  if (hours < 24) {
    return date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit' });
  }
  if (hours < 48) {
    return 'Yesterday';
  }
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}
