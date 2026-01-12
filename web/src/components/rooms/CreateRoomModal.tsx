"use client";

import { useState } from 'react';
import { XIcon, UsersIcon, GlobeIcon, LockKeyIcon, ChatCircleTextIcon, MusicNotesIcon, MonitorPlayIcon } from '@phosphor-icons/react/dist/ssr';
import { roomService } from '@/services/roomService';
import type { RoomType, CreateRoomRequest } from '@/types/messaging';

interface CreateRoomModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: () => void;
}

export default function CreateRoomModal({ isOpen, onClose, onCreated }: CreateRoomModalProps) {
  const [formData, setFormData] = useState<CreateRoomRequest>({
    name: '',
    description: '',
    roomType: 'chat',
    isPublic: true,
    maxMembers: 50,
  });
  const [loading, setLoading] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    try {
      setLoading(true);
      await roomService.createRoom(formData);
      onCreated();
      onClose();
      setFormData({
        name: '',
        description: '',
        roomType: 'chat',
        isPublic: true,
        maxMembers: 50,
      });
    } catch (error) {
      console.error('Failed to create room:', error);
    } finally {
      setLoading(false);
    }
  };

  const getRoomTypeIcon = (type: RoomType) => {
    switch (type) {
      case 'music': return <MusicNotesIcon size={20} />;
      case 'watch': return <MonitorPlayIcon size={20} />;
      default: return <ChatCircleTextIcon size={20} />;
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4 animate-in fade-in duration-200">
      <div className="bg-card border border-border rounded-2xl max-w-lg w-full p-0 shadow-2xl animate-in zoom-in-95 duration-200 overflow-hidden">
        <div className="p-6 border-b border-border flex items-center justify-between bg-muted/5">
          <div>
            <h2 className="text-2xl font-bold text-foreground">Create Room</h2>
            <p className="text-sm text-muted">Set up a new space for your community</p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-muted/10 text-muted hover:text-foreground rounded-full transition-colors"
          >
            <XIcon size={24} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <div>
            <label className="block text-sm font-semibold text-foreground mb-2">
              Room Name
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-4 py-3 rounded-xl border border-border bg-background/50 focus:bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all placeholder:text-muted/40"
              placeholder="e.g. The Chill Lounge"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-foreground mb-2">
              Description <span className="text-muted font-normal">(Optional)</span>
            </label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-4 py-3 rounded-xl border border-border bg-background/50 focus:bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all placeholder:text-muted/40 resize-none"
              placeholder="What is this room about?"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-semibold text-foreground mb-2">
              Room Type
            </label>
            <div className="grid grid-cols-3 gap-3">
              {(['chat', 'music', 'watch'] as RoomType[]).map((type) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => setFormData({ ...formData, roomType: type })}
                  className={`
                    flex flex-col items-center gap-2 p-3 rounded-xl border transition-all
                    ${formData.roomType === type
                      ? 'border-primary bg-primary/5 text-primary ring-1 ring-primary/20'
                      : 'border-border hover:border-primary/30 hover:bg-muted/5 text-muted hover:text-foreground'
                    }
                  `}
                >
                  {getRoomTypeIcon(type)}
                  <span className="text-sm font-medium capitalize">{type}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
             <div>
              <label className="block text-sm font-semibold text-foreground mb-2">
                Privacy
              </label>
              <div className="flex rounded-xl border border-border p-1 bg-muted/5">
                <button
                  type="button"
                  onClick={() => setFormData({ ...formData, isPublic: true })}
                  className={`
                    flex-1 flex items-center justify-center gap-2 py-2 rounded-lg text-sm font-medium transition-all
                    ${formData.isPublic
                      ? 'bg-card text-foreground shadow-sm'
                      : 'text-muted hover:text-foreground'
                    }
                  `}
                >
                  <GlobeIcon size={16} />
                  Public
                </button>
                <button
                  type="button"
                  onClick={() => setFormData({ ...formData, isPublic: false })}
                  className={`
                    flex-1 flex items-center justify-center gap-2 py-2 rounded-lg text-sm font-medium transition-all
                    ${!formData.isPublic
                      ? 'bg-card text-foreground shadow-sm'
                      : 'text-muted hover:text-foreground'
                    }
                  `}
                >
                  <LockKeyIcon size={16} />
                  Private
                </button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-foreground mb-2">
                Max Members
              </label>
              <div className="relative">
                <UsersIcon className="absolute left-3 top-1/2 -translate-y-1/2 text-muted" size={18} />
                <input
                  type="number"
                  value={formData.maxMembers}
                  onChange={(e) => setFormData({ ...formData, maxMembers: parseInt(e.target.value) })}
                  className="w-full pl-10 pr-4 py-2.5 rounded-xl border border-border bg-background focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all"
                  min="2"
                  max="100"
                />
              </div>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-6 py-3 rounded-xl border border-border font-medium hover:bg-muted/5 transition-all text-foreground"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || !formData.name.trim()}
              className="flex-1 px-6 py-3 rounded-xl bg-primary text-white font-medium hover:bg-primary-hover shadow-lg shadow-primary/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform active:scale-95"
            >
              {loading ? 'Creating...' : 'Create Room'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
