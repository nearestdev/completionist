"use client";

import { useEffect, useState } from "react";
import { Queue, X } from "@phosphor-icons/react";
import queueService, { QueueItem } from "@/services/queueService";
import Image from "next/image";

export default function NextUpSidebar() {
  const [items, setItems] = useState<QueueItem[]>([]);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    queueService.getMyQueue().then((data) => setItems(data || [])).catch(() => {});
  }, []);

  const remove = async (itemId: number) => {
    await queueService.removeFromQueue(itemId);
    setItems(items.filter((i) => i.id !== itemId));
  };

  return (
    <>
      <button onClick={() => setOpen(!open)} className="flex items-center gap-2 text-sm text-muted hover:text-foreground">
        <Queue size={18} /> Next Up ({items.length})
      </button>

      {open && (
        <div className="mt-3 space-y-2">
          {items.length === 0 ? (
            <p className="text-xs text-muted">Queue is empty. Add items from your list.</p>
          ) : (
            items.map((item) => (
              <div key={item.id} className="flex items-center gap-2 bg-card rounded-lg border border-border p-2">
                {item.coverImageUrl && (
                  <div className="w-8 h-11 rounded overflow-hidden flex-shrink-0">
                    <Image src={item.coverImageUrl} alt={item.title} width={32} height={44} className="w-full h-full object-cover" />
                  </div>
                )}
                <span className="text-xs text-foreground flex-1 truncate">{item.title}</span>
                <button onClick={() => remove(item.id)} className="text-muted hover:text-red-500">
                  <X size={12} />
                </button>
              </div>
            ))
          )}
        </div>
      )}
    </>
  );
}
