"use client";

import { useState } from "react";
import Button from "@/components/ui/Button";
import api from "@/services/api";

export default function ProfileEditor({ initialHtml, initialTheme }: { initialHtml?: string; initialTheme?: Record<string, string> }) {
  const [html, setHtml] = useState(initialHtml || "");
  const [primaryColor, setPrimaryColor] = useState(initialTheme?.primaryColor || "#4F46E5");
  const [bgColor, setBgColor] = useState(initialTheme?.bgColor || "");
  const [loading, setLoading] = useState(false);
  const [saved, setSaved] = useState(false);

  const handleSave = async () => {
    setLoading(true);
    try {
      await api.put("/me/profile", {
        profileHtml: html,
        profileTheme: JSON.stringify({ primaryColor, bgColor }),
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch {}
    setLoading(false);
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="text-sm font-medium text-foreground block mb-1">Custom Bio (HTML)</label>
        <textarea
          value={html}
          onChange={(e) => setHtml(e.target.value)}
          rows={6}
          placeholder="<p>Your custom bio...</p>"
          className="w-full px-3 py-2 rounded-lg bg-background border border-border text-foreground text-sm font-mono focus:ring-2 focus:ring-primary/50 outline-none resize-none"
        />
      </div>
      <div className="flex gap-4">
        <div>
          <label className="text-sm font-medium text-foreground block mb-1">Accent Color</label>
          <input type="color" value={primaryColor} onChange={(e) => setPrimaryColor(e.target.value)} className="w-10 h-10 rounded cursor-pointer" />
        </div>
        <div>
          <label className="text-sm font-medium text-foreground block mb-1">Background</label>
          <input type="color" value={bgColor} onChange={(e) => setBgColor(e.target.value)} className="w-10 h-10 rounded cursor-pointer" />
        </div>
      </div>
      <Button onClick={handleSave} disabled={loading} variant="primary">
        {saved ? "Saved!" : loading ? "Saving..." : "Save Profile"}
      </Button>
    </div>
  );
}
