"use client";

import { useState, useEffect, useRef, useCallback } from 'react';
import { messagingService } from '@/services/messagingService';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { DirectMessage, Conversation, DMPayload } from '@/types/messaging';

interface ConversationListProps {
  onSelectConversation: (userId: number, username: string) => void;
  selectedUserId?: number;
}

export default function ConversationList({ onSelectConversation, selectedUserId }: ConversationListProps) {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const { subscribe } = useWebSocket();

  useEffect(() => {
    loadConversations();
  }, []);

  useEffect(() => {
    const unsubscribe = subscribe<DMPayload>('dm', (payload) => {
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

        const newConv: Conversation = {
          userId: payload.senderId,
          username: payload.senderUsername,
          lastMessage: payload.content,
          lastMessageAt: payload.createdAt,
          unreadCount: selectedUserId !== payload.senderId ? 1 : 0,
        };
        return [newConv, ...prev];
      });
    });

    return unsubscribe;
  }, [subscribe, selectedUserId]);

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
      <div className="flex items-center justify-center h-full">
        <div className="text-muted">Loading conversations...</div>
      </div>
    );
  }

  if (conversations.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-2 px-4 text-center">
        <div className="text-muted">No conversations yet</div>
        <p className="text-sm text-muted/70">
          Start a conversation by visiting a user&apos;s profile
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-y-auto">
      {conversations.map((conv) => (
        <button
          key={conv.userId}
          onClick={() => handleSelectConversation(conv.userId, conv.username)}
          className={`
            flex items-center gap-3 p-4 hover:bg-primary/5 transition-colors border-b border-border
            ${selectedUserId === conv.userId ? 'bg-primary/10' : ''}
          `}
        >
          <div className="w-12 h-12 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold flex-shrink-0">
            {conv.username.charAt(0).toUpperCase()}
          </div>
          
          <div className="flex-1 min-w-0 text-left">
            <div className="flex items-center justify-between gap-2">
              <h4 className="font-semibold text-foreground truncate">
                {conv.username}
              </h4>
              <span className="text-xs text-muted flex-shrink-0">
                {formatTime(conv.lastMessageAt)}
              </span>
            </div>
            <p className="text-sm text-muted truncate">{conv.lastMessage}</p>
          </div>

          {conv.unreadCount > 0 && (
            <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center text-xs text-white font-bold flex-shrink-0">
              {conv.unreadCount}
            </div>
          )}
        </button>
      ))}
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
