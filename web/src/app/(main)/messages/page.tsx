"use client";

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { ChatCircleTextIcon } from '@phosphor-icons/react/dist/ssr';
import ConversationList from '@/components/messaging/ConversationList';
import ChatBox from '@/components/messaging/ChatBox';
import { messagingService } from '@/services/messagingService';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { DirectMessage, DMPayload } from '@/types/messaging';

export default function MessagesPage() {
  const [selectedUserId, setSelectedUserId] = useState<number>();
  const [selectedUsername, setSelectedUsername] = useState<string>();
  const [messages, setMessages] = useState<DirectMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const { subscribe } = useWebSocket();
  const searchParams = useSearchParams();

  useEffect(() => {
    const userIdParam = searchParams.get('userId');
    const usernameParam = searchParams.get('username');

    if (userIdParam && usernameParam && !selectedUserId) {
      handleSelectConversation(parseInt(userIdParam), usernameParam);
    } else if (selectedUserId) {
      loadMessages(selectedUserId);
    }
  }, [selectedUserId, searchParams]);

  useEffect(() => {
    const unsubscribe = subscribe<DMPayload>('dm', (payload) => {
      if (
        selectedUserId &&
        (payload.senderId === selectedUserId || payload.receiverId === selectedUserId)
      ) {
        setMessages(prev => {
          if (prev.some(m => m.id === payload.id)) {
            return prev;
          }
          const newMessage: DirectMessage = {
            id: payload.id,
            senderId: payload.senderId,
            senderUsername: payload.senderUsername,
            receiverId: payload.receiverId,
            receiverUsername: payload.receiverUsername,
            content: payload.content,
            isRead: payload.isRead,
            createdAt: payload.createdAt,
          };
          return [...prev, newMessage];
        });
      }
    });

    return unsubscribe;
  }, [subscribe, selectedUserId]);

  const [hasMore, setHasMore] = useState(true);

  const loadMessages = async (userId: number) => {
    try {
      setLoading(true);
      const data = await messagingService.getConversation(userId);
      setMessages(data);
      setHasMore(data.length === 50);
    } catch (error) {
      console.error('Failed to load messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadMoreMessages = async () => {
    if (!selectedUserId || messages.length === 0) return;
    const oldestMessage = messages[0];
    try {
      const moreMessages = await messagingService.getConversation(selectedUserId, 50, oldestMessage.id);
      if (moreMessages.length < 50) {
        setHasMore(false);
      }
      setMessages(prev => [...moreMessages, ...prev]);
    } catch (error) {
      console.error('Failed to load more messages:', error);
    }
  };

  const handleSelectConversation = (userId: number, username: string) => {
    setSelectedUserId(userId);
    setSelectedUsername(username);
    setHasMore(true);
  };

  const handleSendMessage = async (content: string) => {
    if (!selectedUserId) return;

    try {
      await messagingService.sendDirectMessage(selectedUserId, content);
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  return (
    <div className="h-[calc(100vh-80px)] flex max-w-7xl mx-auto px-4 py-6 gap-6 w-full">
      {/* Sidebar - Conversation List */}
      <div className="w-80 lg:w-96 flex flex-col bg-card rounded-3xl border border-border shadow-sm overflow-hidden flex-shrink-0">
        <div className="p-5 border-b border-border bg-muted/5">
          <h2 className="text-xl font-bold text-foreground flex items-center gap-2">
            <ChatCircleTextIcon weight="fill" className="text-primary" />
            Messages
          </h2>
        </div>
        <ConversationList 
          onSelectConversation={handleSelectConversation}
          selectedUserId={selectedUserId}
        />
      </div>

      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col bg-card rounded-3xl border border-border shadow-sm overflow-hidden">
        {selectedUserId && selectedUsername ? (
          <>
            <div className="p-4 border-b border-border bg-background/50 backdrop-blur-sm z-10">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary to-primary-hover flex items-center justify-center text-white font-bold shadow-md">
                  {selectedUsername.charAt(0).toUpperCase()}
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-foreground">{selectedUsername}</h3>
                  <div className="flex items-center gap-1.5">
                    <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
                    <span className="text-xs text-muted">Online</span>
                  </div>
                </div>
              </div>
            </div>
            <ChatBox 
              messages={messages}
              onSendMessage={handleSendMessage}
              onLoadMore={loadMoreMessages}
              hasMore={hasMore}
              loading={loading}
            />
          </>
        ) : (
          <div className="flex flex-col items-center justify-center h-full bg-muted/5">
             <div className="w-24 h-24 rounded-full bg-primary/5 flex items-center justify-center mb-6 text-primary animate-in zoom-in-50 duration-500">
               <ChatCircleTextIcon size={48} weight="thin" />
             </div>
            <div className="text-center max-w-xs">
              <h3 className="text-2xl font-bold text-foreground mb-2">Your Messages</h3>
              <p className="text-muted">
                Select a conversation from the sidebar to start chatting or connect with someone new.
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
