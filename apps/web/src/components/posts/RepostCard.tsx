"use client";

import { useEffect, useState } from "react";
import { ArrowBendUpRight } from "@phosphor-icons/react";
import postService from "@/services/postService";
import { PostWithDetails } from "@/types/posts";

export default function RepostCard({ sharedPostId }: { sharedPostId: number }) {
  const [original, setOriginal] = useState<PostWithDetails | null>(null);

  useEffect(() => {
    postService.getPostById(sharedPostId).then(setOriginal).catch(() => {});
  }, [sharedPostId]);

  if (!original) {
    return (
      <div className="border border-border rounded-lg p-3 mt-2 bg-background/50">
        <p className="text-sm text-muted italic">Original post unavailable</p>
      </div>
    );
  }

  return (
    <div className="border border-border rounded-lg p-3 mt-2 bg-background/50">
      <div className="flex items-center gap-2 text-xs text-muted mb-1">
        <ArrowBendUpRight size={12} />
        <span className="font-medium">{original.username}</span>
      </div>
      <p className="text-sm text-foreground line-clamp-3">{original.content}</p>
    </div>
  );
}
