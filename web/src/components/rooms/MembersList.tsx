"use client";

import { RoomMember } from "@/types/messaging";
import { User, Circle } from "@phosphor-icons/react";

interface MembersListProps {
  members: RoomMember[];
  onlineUsers: Set<number>;
}

export default function MembersList({ members, onlineUsers }: MembersListProps) {
  const sortedMembers = [...members].sort((a, b) => {
    const aOnline = onlineUsers.has(a.userId);
    const bOnline = onlineUsers.has(b.userId);
    if (aOnline && !bOnline) return -1;
    if (!aOnline && bOnline) return 1;
    return a.username.localeCompare(b.username);
  });

  return (
    <div className="h-full flex flex-col bg-card/50">
      <div className="p-4 border-b border-border">
        <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">
          Members — {members.length}
        </h3>
      </div>
      <div className="flex-1 overflow-y-auto p-2 space-y-1">
        {sortedMembers.map((member) => {
          const isOnline = onlineUsers.has(member.userId);
          return (
            <div
              key={member.userId}
              className="flex items-center gap-3 p-2 rounded-lg hover:bg-white/5 transition-colors group"
            >
              <div className="relative">
                <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary group-hover:bg-primary/30 transition-colors">
                  <User size={16} weight="fill" />
                </div>
                <div className="absolute -bottom-0.5 -right-0.5 bg-background rounded-full p-0.5">
                  <Circle
                    size={10}
                    weight="fill"
                    className={isOnline ? "text-green-500" : "text-gray-500/50"}
                  />
                </div>
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <span className={`text-sm font-medium truncate ${isOnline ? 'text-foreground' : 'text-muted-foreground'}`}>
                    {member.username}
                  </span>
                  {member.isModerator && (
                    <span className="text-[10px] bg-primary/10 text-primary px-1.5 py-0.5 rounded uppercase font-bold tracking-wider">
                      MOD
                    </span>
                  )}
                </div>
                {isOnline && (
                  <p className="text-[10px] text-green-500/80 font-medium">Online</p>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
