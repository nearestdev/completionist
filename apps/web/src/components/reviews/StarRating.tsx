"use client";

import { useState } from "react";
import { Star } from "@phosphor-icons/react";

interface StarRatingProps {
  value: number;
  onChange?: (value: number) => void;
  max?: number;
  readOnly?: boolean;
  size?: number;
}

export default function StarRating({ value, onChange, max = 10, readOnly = false, size = 20 }: StarRatingProps) {
  const [hovered, setHovered] = useState(0);

  return (
    <div className="flex gap-0.5">
      {Array.from({ length: max }, (_, i) => i + 1).map((star) => (
        <button
          key={star}
          type="button"
          disabled={readOnly}
          className={`transition-colors ${readOnly ? "cursor-default" : "cursor-pointer hover:scale-110"}`}
          onMouseEnter={() => !readOnly && setHovered(star)}
          onMouseLeave={() => setHovered(0)}
          onClick={() => onChange?.(star)}
        >
          <Star
            size={size}
            weight={(hovered || value) >= star ? "fill" : "regular"}
            className={(hovered || value) >= star ? "text-yellow-400" : "text-muted/30"}
          />
        </button>
      ))}
    </div>
  );
}
