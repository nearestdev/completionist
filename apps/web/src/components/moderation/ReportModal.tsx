"use client";

import { useState } from "react";
import { X } from "@phosphor-icons/react";
import Button from "@/components/ui/Button";
import moderationService from "@/services/moderationService";
import { ModerationEntityType } from "@/types/moderation";

interface ReportModalProps {
  entityType: ModerationEntityType;
  entityId: number;
  onClose: () => void;
}

export default function ReportModal({ entityType, entityId, onClose }: ReportModalProps) {
  const [reason, setReason] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!reason.trim()) return;
    setLoading(true);
    try {
      await moderationService.reportContent({ entityType, entityId, reason: reason.trim() });
      setSubmitted(true);
    } catch {}
    setLoading(false);
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl p-6 w-full max-w-sm border border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-foreground">Report Content</h2>
          <button onClick={onClose} className="text-muted hover:text-foreground"><X size={20} /></button>
        </div>
        {submitted ? (
          <div className="text-center py-4">
            <p className="text-accent font-semibold">Report submitted. Thank you.</p>
            <Button variant="ghost" onClick={onClose} className="mt-4">Close</Button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Why are you reporting this content?"
              rows={3}
              className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm focus:ring-2 focus:ring-primary/50 outline-none resize-none"
            />
            <div className="flex gap-3">
              <Button type="button" variant="ghost" onClick={onClose} className="flex-1">Cancel</Button>
              <Button type="submit" variant="primary" disabled={loading || !reason.trim()} className="flex-1">
                {loading ? "Submitting..." : "Report"}
              </Button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
