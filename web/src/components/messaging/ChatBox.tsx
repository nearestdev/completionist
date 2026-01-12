"use client";

import { useState, useEffect, useRef } from 'react';
import { PaperPlaneIcon } from '@phosphor-icons/react/dist/ssr';
import { useAuth } from '@/hooks/useAuth';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { DirectMessage, RoomMessage } from '@/types/messaging';

interface ChatBoxProps {
  messages: Array<DirectMessage | RoomMessage>;
  onSendMessage: (content: string) => void;
  placeholder?: string;
  loading?: boolean;
}

export default function ChatBox({ messages, onSendMessage, placeholder = "Type a message...", loading }: ChatBoxProps) {
  const [input, setInput] = useState('');
  const { user } = useAuth();
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim()) {
      onSendMessage(input.trim());
      setInput('');
    }
  };

  const isDirectMessage = (msg: DirectMessage | RoomMessage): msg is DirectMessage => {
    return 'senderId' in msg;
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-muted">Loading messages...</div>
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-muted">No messages yet. Start the conversation!</div>
          </div>
        ) : (
          messages.map((msg, index) => {
            const isSender = isDirectMessage(msg) 
              ? msg.senderId === user?.id
              : msg.userId === user?.id;
            
            const username = isDirectMessage(msg) 
              ? msg.senderUsername
              : msg.username;

            return (
              <div
                key={isDirectMessage(msg) ? msg.id : `room-${msg.id}`}
                className={`flex ${isSender ? 'justify-end' : 'justify-start'}`}
              >
                <div className={`max-w-[70%] ${isSender ? 'items-end' : 'items-start'} flex flex-col gap-1`}>
                  {!isSender && (
                    <span className="text-xs text-muted px-2">{username}</span>
                  )}
                  <div
                    className={`
                      px-4 py-2 rounded-2xl
                      ${isSender 
                        ? 'bg-primary text-white rounded-br-sm' 
                        : 'bg-card border border-border rounded-bl-sm'
                      }
                    `}
                  >
                    <p className="text-sm break-words">{msg.content}</p>
                  </div>
                  <span className="text-xs text-muted px-2">
                    {new Date(msg.createdAt).toLocaleTimeString('en-US', {
                      hour: 'numeric',
                      minute: '2-digit'
                    })}
                  </span>
                </div>
              </div>
            );
          })
        )}
        <div ref={messagesEndRef} />
      </div>

      <form onSubmit={handleSubmit} className="border-t border-border p-4">
        <div className="flex gap-2">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={placeholder}
            className="flex-1 px-4 py-2 rounded-full border border-border bg-background focus:outline-none focus:ring-2 focus:ring-primary"
          />
          <button
            type="submit"
            disabled={!input.trim()}
            className="px-6 py-2 bg-primary text-white rounded-full hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2"
          >
            <PaperPlaneIcon size={20} weight="fill" />
          </button>
        </div>
      </form>
    </div>
  );
}
