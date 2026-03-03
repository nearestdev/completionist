"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { PlusIcon, MagnifyingGlassIcon } from '@phosphor-icons/react/dist/ssr';
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
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
            Discover Rooms
          </h1>
          <p className="text-muted mt-1">Join the conversation or start your own community.</p>
        </div>
        <button
          onClick={() => setIsCreateModalOpen(true)}
          className="flex items-center justify-center gap-2 px-6 py-3 bg-primary hover:bg-primary-hover text-white rounded-xl shadow-lg hover:shadow-primary/25 transition-all transform hover:-translate-y-0.5 active:translate-y-0"
        >
          <PlusIcon size={20} weight="bold" />
          <span>Create Room</span>
        </button>
      </div>

      {/* Filter Tabs */}
      <div className="flex items-center gap-2 mb-8 overflow-x-auto pb-2 scrollbar-hide">
        <div className="p-1 bg-muted/10 border border-border rounded-xl flex items-center gap-1">
          {(['all', 'chat', 'music', 'watch'] as const).map((type) => (
            <button
              key={type}
              onClick={() => setSelectedType(type)}
              className={`
                px-4 py-2 rounded-lg text-sm font-medium transition-all capitalize whitespace-nowrap
                ${selectedType === type
                  ? 'bg-card text-foreground shadow-sm'
                  : 'text-muted hover:text-foreground hover:bg-muted/10'
                }
              `}
            >
              {type}
            </button>
          ))}
        </div>
      </div>

      {/* Content Grid */}
      {loading ? (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <div key={i} className="h-48 rounded-2xl bg-muted/10 animate-pulse border border-border/50" />
          ))}
        </div>
      ) : rooms.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-center bg-muted/5 rounded-3xl border border-dashed border-border/60">
          <div className="w-16 h-16 rounded-full bg-muted/10 flex items-center justify-center mb-4 text-muted">
            <MagnifyingGlassIcon size={32} />
          </div>
          <h3 className="text-xl font-semibold text-foreground mb-2">No rooms found</h3>
          <p className="text-muted/70 max-w-sm mb-6">
            There are no {selectedType !== 'all' ? selectedType : ''} rooms created yet. 
            Why not start one yourself?
          </p>
          <button
            onClick={() => setIsCreateModalOpen(true)}
            className="px-6 py-2 bg-card border border-border hover:border-primary/50 text-foreground rounded-xl hover:bg-muted/5 transition-all font-medium shadow-sm"
          >
            Create First Room
          </button>
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
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
