"use client";

import { useState, useEffect, useRef, useLayoutEffect } from 'react';
import { 
  PaperPlaneIcon, 
  WifiSlashIcon, 
  PaperclipIcon, 
  XIcon, 
  PushPinIcon, 
  SmileyIcon, 
  ArrowUUpLeftIcon
} from '@phosphor-icons/react/dist/ssr';
import { useAuth } from '@/hooks/useAuth';
import type { DirectMessage, RoomMessage, Reaction, Attachment } from '@/types/messaging';
import { messagingService } from '@/services/messagingService';
import { roomService } from '@/services/roomService';
import Image from 'next/image';

interface ChatBoxProps {
  messages: Array<DirectMessage | RoomMessage>;
  onSendMessage: (content: string, replyToId?: number) => void;
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
  const [replyTo, setReplyTo] = useState<DirectMessage | RoomMessage | null>(null);
  const [highlightedMessageId, setHighlightedMessageId] = useState<number | null>(null);
  const [emojiPickerMessageId, setEmojiPickerMessageId] = useState<number | null>(null);
  const { user } = useAuth();
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const prevMessagesLengthRef = useRef(messages.length);
  const beforeLoadScrollHeightRef = useRef<number>(0);
  const shouldRestoreScrollRef = useRef(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const messageRefs = useRef<Map<number, HTMLDivElement>>(new Map());

  function scrollToBottom() {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }

  function scrollToMessage(messageId: number) {
    const messageElement = messageRefs.current.get(messageId);
    if (messageElement && scrollContainerRef.current) {
      messageElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
      setHighlightedMessageId(messageId);
      setTimeout(() => setHighlightedMessageId(null), 2000);
    }
  }

  useLayoutEffect(() => {
    if (messages.length > prevMessagesLengthRef.current) {
      if (shouldRestoreScrollRef.current && scrollContainerRef.current) {
        const newScrollHeight = scrollContainerRef.current.scrollHeight;
        const diff = newScrollHeight - beforeLoadScrollHeightRef.current;
        scrollContainerRef.current.scrollTop = diff;
        shouldRestoreScrollRef.current = false;
      } else if (!isFetchingMore) {
        scrollToBottom();
      }
    }
    prevMessagesLengthRef.current = messages.length;
  }, [messages, isFetchingMore]);

  const getReplyMessageContent = (replyToId: number): string => {
    const replyMsg = messages.find(m => m.id === replyToId);
    if (!replyMsg) return "Message not found";
    const imageMatch = replyMsg.content.match(/^\[(.*?)\]\((.*?)\)$/);
    const isImage = imageMatch && (/\.(jpg|jpeg|png|gif|webp)$/i.test(imageMatch[2]) || /\.(jpg|jpeg|png|gif|webp)\?/i.test(imageMatch[2]));
    if (isImage) return "📷 Image";
    return replyMsg.content.length > 100 ? replyMsg.content.substring(0, 100) + "..." : replyMsg.content;
  };

  const getReplyMessageSender = (msg: DirectMessage | RoomMessage): string => {
    if (!msg.replyToId) return "";
    const replyMsg = messages.find(m => m.id === msg.replyToId);
    if (!replyMsg) return "";
    return isDirectMessage(replyMsg) ? replyMsg.senderUsername : replyMsg.username;
  };

  const handleScroll = async (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight } = e.currentTarget;
    if (scrollTop === 0 && hasMore && !loading && !isFetchingMore && onLoadMore) {
       setIsFetchingMore(true);
       beforeLoadScrollHeightRef.current = scrollHeight;
       shouldRestoreScrollRef.current = true;
       try {
         await onLoadMore();
       } finally {
         setIsFetchingMore(false);
       }
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (input.trim() && !disabled) {
      onSendMessage(input.trim(), replyTo?.id);
      setInput('');
      setReplyTo(null);
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      try {
        const url = await messagingService.uploadFile(file);
        const content = `[${file.name}](${url})`;
        onSendMessage(content, replyTo?.id);
        setReplyTo(null);
      } catch (error) {
        console.error("File upload failed", error);
      }
    }
  };

  const isDirectMessage = (msg: DirectMessage | RoomMessage): msg is DirectMessage => {
    return 'senderId' in msg;
  };

  const getMessageId = (msg: DirectMessage | RoomMessage) => msg.id;
  const getSenderId = (msg: DirectMessage | RoomMessage) => isDirectMessage(msg) ? msg.senderId : msg.userId;

  const handleReact = async (msg: DirectMessage | RoomMessage, reaction: string) => {
      try {
          if (isDirectMessage(msg)) {
              await messagingService.reactToMessage(msg.id, reaction);
          } else {
              await roomService.reactToMessage((msg as RoomMessage).roomId, msg.id, reaction);
          }
      } catch (error) {
          console.error("Failed to react", error);
      }
  };

  const handlePin = async (msg: DirectMessage | RoomMessage) => {
      try {
          if (isDirectMessage(msg)) {
              await messagingService.pinMessage(msg.id);
          } else {
              await roomService.pinMessage((msg as RoomMessage).roomId, msg.id);
          }
      } catch (error) {
          console.error("Failed to pin", error);
      }
  };

  return (
    <div className="flex flex-col flex-1 min-h-0 bg-muted/5 relative">
      {(() => {
        const pinnedMessages = messages.filter(m => m.isPinned);
        if (pinnedMessages.length === 0) return null;
        
        return (
          <div className="bg-yellow-500/10 border-b border-yellow-500/30 px-6 py-3">
            <div className="flex items-center gap-2 mb-2">
              <PushPinIcon size={16} className="text-yellow-600" weight="fill" />
              <span className="text-sm font-semibold text-yellow-700">Pinned Messages ({pinnedMessages.length})</span>
            </div>
            <div className="space-y-2 max-h-32 overflow-y-auto">
              {pinnedMessages.map(msg => {
                const username = isDirectMessage(msg) ? msg.senderUsername : msg.username;
                const imageMatch = msg.content.match(/^\[(.*?)\]\((.*?)\)$/);
                const isImage = imageMatch && (/\.(jpg|jpeg|png|gif|webp)$/i.test(imageMatch[2]) || /\.(jpg|jpeg|png|gif|webp)\?/i.test(imageMatch[2]));
                const displayContent = isImage ? "📷 Image" : (msg.content.length > 60 ? msg.content.substring(0, 60) + "..." : msg.content);
                
                return (
                  <button
                    key={msg.id}
                    onClick={() => scrollToMessage(msg.id)}
                    className="flex items-start gap-2 text-left p-2 rounded-lg hover:bg-yellow-500/20 transition-colors w-full"
                  >
                    <span className="text-xs font-semibold text-yellow-700 flex-shrink-0">{username}:</span>
                    <span className="text-xs text-muted-foreground truncate">{displayContent}</span>
                  </button>
                );
              })}
            </div>
          </div>
        );
      })()}
      
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
              <div className="text-sm text-muted font-medium">Loading...</div>
            </div>
          </div>
        ) : messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-muted text-center max-w-sm px-4">
              <p className="mb-2 text-lg">👋</p>
              <p>No messages yet.</p>
            </div>
          </div>
        ) : (
          messages.map((msg, index) => {
            const isSender = getSenderId(msg) === user?.id;
            const username = isDirectMessage(msg) ? msg.senderUsername : msg.username;
            const showUsername = !isSender && (index === 0 || getSenderId(messages[index - 1]) !== getSenderId(msg));
            const isReplyingToHimself = msg.replyToId && (() => {
                const replyMsg = messages.find(m => m.id === msg.replyToId);
                return replyMsg && getSenderId(replyMsg) === user?.id;
            })();

            return (
              <div
                key={isDirectMessage(msg) ? `dm-${msg.id}` : `room-${msg.id}`}
                ref={(el) => {
                  if (el) messageRefs.current.set(msg.id, el);
                  else messageRefs.current.delete(msg.id);
                }}
                className={`flex w-full ${isSender ? 'justify-end' : 'justify-start'} group transition-all ${
                  highlightedMessageId === msg.id ? 'scale-[1.02]' : ''
                }`}
              >
                <div className={`max-w-[70%] flex flex-col ${isSender ? 'items-end' : 'items-start'} gap-1`}>
                  {showUsername && (
                    <span className="text-xs font-semibold text-muted-foreground px-1">{username}</span>
                  )}
                  
                  {msg.replyToId && (
                      <button 
                        onClick={() => scrollToMessage(msg.replyToId!)}
                        className="text-xs bg-muted/30 hover:bg-muted/50 px-3 py-1.5 rounded-lg mb-1 flex items-center gap-2 transition-all border border-border/30 max-w-full group/reply"
                      >
                          <ArrowUUpLeftIcon size={12} className="flex-shrink-0" />
                          <div className="flex flex-col items-start min-w-0">
                            <span className="font-semibold text-[10px] text-primary">{getReplyMessageSender(msg)}</span>
                            <span className="text-muted-foreground truncate max-w-full">{getReplyMessageContent(msg.replyToId)}</span>
                          </div>
                      </button>
                  )}

                  <div className="relative group/bubble">
                      <div
                        className={`
                        px-5 py-3 shadow-sm transition-all relative
                        ${isSender 
                            ? 'bg-gradient-to-br from-primary to-primary-hover text-white rounded-2xl rounded-tr-sm' 
                            : 'bg-card text-foreground border border-border/50 rounded-2xl rounded-tl-sm'
                        }
                        ${msg.isPinned ? 'border-yellow-500 border-2' : ''}
                        ${highlightedMessageId === msg.id ? 'ring-2 ring-primary ring-offset-2' : ''}
                        `}
                    >
                        {msg.isPinned && <PushPinIcon size={12} className="absolute -top-2 -right-2 bg-yellow-500 text-white rounded-full p-0.5" />}
                        {(() => {
                            const imageMatch = msg.content.match(/^\[(.*?)\]\((.*?)\)$/);
                            if (imageMatch) {
                                const [_, alt, url] = imageMatch;
                                const isImage = /\.(jpg|jpeg|png|gif|webp)$/i.test(url) || /\.(jpg|jpeg|png|gif|webp)\?/i.test(url);
                                if (isImage) {
                                  return (
                                    <div className="relative group/image">
                                        <Image 
                                            src={url} 
                                            alt={alt} 
                                            width={300} 
                                            height={200} 
                                            className="rounded-lg max-w-full h-auto object-cover cursor-pointer hover:opacity-90 transition-opacity"
                                            unoptimized
                                        />
                                        <a href={url} target="_blank" rel="noopener noreferrer" className="absolute bottom-2 right-2 bg-black/50 text-white p-1 rounded opacity-0 group-hover/image:opacity-100 transition-opacity text-xs hover:bg-black/70">
                                            Open
                                        </a>
                                    </div>
                                  );
                                }
                            }
                            return <p className="text-sm break-words leading-relaxed whitespace-pre-wrap">{msg.content}</p>;
                        })()}
                    </div>
                    
                    <div className={`absolute ${isSender ? '-left-20' : '-right-20'} top-1/2 -translate-y-1/2 opacity-0 group-hover/bubble:opacity-100 flex gap-1 bg-background/80 backdrop-blur rounded-lg p-1 shadow-sm transition-opacity`}>
                        <button onClick={() => setReplyTo(msg)} title="Reply" className="p-1 hover:bg-muted rounded"><ArrowUUpLeftIcon size={14} /></button>
                        <button onClick={() => handlePin(msg)} title={msg.isPinned ? "Unpin" : "Pin"} className="p-1 hover:bg-muted rounded"><PushPinIcon size={14} /></button>
                        <div className="relative">
                          <button 
                            onClick={() => setEmojiPickerMessageId(emojiPickerMessageId === msg.id ? null : msg.id)} 
                            title="React" 
                            className="p-1 hover:bg-muted rounded"
                          >
                            <SmileyIcon size={14} />
                          </button>
                          {emojiPickerMessageId === msg.id && (
                            <div className="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 bg-background border border-border rounded-lg shadow-lg p-2 flex gap-1 z-50">
                              {['👍', '❤️', '😂', '😮', '😢', '🙏', '🔥', '🎉'].map(emoji => (
                                <button
                                  key={emoji}
                                  onClick={() => {
                                    handleReact(msg, emoji);
                                    setEmojiPickerMessageId(null);
                                  }}
                                  className="text-lg hover:scale-125 transition-transform p-1"
                                  title={emoji}
                                >
                                  {emoji}
                                </button>
                              ))}
                            </div>
                          )}
                        </div>
                    </div>
                  </div>

                  {msg.reactions && msg.reactions.length > 0 && (
                    <div className="flex flex-wrap gap-1 mt-1">
                      {(() => {
                        const reactionGroups = msg.reactions.reduce((acc, r) => {
                          if (!acc[r.reaction]) acc[r.reaction] = [];
                          acc[r.reaction].push(r);
                          return acc;
                        }, {} as Record<string, typeof msg.reactions>);
                        
                        return Object.entries(reactionGroups).map(([emoji, reactions]) => (
                          <button
                            key={emoji}
                            onClick={() => handleReact(msg, emoji)}
                            className="flex items-center gap-1 px-2 py-0.5 rounded-full bg-muted/30 hover:bg-muted/50 transition-all text-xs border border-border/30"
                            title={reactions.map(r => r.username).join(', ')}
                          >
                            <span>{emoji}</span>
                            <span className="text-muted-foreground font-medium">{reactions.length}</span>
                          </button>
                        ));
                      })()}
                    </div>
                  )}

                  <div className="flex items-center gap-2">
                      <span className="text-[10px] text-muted px-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        {new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                      </span>
                  </div>
                </div>
              </div>
            );
          })
        )}
        <div ref={messagesEndRef} />
      </div>

      <div className="p-4 bg-card/80 backdrop-blur-md border-t border-border z-20">
        {replyTo && (
            <div className="flex items-center justify-between bg-muted/10 p-3 rounded-lg mb-2 border-l-4 border-primary">
                <div className="flex flex-col gap-1 flex-1 min-w-0">
                  <span className="text-xs font-semibold text-primary">
                    Replying to {isDirectMessage(replyTo) ? replyTo.senderUsername : replyTo.username}
                  </span>
                  <span className="text-sm text-muted-foreground truncate">
                    {getReplyMessageContent(replyTo.id)}
                  </span>
                </div>
                <button onClick={() => setReplyTo(null)} className="hover:text-foreground ml-2 flex-shrink-0"><XIcon size={18} /></button>
            </div>
        )}
        <form onSubmit={handleSubmit} className="relative max-w-4xl mx-auto flex gap-2 items-end">
           <input 
            type="file" 
            ref={fileInputRef} 
            className="hidden" 
            onChange={handleFileUpload} 
           />
           <button 
            type="button" 
            onClick={() => fileInputRef.current?.click()}
            className="p-3 text-muted-foreground hover:text-foreground hover:bg-muted/50 rounded-xl transition-all"
           >
            <PaperclipIcon size={20} />
           </button>

          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={disabled}
            placeholder={disabled ? "Connecting..." : placeholder}
            className="flex-1 px-4 py-3 rounded-2xl border border-border bg-background focus:bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          />
         
          <button
            type="submit"
            disabled={!input.trim() || disabled}
            className="p-3 bg-primary text-white rounded-xl hover:bg-primary-hover disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md active:scale-95"
          >
            <PaperPlaneIcon size={20} weight="fill" />
          </button>
        </form>
      </div>
    </div>
  );
}
