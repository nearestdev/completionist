"use client";

import { useEffect, useState } from "react";
import badgeService from "@/services/badgeService";
import { BadgeWithEarned } from "@/types/badge";
import BadgeCard from "./BadgeCard";

export default function BadgeShelf({ userId }: { userId: number }) {
  const [badges, setBadges] = useState<BadgeWithEarned[]>([]);

  useEffect(() => {
    badgeService.getUserBadges(userId).then((data) => setBadges(data || [])).catch(() => {});
  }, [userId]);

  if (badges.length === 0) return null;

  const earned = badges.filter((b) => b.earnedAt);

  return (
    <div>
      <h3 className="text-sm font-semibold text-foreground mb-3">
        Badges ({earned.length}/{badges.length})
      </h3>
      <div className="flex gap-3 overflow-x-auto pb-2">
        {badges.map((badge) => (
          <div key={badge.id} className="flex-shrink-0 w-24">
            <BadgeCard badge={badge} />
          </div>
        ))}
      </div>
    </div>
  );
}
