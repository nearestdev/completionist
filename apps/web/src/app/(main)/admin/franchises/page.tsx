"use client";

import FranchiseManager from "@/components/admin/FranchiseManager";

export default function AdminFranchisesPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold text-foreground mb-8">Franchise Management</h1>
      <FranchiseManager />
    </div>
  );
}
