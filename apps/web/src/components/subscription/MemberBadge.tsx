"use client";

import { Crown } from "@phosphor-icons/react";

export default function MemberBadge() {
  return (
    <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 rounded-full bg-gradient-to-r from-primary/20 to-accent/20 text-primary" title="Member">
      <Crown size={12} weight="fill" />
    </span>
  );
}
