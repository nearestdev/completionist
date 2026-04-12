"use client";

import { Star, Quotes } from "@phosphor-icons/react";
import { ReviewWithUser } from "@/types/review";

export default function ReviewCard({ review }: { review: ReviewWithUser }) {
  return (
    <div className="bg-card rounded-xl border border-border p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-semibold text-foreground">{review.username}</span>
        <div className="flex items-center gap-1">
          <Star size={14} weight="fill" className="text-yellow-400" />
          <span className="text-sm font-bold text-foreground">{review.rating}/10</span>
        </div>
      </div>

      {review.tags && review.tags.length > 0 && (
        <div className="flex flex-wrap gap-1 mb-2">
          {review.tags.map((tag) => (
            <span key={tag} className="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary">{tag}</span>
          ))}
        </div>
      )}

      {review.favoriteQuote && (
        <div className="flex items-start gap-2 mb-2 bg-background/50 rounded-lg p-2">
          <Quotes size={14} className="text-muted flex-shrink-0 mt-0.5" />
          <p className="text-xs text-muted italic">{review.favoriteQuote}</p>
        </div>
      )}

      {review.reviewText && (
        <p className="text-sm text-foreground">{review.reviewText}</p>
      )}

      <p className="text-xs text-muted mt-2">
        {review.completedAt ? `Completed ${new Date(review.completedAt).toLocaleDateString()}` : ""}
      </p>
    </div>
  );
}
