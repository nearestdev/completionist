"use client";

import { useEffect, useState } from "react";
import { Sparkle } from "@phosphor-icons/react";
import suggestionService from "@/services/suggestionService";
import { MediaItem } from "@/types/list";
import Image from "next/image";

export default function SuggestionRail() {
  const [items, setItems] = useState<MediaItem[]>([]);

  useEffect(() => {
    suggestionService.getSuggestions(10).then((data) => setItems(data || [])).catch(() => {});
  }, []);

  if (items.length === 0) return null;

  return (
    <div className="mb-6">
      <div className="flex items-center gap-2 mb-3">
        <Sparkle size={18} weight="fill" className="text-accent" />
        <h3 className="text-sm font-semibold text-foreground">Suggested for You</h3>
      </div>
      <div className="flex gap-3 overflow-x-auto pb-2">
        {items.map((item) => (
          <div key={item.id} className="flex-shrink-0 w-28">
            <div className="w-28 h-40 rounded-lg overflow-hidden bg-muted/10 mb-1">
              {item.coverImageUrl ? (
                <Image src={item.coverImageUrl} alt={item.title} width={112} height={160} className="w-full h-full object-cover" />
              ) : (
                <div className="w-full h-full flex items-center justify-center text-muted text-xs">{item.itemType}</div>
              )}
            </div>
            <p className="text-xs text-foreground font-medium truncate">{item.title}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
