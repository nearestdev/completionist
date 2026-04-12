"use client";

import { useEffect, useState } from "react";
import { Plus, Lock, FolderOpen } from "@phosphor-icons/react";
import collectionService from "@/services/collectionService";
import { Collection } from "@/types/collection";
import Button from "@/components/ui/Button";
import Link from "next/link";

export default function CollectionsPage() {
  const [collections, setCollections] = useState<Collection[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [name, setName] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    collectionService.getMyCollections().then((data) => setCollections(data || [])).catch(() => {});
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    setLoading(true);
    try {
      const c = await collectionService.create({ name: name.trim(), isPrivate });
      setCollections([c, ...collections]);
      setShowCreate(false);
      setName("");
      setIsPrivate(false);
    } catch {}
    setLoading(false);
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-foreground">Collections</h1>
        <Button variant="primary" size="sm" onClick={() => setShowCreate(!showCreate)}>
          <Plus size={16} weight="bold" className="mr-1" /> New
        </Button>
      </div>

      {showCreate && (
        <form onSubmit={handleCreate} className="bg-card rounded-xl border border-border p-4 mb-6 flex gap-3 items-end">
          <div className="flex-1">
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Collection name"
              className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm focus:ring-2 focus:ring-primary/50 outline-none"
            />
          </div>
          <label className="flex items-center gap-1 text-sm text-muted cursor-pointer">
            <input type="checkbox" checked={isPrivate} onChange={(e) => setIsPrivate(e.target.checked)} />
            Private
          </label>
          <Button type="submit" variant="primary" size="sm" disabled={loading}>Create</Button>
        </form>
      )}

      {collections.length === 0 ? (
        <p className="text-muted text-center py-12">No collections yet. Create one to organize your list!</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {collections.map((c) => (
            <Link key={c.id} href={`/collections/${c.id}`}>
              <div className="bg-card rounded-xl border border-border p-5 hover:shadow-lg hover:border-primary/30 transition-all">
                <div className="flex items-center gap-3">
                  <FolderOpen size={24} className="text-primary" />
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-foreground truncate">{c.name}</h3>
                    {c.description && <p className="text-sm text-muted truncate">{c.description}</p>}
                  </div>
                  {c.isPrivate && <Lock size={14} className="text-muted" />}
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
