"use client";

import { useEffect, useRef, useState } from "react";
import adService from "@/services/adService";
import { AdCampaign, AdPlacement } from "@/types/ad";
import { useAuth } from "@/hooks/useAuth";

export default function AdSlot({ placement = "feed" }: { placement?: AdPlacement }) {
  const { user } = useAuth();
  const [ad, setAd] = useState<AdCampaign | null>(null);
  const ref = useRef<HTMLDivElement>(null);
  const impressionSent = useRef(false);

  const isMember = user?.role === "member" || user?.role === "admin";

  useEffect(() => {
    if (isMember) return;
    adService.getActiveAd(placement).then(setAd).catch(() => {});
  }, [placement, isMember]);

  useEffect(() => {
    if (!ad || !ref.current || impressionSent.current) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !impressionSent.current) {
          impressionSent.current = true;
          adService.recordImpression(ad.id, "view");
        }
      },
      { threshold: 0.5 }
    );
    observer.observe(ref.current);
    return () => observer.disconnect();
  }, [ad]);

  if (isMember || !ad) return null;

  return (
    <div ref={ref} className="rounded-xl border border-border overflow-hidden bg-card">
      <a
        href={ad.targetUrl}
        target="_blank"
        rel="noopener noreferrer"
        onClick={() => adService.recordImpression(ad.id, "click")}
      >
        {ad.imageUrl && (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={ad.imageUrl} alt={ad.name} className="w-full h-auto" />
        )}
        <div className="p-3">
          <div dangerouslySetInnerHTML={{ __html: ad.contentHtml }} className="text-sm text-foreground" />
          <p className="text-[10px] text-muted mt-1">Sponsored</p>
        </div>
      </a>
    </div>
  );
}
