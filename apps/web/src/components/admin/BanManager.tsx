"use client";

import { useState } from "react";
import { Gavel } from "@phosphor-icons/react";
import api from "@/services/api";
import Button from "@/components/ui/Button";

export default function BanManager() {
  const [userId, setUserId] = useState("");
  const [reason, setReason] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);

  const handleBan = async () => {
    if (!userId || !reason) return;
    setLoading(true);
    try {
      await api.post(`/admin/users/${userId}/ban`, { reason });
      setResult(`User ${userId} banned.`);
      setUserId("");
      setReason("");
    } catch {
      setResult("Failed to ban user.");
    }
    setLoading(false);
  };

  const handleUnban = async () => {
    if (!userId) return;
    setLoading(true);
    try {
      await api.delete(`/admin/users/${userId}/ban`);
      setResult(`User ${userId} unbanned.`);
      setUserId("");
    } catch {
      setResult("Failed to unban user.");
    }
    setLoading(false);
  };

  return (
    <div className="bg-card rounded-xl border border-border p-6">
      <h3 className="text-lg font-semibold text-foreground flex items-center gap-2 mb-4">
        <Gavel size={20} /> Ban Manager
      </h3>
      <div className="space-y-3">
        <input
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          placeholder="User ID"
          className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none"
        />
        <input
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          placeholder="Ban reason"
          className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm outline-none"
        />
        <div className="flex gap-3">
          <Button variant="primary" size="sm" onClick={handleBan} disabled={loading}>Ban</Button>
          <Button variant="ghost" size="sm" onClick={handleUnban} disabled={loading}>Unban</Button>
        </div>
        {result && <p className="text-sm text-muted">{result}</p>}
      </div>
    </div>
  );
}
