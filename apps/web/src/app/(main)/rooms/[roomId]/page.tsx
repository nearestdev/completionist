"use client";

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import ChatRoom from '@/components/rooms/ChatRoom';
import MusicRoom from '@/components/rooms/MusicRoom';
import WatchRoom from '@/components/rooms/WatchRoom';
import { roomService } from '@/services/roomService';
import { useWebSocket } from '@/hooks/useWebSocket';
import type { Room, RoomJoinPayload } from '@/types/messaging';

export default function RoomPage() {
  const params = useParams();
  const router = useRouter();
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { sendMessage } = useWebSocket();
  const roomId = parseInt(params.roomId as string);

  useEffect(() => {
    if (!roomId || isNaN(roomId)) {
      setError('Invalid room ID');
      setLoading(false);
      return;
    }

    loadRoom();
  }, [roomId]);

  useEffect(() => {
    if (!room || !sendMessage) return;

    sendMessage<RoomJoinPayload>('room_join', {
      roomId: room.id,
      userId: 0,
      username: '',
    });
  }, [room, sendMessage]);

  const loadRoom = async () => {
    try {
      setLoading(true);
      const data = await roomService.getRoom(roomId);
      
      if (!data.isMember) {
        await roomService.joinRoom(data.id);
        const updatedRoom = await roomService.getRoom(roomId);
        setRoom(updatedRoom);
      } else {
        setRoom(data);
      }
    } catch (err) {
      console.error('Failed to load room:', err);
      setError('Failed to load room');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-muted">Loading room...</div>
      </div>
    );
  }

  if (error || !room) {
    return (
      <div className="flex flex-col items-center justify-center h-screen gap-4">
        <div className="text-xl font-semibold text-muted">{error || 'Room not found'}</div>
        <button
          onClick={() => router.push('/rooms')}
          className="px-6 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-all"
        >
          Back to Rooms
        </button>
      </div>
    );
  }

  return (
    <div className="h-[calc(100vh-80px)]">
      {room.roomType === 'chat' && <ChatRoom room={room} />}
      {room.roomType === 'music' && <MusicRoom room={room} />}
      {room.roomType === 'watch' && <WatchRoom room={room} />}
    </div>
  );
}
