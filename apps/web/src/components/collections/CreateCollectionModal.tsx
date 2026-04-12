"use client";

import { useState } from "react";
import { X } from "@phosphor-icons/react";
import Button from "@/components/ui/Button";
import collectionService from "@/services/collectionService";
import { Collection } from "@/types/collection";

interface CreateCollectionModalProps {
  onClose: () => void;
  onCreated: (collection: Collection) => void;
}

export default function CreateCollectionModal({ onClose, onCreated }: CreateCollectionModalProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setLoading(true);
    try {
      const c = await collectionService.create({ name: name.trim(), description: description || undefined, isPrivate });
      onCreated(c);
      onClose();
    } catch {}
    setLoading(false);
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl p-6 w-full max-w-sm border border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-foreground">New Collection</h2>
          <button onClick={onClose} className="text-muted hover:text-foreground"><X size={20} /></button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Collection name"
            className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none" />
          <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Description (optional)"
            className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none" />
          <label className="flex items-center gap-2 text-sm text-muted cursor-pointer">
            <input type="checkbox" checked={isPrivate} onChange={(e) => setIsPrivate(e.target.checked)} /> Private
          </label>
          <Button type="submit" variant="primary" disabled={loading} className="w-full">
            {loading ? "Creating..." : "Create"}
          </Button>
        </form>
      </div>
    </div>
  );
}
