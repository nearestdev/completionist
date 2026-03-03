"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { SignOutIcon, UsersIcon, PlayIcon, PauseIcon } from '@phosphor-icons/react/dist/ssr';
import ChatBox from '@/components/messaging/ChatBox';
import { roomService } from '@/services/roomService';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { Room, RoomMessage, RoomChatPayload, RoomState, RoomStatePayload } from '@/types/messaging';

interface WatchRoomProps {
  room: Room;
}

export default function WatchRoom({ room }: WatchRoomProps) {
  const router = useRouter();
  const [messages, setMessages] = useState<RoomMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const [members, setMembers] = useState<number>(room.memberCount || 0);
  const [roomState, setRoomState] = useState<RoomState | null>(null);
  const [videoUrl, setVideoUrl] = useState('');
  const { subscribe, sendMessage } = useWebSocket();

  useEffect(() => {
    loadMessages();
    loadRoomState();
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

    const unsubscribeState = subscribe<RoomStatePayload>('room_state', (payload) => {
      if (payload.roomId === room.id) {
        setRoomState({
          roomId: payload.roomId,
          currentMediaUrl: payload.currentMediaUrl,
          currentMediaTitle: payload.currentMediaTitle,
          currentPositionMs: payload.currentPositionMs,
          isPlaying: payload.isPlaying,
          updatedBy: payload.updatedBy,
          updatedByUsername: payload.updatedByUsername,
          updatedAt: payload.updatedAt,
        });
      }
    });

    return () => {
      unsubscribeChat();
      unsubscribeState();
    };
  }, [subscribe, room.id]);

  const loadMessages = async () => {
    try {
      setLoading(true);
      const data = await roomService.getRoomMessages(room.id);
      setMessages(data);
    } catch (error) {
      console.error('Failed to load messages:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadRoomState = async () => {
    try {
      const state = await roomService.getRoomState(room.id);
      setRoomState(state);
    } catch (error) {
      console.error('Failed to load room state:', error);
    }
  };

  const handleSendMessage = (content: string) => {
    sendMessage('room_chat', {
      roomId: room.id,
      content,
    });
  };

  const handleLoadVideo = async () => {
    if (!videoUrl.trim()) return;

    try {
      await roomService.updateRoomState(room.id, {
        currentMediaUrl: videoUrl,
        currentMediaTitle: extractTitleFromUrl(videoUrl),
        currentPositionMs: 0,
        isPlaying: true,
      });
      setVideoUrl('');
    } catch (error) {
      console.error('Failed to load video:', error);
    }
  };

  const handlePlayPause = async () => {
    if (!roomState) return;

    try {
      await roomService.updateRoomState(room.id, {
        isPlaying: !roomState.isPlaying,
      });
    } catch (error) {
      console.error('Failed to update playback:', error);
    }
  };

  const handleLeaveRoom = async () => {
    try {
      await roomService.leaveRoom(room.id);
      router.push('/rooms');
    } catch (error) {
      console.error('Failed to leave room:', error);
    }
  };

  const extractYouTubeId = (url: string): string | null => {
    const match = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\s]+)/);
    return match ? match[1] : null;
  };

  const extractTitleFromUrl = (url: string): string => {
    return url.split('/').pop() || 'Untitled';
  };

  const youtubeId = roomState?.currentMediaUrl ? extractYouTubeId(roomState.currentMediaUrl) : null;

  return (
    <div className="flex h-full">
      <div className="flex-1 flex flex-col">
        <div className="p-4 border-b border-border">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-xl font-bold text-foreground">{room.name}</h2>
              {room.description && (
                <p className="text-sm text-muted">{room.description}</p>
              )}
            </div>
            <button
              onClick={handleLeaveRoom}
              className="flex items-center gap-2 px-4 py-2 text-sm border border-border rounded-xl hover:bg-red-500/10 hover:border-red-500 hover:text-red-500 transition-all"
            >
              <SignOutIcon size={18} />
              Leave
            </button>
          </div>

          <div className="flex gap-2">
            <input
              type="text"
              value={videoUrl}
              onChange={(e) => setVideoUrl(e.target.value)}
              placeholder="Paste YouTube or Vimeo URL..."
              className="flex-1 px-4 py-2 rounded-xl border border-border bg-background focus:outline-none focus:ring-2 focus:ring-primary"
            />
            <button
              onClick={handleLoadVideo}
              disabled={!videoUrl.trim()}
              className="px-6 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              Load
            </button>
          </div>
        </div>

        <div className="flex-1 flex flex-col justify-center p-6">
          {youtubeId ? (
            <div className="max-w-5xl mx-auto w-full">
              <div className="aspect-video w-full bg-black rounded-xl overflow-hidden">
                <iframe
                  width="100%"
                  height="100%"
                  src={`https://www.youtube.com/embed/${youtubeId}?autoplay=${roomState?.isPlaying ? '1' : '0'}`}
                  title="YouTube video player"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowFullScreen
                />
              </div>

              <div className="mt-4 p-4 bg-card border border-border rounded-xl">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-foreground truncate">
                      {roomState?.currentMediaTitle || 'No title'}
                    </h3>
                    <p className="text-sm text-muted">
                      Shared by {roomState?.updatedByUsername}
                    </p>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="flex items-center gap-2 text-muted">
                      <UsersIcon size={20} />
                      <span className="text-sm">{members} watching</span>
                    </div>
                    <button
                      onClick={handlePlayPause}
                      className="p-3 bg-primary text-white rounded-full hover:bg-primary/90 transition-all"
                    >
                      {roomState?.isPlaying ? (
                        <PauseIcon size={24} weight="fill" />
                      ) : (
                        <PlayIcon size={24} weight="fill" />
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="flex items-center justify-center h-full text-muted">
              <div className="text-center">
                <p>No video playing</p>
                <p className="text-sm text-muted/70">Share a YouTube link to start watching together</p>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="w-80 flex flex-col border-l border-border">
        <div className="p-4 border-b border-border">
          <h3 className="font-semibold text-foreground">Chat</h3>
        </div>
        <div className="flex-1 overflow-hidden">
          <ChatBox
            messages={messages}
            onSendMessage={handleSendMessage}
            loading={loading}
            placeholder="Chat while watching..."
          />
        </div>
      </div>
    </div>
  );
}
