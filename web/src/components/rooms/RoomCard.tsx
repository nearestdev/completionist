"use client";

import type { Room, RoomType } from '@/types/messaging';
import { UsersIcon, MusicNotesIcon, MonitorPlayIcon, LockKeyIcon, GlobeIcon } from '@phosphor-icons/react/dist/ssr';

interface RoomCardProps {
  room: Room;
  onJoin: (roomId: number) => void;
}

export default function RoomCard({ room, onJoin }: RoomCardProps) {
  const getRoomTypeIcon = (type: RoomType) => {
    switch (type) {
      case 'music':
        return <MusicNotesIcon size={24} weight="fill" />;
      case 'watch':
        return <MonitorPlayIcon size={24} weight="fill" />;
      default:
        return <UsersIcon size={24} weight="fill" />;
    }
  };

  const getRoomTypeBadge = (type: RoomType) => {
    const colors = {
      chat: 'bg-blue-500/20 text-blue-500',
      music: 'bg-purple-500/20 text-purple-500',
      watch: 'bg-pink-500/20 text-pink-500',
    };
    return colors[type];
  };

  return (
    <div className="bg-card border border-border rounded-xl p-4 hover:border-primary/50 transition-all">
      <div className="flex items-start gap-4">
        <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center text-primary flex-shrink-0">
          {getRoomTypeIcon(room.roomType)}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2 mb-2">
            <h3 className="font-semibold text-foreground truncate">{room.name}</h3>
            <div className="flex items-center gap-1 text-muted flex-shrink-0">
              {room.isPublic ? (
                <GlobeIcon size={16} />
              ) : (
                <LockKeyIcon size={16} />
              )}
            </div>
          </div>

          {room.description && (
            <p className="text-sm text-muted line-clamp-2 mb-3">{room.description}</p>
          )}

          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <span className={`text-xs px-2 py-1 rounded-full font-medium ${getRoomTypeBadge(room.roomType)}`}>
                {room.roomType}
              </span>
              <span className="text-xs text-muted">
                {room.memberCount || 0} / {room.maxMembers} members
              </span>
            </div>

            <button
              onClick={() => onJoin(room.id)}
              className="px-4 py-1.5 text-sm bg-primary text-white rounded-full hover:bg-primary/90 transition-all"
            >
              Join
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
