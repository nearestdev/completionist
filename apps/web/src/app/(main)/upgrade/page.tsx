"use client";

import { useEffect } from "react";
import { useSearchParams } from "next/navigation";
import PlanCard from "@/components/subscription/PlanCard";
import { useAuth } from "@/hooks/useAuth";

export default function UpgradePage() {
  const { user } = useAuth();
  const searchParams = useSearchParams();
  const success = searchParams.get("success");

  const isMember = user?.role === "member" || user?.role === "admin";

  return (
    <div className="container mx-auto px-4 py-12 max-w-2xl">
      {success === "true" ? (
        <div className="text-center">
          <div className="text-6xl mb-4">🎉</div>
          <h1 className="text-3xl font-bold text-foreground mb-2">Welcome, Member!</h1>
          <p className="text-muted">Your membership is now active. Enjoy all the perks!</p>
        </div>
      ) : isMember ? (
        <div className="text-center">
          <h1 className="text-3xl font-bold text-foreground mb-2">You&apos;re already a member</h1>
          <p className="text-muted">Manage your subscription in Settings.</p>
        </div>
      ) : (
        <>
          <h1 className="text-3xl font-bold text-foreground text-center mb-8">Upgrade to Member</h1>
          <PlanCard />
        </>
      )}
    </div>
  );
}
