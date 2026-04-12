"use client";

import { useState, useRef, useEffect } from "react";
import { ShareNetwork, Copy, ChatCircle, ArrowBendUpRight } from "@phosphor-icons/react";
import api from "@/services/api";

interface ShareMenuProps {
  postId: number;
  onRepost?: () => void;
}

export default function ShareMenu({ postId, onRepost }: ShareMenuProps) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  const copyLink = () => {
    navigator.clipboard.writeText(`${window.location.origin}/posts/${postId}`);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const repost = async () => {
    try {
      await api.post(`/posts/${postId}/share`, { content: "" });
      onRepost?.();
      setOpen(false);
    } catch {}
  };

  return (
    <div className="relative" ref={ref}>
      <button onClick={() => setOpen(!open)} className="text-muted hover:text-foreground transition-colors p-1">
        <ShareNetwork size={18} />
      </button>
      {open && (
        <div className="absolute right-0 top-8 bg-card border border-border rounded-xl shadow-lg py-1 w-48 z-30">
          <button onClick={repost} className="flex items-center gap-2 w-full px-4 py-2 text-sm text-foreground hover:bg-primary/5">
            <ArrowBendUpRight size={16} /> Repost
          </button>
          <button onClick={copyLink} className="flex items-center gap-2 w-full px-4 py-2 text-sm text-foreground hover:bg-primary/5">
            <Copy size={16} /> {copied ? "Copied!" : "Copy Link"}
          </button>
        </div>
      )}
    </div>
  );
}
