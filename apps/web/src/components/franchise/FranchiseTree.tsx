"use client";

import { CheckCircle, Circle } from "@phosphor-icons/react";
import { FranchiseItemWithMedia } from "@/types/franchise";
import Image from "next/image";

export default function FranchiseTree({ items }: { items: FranchiseItemWithMedia[] }) {
  if (!items || items.length === 0) {
    return <p className="text-muted text-sm">No items in this franchise yet.</p>;
  }

  return (
    <div className="space-y-1">
      {items.map((item, index) => (
        <div key={item.id} className="flex items-start gap-3">
          <div className="flex flex-col items-center">
            {item.isCompleted ? (
              <CheckCircle size={24} weight="fill" className="text-accent flex-shrink-0" />
            ) : (
              <Circle size={24} className="text-muted flex-shrink-0" />
            )}
            {index < items.length - 1 && <div className="w-0.5 h-8 bg-border" />}
          </div>
          <div className="flex items-center gap-3 pb-6">
            {item.coverImageUrl && (
              <div className="w-10 h-14 rounded overflow-hidden flex-shrink-0">
                <Image src={item.coverImageUrl} alt={item.title} width={40} height={56} className="w-full h-full object-cover" />
              </div>
            )}
            <div>
              <p className={`text-sm font-medium ${item.isCompleted ? "text-foreground" : "text-muted"}`}>
                {item.title}
              </p>
              <p className="text-xs text-muted">
                {item.itemType}{item.relationship ? ` · ${item.relationship}` : ""}
              </p>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
