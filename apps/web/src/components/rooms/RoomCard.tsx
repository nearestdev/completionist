"use client";

import type { Room, RoomType } from '@/types/messaging';
import { UsersIcon, MusicNotesIcon, MonitorPlayIcon, LockKeyIcon, GlobeIcon, CaretRightIcon } from '@phosphor-icons/react/dist/ssr';

interface RoomCardProps {
  room: Room;
  onJoin: (roomId: number) => void;
}

export default function RoomCard({ room, onJoin }: RoomCardProps) {
  const getRoomIcon = (type: RoomType) => {
    switch (type) {
      case 'music':
        return <MusicNotesIcon size={24} weight="fill" />;
      case 'watch':
        return <MonitorPlayIcon size={24} weight="fill" />;
      default:
        return <UsersIcon size={24} weight="fill" />;
    }
  };

  const getGradient = (type: RoomType) => {
    switch (type) {
      case 'music':
        return 'from-purple-500/20 to-pink-500/20 text-purple-500';
      case 'watch':
        return 'from-blue-500/20 to-cyan-500/20 text-blue-500';
      default:
        return 'from-primary/20 to-primary-hover/20 text-primary';
    }
  };

  return (
    <div 
      className="group bg-card border border-border hover:border-primary/50 rounded-2xl p-5 transition-all hover:shadow-lg hover:shadow-primary/5 cursor-pointer relative overflow-hidden"
      onClick={() => onJoin(room.id)}
    >
      <div className="relative z-10">
        <div className="flex items-start justify-between mb-4">
          <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${getGradient(room.roomType)} flex items-center justify-center shadow-inner`}>
            {getRoomIcon(room.roomType)}
          </div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-medium px-2 py-1 rounded-full bg-muted/10 text-muted border border-border/50">
              {room.memberCount || 0} / {room.maxMembers}
            </span>
          </div>
        </div>

        <div className="mb-4">
          <div className="flex items-center justify-between gap-2 mb-1">
            <h3 className="font-semibold text-lg text-foreground group-hover:text-primary transition-colors truncate">
              {room.name}
            </h3>
            {room.isPublic ? (
                 <GlobeIcon size={16} className="text-muted/50 flex-shrink-0" />
              ) : (
                <LockKeyIcon size={16} className="text-muted/50 flex-shrink-0" />
              )}
          </div>
          {room.description && (
            <p className="text-sm text-muted line-clamp-2 min-h-[40px]">
              {room.description}
            </p>
          )}
          {!room.description && (
             <p className="text-sm text-muted/40 italic min-h-[40px]">
               No description provided
             </p>
          )}
        </div>

        <div className="flex items-center justify-between pt-4 border-t border-border/50">
           <span className={`text-xs font-semibold uppercase tracking-wider ${
              room.roomType === 'music' ? 'text-purple-500' :
              room.roomType === 'watch' ? 'text-blue-500' : 'text-primary'
           }`}>
             {room.roomType}
           </span>
           <button 
             className="flex items-center gap-1 text-sm font-medium text-foreground group-hover:text-primary transition-colors group/btn"
             onClick={(e) => {
                e.stopPropagation();
                onJoin(room.id);
             }}
           >
             Join Room
             <CaretRightIcon className="w-4 h-4 transition-transform group-hover/btn:translate-x-1" />
           </button>
        </div>
      </div>
    </div>
  );
}
