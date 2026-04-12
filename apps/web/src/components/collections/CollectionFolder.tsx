"use client";

import Link from "next/link";
import { FolderOpen, Lock } from "@phosphor-icons/react";
import { Collection } from "@/types/collection";

export default function CollectionFolder({ collection }: { collection: Collection }) {
  return (
    <Link href={`/collections/${collection.id}`}>
      <div className="bg-card rounded-xl border border-border p-4 hover:shadow-lg hover:border-primary/30 transition-all flex items-center gap-3">
        <FolderOpen size={24} className="text-primary flex-shrink-0" />
        <div className="flex-1 min-w-0">
          <h4 className="font-semibold text-foreground truncate">{collection.name}</h4>
          {collection.description && <p className="text-xs text-muted truncate">{collection.description}</p>}
        </div>
        {collection.isPrivate && <Lock size={14} className="text-muted flex-shrink-0" />}
      </div>
    </Link>
  );
}
