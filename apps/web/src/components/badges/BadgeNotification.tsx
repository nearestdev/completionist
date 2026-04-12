"use client";

import { useEffect, useState } from "react";
import { Trophy } from "@phosphor-icons/react";

interface BadgeNotificationProps {
  badgeName: string;
  onDismiss: () => void;
}

export default function BadgeNotification({ badgeName, onDismiss }: BadgeNotificationProps) {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false);
      onDismiss();
    }, 5000);
    return () => clearTimeout(timer);
  }, [onDismiss]);

  if (!visible) return null;

  return (
    <div className="fixed bottom-6 right-6 z-50 animate-in slide-in-from-bottom-4 fade-in">
      <div className="bg-card border-2 border-primary/30 rounded-xl p-4 shadow-lg flex items-center gap-3 max-w-xs">
        <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center">
          <Trophy size={20} weight="fill" className="text-primary" />
        </div>
        <div>
          <p className="text-xs text-muted">Badge Earned!</p>
          <p className="text-sm font-bold text-foreground">{badgeName}</p>
        </div>
      </div>
    </div>
  );
}
