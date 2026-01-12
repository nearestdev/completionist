"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { PlusIcon } from '@phosphor-icons/react/dist/ssr';
import RoomCard from '@/components/rooms/RoomCard';
import CreateRoomModal from '@/components/rooms/CreateRoomModal';
import { roomService } from '@/services/roomService';
import type { Room, RoomType } from '@/types/messaging';

export default function RoomsPage() {
  const router = useRouter();
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedType, setSelectedType] = useState<RoomType | 'all'>('all');
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  useEffect(() => {
    loadRooms();
  }, [selectedType]);

  const loadRooms = async () => {
    try {
      setLoading(true);
      const type = selectedType === 'all' ? undefined : selectedType;
      const data = await roomService.getRooms(type);
      setRooms(data);
    } catch (error) {
      console.error('Failed to load rooms:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleJoinRoom = async (roomId: number) => {
    try {
      await roomService.joinRoom(roomId);
      router.push(`/rooms/${roomId}`);
    } catch (error) {
      console.error('Failed to join room:', error);
    }
  };

  const handleRoomCreated = () => {
    loadRooms();
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-foreground">Rooms</h1>
        <button
          onClick={() => setIsCreateModalOpen(true)}
          className="flex items-center gap-2 px-6 py-3 bg-primary text-white rounded-xl hover:bg-primary/90 transition-all"
        >
          <PlusIcon size={20} weight="bold" />
          Create Room
        </button>
      </div>

      <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
        {(['all', 'chat', 'music', 'watch'] as const).map((type) => (
          <button
            key={type}
            onClick={() => setSelectedType(type)}
            className={`
              px-4 py-2 rounded-xl border transition-all capitalize whitespace-nowrap
              ${selectedType === type
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-border hover:border-primary/50'
              }
            `}
          >
            {type}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <div className="text-muted">Loading rooms...</div>
        </div>
      ) : rooms.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <h3 className="text-xl font-semibold text-muted mb-2">No rooms found</h3>
          <p className="text-sm text-muted/70 mb-4">
            Be the first to create a room!
          </p>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="px-6 py-2 bg-primary text-white rounded-xl hover:bg-primary/90 transition-all"
          >
            Create Room
          </button>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {rooms.map((room) => (
            <RoomCard key={room.id} room={room} onJoin={handleJoinRoom} />
          ))}
        </div>
      )}

      <CreateRoomModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreated={handleRoomCreated}
      />
    </div>
  );
}
