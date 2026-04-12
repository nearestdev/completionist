"use client";

import { useState } from "react";
import { X } from "@phosphor-icons/react";
import Button from "@/components/ui/Button";
import StarRating from "./StarRating";
import reviewService from "@/services/reviewService";

interface ReviewModalProps {
  mediaItemId: string;
  mediaTitle: string;
  onClose: () => void;
  onSaved?: () => void;
}

export default function ReviewModal({ mediaItemId, mediaTitle, onClose, onSaved }: ReviewModalProps) {
  const [rating, setRating] = useState(0);
  const [tags, setTags] = useState("");
  const [favoriteQuote, setFavoriteQuote] = useState("");
  const [reviewText, setReviewText] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (rating < 1) {
      setError("Please select a rating");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      await reviewService.create({
        mediaItemId,
        rating,
        tags: tags ? tags.split(",").map((t) => t.trim()).filter(Boolean) : undefined,
        favoriteQuote: favoriteQuote || undefined,
        reviewText: reviewText || undefined,
        completedAt: new Date().toISOString(),
      });
      onSaved?.();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to save review");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl p-6 w-full max-w-md border border-border animate-in zoom-in-95">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-foreground">Review: {mediaTitle}</h2>
          <button onClick={onClose} className="text-muted hover:text-foreground">
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-sm font-medium text-foreground block mb-2">Rating</label>
            <StarRating value={rating} onChange={setRating} />
          </div>

          <div>
            <label className="text-sm font-medium text-foreground block mb-1">Tags (comma-separated)</label>
            <input
              type="text"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="Must Rewatch, Comfort Food, Overrated"
              className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm focus:ring-2 focus:ring-primary/50 outline-none"
            />
          </div>

          <div>
            <label className="text-sm font-medium text-foreground block mb-1">Favorite Quote</label>
            <input
              type="text"
              value={favoriteQuote}
              onChange={(e) => setFavoriteQuote(e.target.value)}
              placeholder="A memorable line..."
              className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm focus:ring-2 focus:ring-primary/50 outline-none"
            />
          </div>

          <div>
            <label className="text-sm font-medium text-foreground block mb-1">Review</label>
            <textarea
              value={reviewText}
              onChange={(e) => setReviewText(e.target.value)}
              placeholder="What did you think?"
              rows={3}
              className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm focus:ring-2 focus:ring-primary/50 outline-none resize-none"
            />
          </div>

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <div className="flex gap-3">
            <Button type="button" variant="ghost" onClick={onClose} className="flex-1">
              Cancel
            </Button>
            <Button type="submit" variant="primary" disabled={loading} className="flex-1">
              {loading ? "Saving..." : "Save Review"}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
