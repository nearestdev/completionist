"use client";

import { useState, useEffect, useRef, useLayoutEffect } from 'react';
import { PaperPlaneIcon, WifiSlashIcon } from '@phosphor-icons/react/dist/ssr';
import { useAuth } from '@/hooks/useAuth';
import type { DirectMessage, RoomMessage } from '@/types/messaging';

interface ChatBoxProps {
  messages: Array<DirectMessage | RoomMessage>;
  onSendMessage: (content: string) => void;
  onLoadMore?: () => Promise<void>;
  hasMore?: boolean;
  placeholder?: string;
  loading?: boolean;
  disabled?: boolean;
}

export default function ChatBox({ 
  messages, 
  onSendMessage, 
  onLoadMore,
  hasMore = false,
  placeholder = "Type a message...", 
  loading, 
  disabled 
}: ChatBoxProps) {
  const [input, setInput] = useState('');
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const { user } = useAuth();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const prevMessagesLengthRef = useRef(messages.length);
  const beforeLoadScrollHeightRef = useRef<number>(0);

  useLayoutEffect(() => {
    // If message count increased
    if (messages.length > prevMessagesLengthRef.current) {
      if (isFetchingMore && scrollContainerRef.current) {
        // Restore scroll position after loading previous messages
        const newScrollHeight = scrollContainerRef.current.scrollHeight;
        const diff = newScrollHeight - beforeLoadScrollHeightRef.current;
        scrollContainerRef.current.scrollTop = diff;
        setIsFetchingMore(false);
      } else if (!isFetchingMore) {
        // New message received or sent, scroll to bottom
        scrollToBottom();
      }
    }
    prevMessagesLengthRef.current = messages.length;
  }, [messages, isFetchingMore]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleScroll = async (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight } = e.currentTarget;
    if (scrollTop === 0 && hasMore && !loading && !isFetchingMore && onLoadMore) {
       setIsFetchingMore(true);
       beforeLoadScrollHeightRef.current = scrollHeight;
       await onLoadMore();
       // Scroll restoration is handled in useLayoutEffect
    }
  };


  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim() && !disabled) {
      onSendMessage(input.trim());
      setInput('');
    }
  };

  const isDirectMessage = (msg: DirectMessage | RoomMessage): msg is DirectMessage => {
    return 'senderId' in msg;
  };

  return (
    <div className="flex flex-col flex-1 min-h-0 bg-muted/5 relative">
      <div 
        ref={scrollContainerRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto px-6 py-6 space-y-6"
      >
        {isFetchingMore && (
          <div className="flex justify-center p-2">
            <div className="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin" />
          </div>
        )}
        {loading && !isFetchingMore ? (
          <div className="flex items-center justify-center h-full">
            <div className="flex flex-col items-center gap-3">
              <div className="w-8 h-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
              <div className="text-sm text-muted font-medium">Loading messages...</div>
            </div>
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-muted text-center max-w-sm px-4">
              <p className="mb-2 text-lg">👋</p>
              <p>No messages yet. Start the conversation!</p>
            </div>
          </div>
        ) : (
          messages.map((msg, index) => {
            const isSender = isDirectMessage(msg) 
              ? msg.senderId === user?.id
              : msg.userId === user?.id;
            
            const username = isDirectMessage(msg) 
              ? msg.senderUsername
              : msg.username;

             const showUsername = !isSender && (index === 0 || (
               isDirectMessage(messages[index - 1]) 
                 ? (messages[index - 1] as DirectMessage).senderId !== (msg as DirectMessage).senderId
                 : (messages[index - 1] as RoomMessage).userId !== (msg as RoomMessage).userId
             ));

            return (
              <div
                key={isDirectMessage(msg) ? msg.id : `room-${msg.id}`}
                className={`flex w-full ${isSender ? 'justify-end' : 'justify-start'} group`}
              >
                <div className={`max-w-[70%] flex flex-col ${isSender ? 'items-end' : 'items-start'} gap-1`}>
                  {showUsername && (
                    <span className="text-xs font-semibold text-muted-foreground px-1">{username}</span>
                  )}
                  <div
                    className={`
                      px-5 py-3 shadow-sm transition-all relative
                      ${isSender 
                        ? 'bg-gradient-to-br from-primary to-primary-hover text-white rounded-2xl rounded-tr-sm' 
                        : 'bg-card text-foreground border border-border/50 rounded-2xl rounded-tl-sm'
                      }
                    `}
                  >
                    <p className="text-sm break-words leading-relaxed">{msg.content}</p>
                  </div>
                  <span className={`
                    text-[10px] text-muted px-1 opacity-0 group-hover:opacity-100 transition-opacity
                  `}>
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

      <div className="p-4 bg-card/80 backdrop-blur-md border-t border-border z-20">
        <form onSubmit={handleSubmit} className="relative max-w-4xl mx-auto">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={disabled}
            placeholder={disabled ? "Connecting..." : placeholder}
            className="w-full pl-6 pr-14 py-4 rounded-2xl border border-border bg-background focus:bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          />
          <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
            {disabled && (
               <div className="p-2 text-muted animate-pulse">
                <WifiSlashIcon size={20} />
               </div>
            )}
            <button
              type="submit"
              disabled={!input.trim() || disabled}
              className="p-2 bg-primary text-white rounded-xl hover:bg-primary-hover disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md active:scale-95"
            >
              <PaperPlaneIcon size={20} weight="fill" />
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
