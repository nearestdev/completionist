"use client";

import { useState } from "react";
import { Crown, Check } from "@phosphor-icons/react";
import Button from "@/components/ui/Button";
import subscriptionService from "@/services/subscriptionService";

const features = [
  "No advertisements",
  "Custom profile themes & HTML",
  "Priority support",
  "Early access to features",
  "Member badge on profile",
];

export default function PlanCard() {
  const [loading, setLoading] = useState(false);

  const handleUpgrade = async () => {
    setLoading(true);
    try {
      const { url } = await subscriptionService.createCheckout();
      window.location.href = url;
    } catch {
      setLoading(false);
    }
  };

  return (
    <div className="bg-card rounded-2xl border-2 border-primary/30 p-8 max-w-sm mx-auto text-center">
      <div className="w-16 h-16 rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center mx-auto mb-4">
        <Crown size={32} weight="fill" className="text-white" />
      </div>
      <h2 className="text-2xl font-bold text-foreground mb-1">Member</h2>
      <p className="text-muted text-sm mb-6">Unlock the full Completionist experience</p>

      <ul className="text-left space-y-3 mb-8">
        {features.map((f) => (
          <li key={f} className="flex items-center gap-2 text-sm text-foreground">
            <Check size={16} weight="bold" className="text-accent flex-shrink-0" />
            {f}
          </li>
        ))}
      </ul>

      <Button variant="primary" size="lg" onClick={handleUpgrade} disabled={loading} className="w-full">
        {loading ? "Redirecting..." : "Upgrade Now"}
      </Button>
    </div>
  );
}
