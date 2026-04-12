"use client";

import { useEffect, useState } from "react";
import { Plus } from "@phosphor-icons/react";
import api from "@/services/api";
import franchiseService from "@/services/franchiseService";
import { Franchise } from "@/types/franchise";
import Button from "@/components/ui/Button";

export default function FranchiseManager() {
  const [franchises, setFranchises] = useState<Franchise[]>([]);
  const [showCreate, setShowCreate] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  useEffect(() => {
    franchiseService.list(100, 0).then((data) => setFranchises(data || [])).catch(() => {});
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    try {
      const res = await api.post<Franchise>("/admin/franchises", { name: name.trim(), description: description || undefined });
      setFranchises([res.data, ...franchises]);
      setShowCreate(false);
      setName("");
      setDescription("");
    } catch {}
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-foreground">Franchises</h3>
        <Button size="sm" variant="primary" onClick={() => setShowCreate(!showCreate)}>
          <Plus size={14} className="mr-1" /> New
        </Button>
      </div>

      {showCreate && (
        <form onSubmit={handleCreate} className="bg-card rounded-xl border border-border p-4 mb-4 space-y-3">
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Franchise name"
            className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none" />
          <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Description (optional)"
            className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none" />
          <Button type="submit" variant="primary" size="sm">Create</Button>
        </form>
      )}

      {franchises.length === 0 ? (
        <p className="text-muted">No franchises yet.</p>
      ) : (
        <div className="space-y-2">
          {franchises.map((f) => (
            <div key={f.id} className="bg-card rounded-lg border border-border p-3 flex items-center justify-between">
              <span className="font-medium text-foreground">{f.name}</span>
              <span className="text-xs text-muted">{new Date(f.createdAt).toLocaleDateString()}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
