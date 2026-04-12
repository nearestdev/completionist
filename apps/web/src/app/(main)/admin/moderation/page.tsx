"use client";

import ModerationQueue from "@/components/admin/ModerationQueue";
import BanManager from "@/components/admin/BanManager";

export default function AdminModerationPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold text-foreground mb-8">Moderation</h1>
      <div className="space-y-8">
        <section>
          <h2 className="text-lg font-semibold text-foreground mb-4">Review Queue</h2>
          <ModerationQueue />
        </section>
        <section>
          <BanManager />
        </section>
      </div>
    </div>
  );
}
