"use client";

import { useState, useEffect } from 'react';
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

  useEffect(() => {
    if (selectedUserId) {
      loadMessages(selectedUserId);
    }
  }, [selectedUserId]);

  useEffect(() => {
    const unsubscribe = subscribe<DMPayload>('dm', (payload) => {
      if (
        selectedUserId &&
        (payload.senderId === selectedUserId || payload.receiverId === selectedUserId)
      ) {
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
        setMessages(prev => [...prev, newMessage]);
      }
    });

    return unsubscribe;
  }, [subscribe, selectedUserId]);

  const loadMessages = async (userId: number) => {
    try {
      setLoading(true);
      const data = await messagingService.getConversation(userId);
      setMessages(data);
    } catch (error) {
      console.error('Failed to load messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectConversation = (userId: number, username: string) => {
    setSelectedUserId(userId);
    setSelectedUsername(username);
  };

  const handleSendMessage = async (content: string) => {
    if (!selectedUserId) return;

    try {
      const newMessage = await messagingService.sendDirectMessage(selectedUserId, content);
      setMessages(prev => [...prev, newMessage]);
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  return (
    <div className="h-[calc(100vh-80px)] flex">
      <div className="w-80 border-r border-border flex flex-col">
        <div className="p-4 border-b border-border">
          <h2 className="text-xl font-bold text-foreground">Messages</h2>
        </div>
        <ConversationList 
          onSelectConversation={handleSelectConversation}
          selectedUserId={selectedUserId}
        />
      </div>

      <div className="flex-1 flex flex-col">
        {selectedUserId && selectedUsername ? (
          <>
            <div className="p-4 border-b border-border">
              <h3 className="text-lg font-semibold text-foreground">{selectedUsername}</h3>
            </div>
            <ChatBox 
              messages={messages}
              onSendMessage={handleSendMessage}
              loading={loading}
            />
          </>
        ) : (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <h3 className="text-xl font-semibold text-muted mb-2">Select a conversation</h3>
              <p className="text-sm text-muted/70">
                Choose a conversation from the list to start messaging
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
