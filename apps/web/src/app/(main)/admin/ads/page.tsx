"use client";

import AdManager from "@/components/admin/AdManager";

export default function AdminAdsPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold text-foreground mb-8">Ad Campaigns</h1>
      <AdManager />
    </div>
  );
}
