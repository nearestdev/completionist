"use client";

import Link from "next/link";
import Image from "next/image";
import { Franchise } from "@/types/franchise";

export default function FranchiseCard({ franchise }: { franchise: Franchise }) {
  return (
    <Link href={`/franchises/${franchise.id}`}>
      <article className="bg-card rounded-xl border border-border overflow-hidden transition-all hover:shadow-lg hover:border-primary/30">
        {franchise.coverImageUrl ? (
          <div className="w-full h-32 relative">
            <Image src={franchise.coverImageUrl} alt={franchise.name} fill className="object-cover" />
          </div>
        ) : (
          <div className="w-full h-32 bg-gradient-to-br from-primary/10 to-accent/10 flex items-center justify-center">
            <span className="text-3xl">🎬</span>
          </div>
        )}
        <div className="p-4">
          <h3 className="font-semibold text-foreground truncate">{franchise.name}</h3>
          {franchise.description && (
            <p className="text-sm text-muted mt-1 line-clamp-2">{franchise.description}</p>
          )}
        </div>
      </article>
    </Link>
  );
}
