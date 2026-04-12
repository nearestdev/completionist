"use client";

import { useEffect, useState } from "react";
import { X, ArrowRight } from "@phosphor-icons/react";
import listService from "@/services/listService";
import Button from "@/components/ui/Button";

interface SpinnerOverlayProps {
  onClose: () => void;
}

export default function SpinnerOverlay({ onClose }: SpinnerOverlayProps) {
  const [items, setItems] = useState<any[]>([]);
  const [picked, setPicked] = useState<any | null>(null);
  const [spinning, setSpinning] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    listService.getMyListItems().then((data: any[]) => setItems((data || []).filter((i: any) => i.status === "planning"))).catch(() => {});
  }, []);

  const spin = () => {
    if (items.length === 0) return;
    setSpinning(true);
    setPicked(null);

    let ticks = 0;
    const totalTicks = 20 + Math.floor(Math.random() * 10);
    const interval = setInterval(() => {
      setCurrentIndex(Math.floor(Math.random() * items.length));
      ticks++;
      if (ticks >= totalTicks) {
        clearInterval(interval);
        const finalIndex = Math.floor(Math.random() * items.length);
        setCurrentIndex(finalIndex);
        setPicked(items[finalIndex]);
        setSpinning(false);
      }
    }, 80 + ticks * 5);
  };

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-card rounded-2xl p-8 w-full max-w-md border border-border text-center">
        <div className="flex justify-end mb-4">
          <button onClick={onClose} className="text-muted hover:text-foreground">
            <X size={20} />
          </button>
        </div>

        <h2 className="text-2xl font-bold text-foreground mb-6">What should I do next?</h2>

        {items.length === 0 ? (
          <p className="text-muted">No items in your planning list to pick from.</p>
        ) : (
          <>
            <div className="h-32 flex items-center justify-center mb-6">
              {spinning || picked ? (
                <div className="text-center">
                  <div className={`text-lg font-bold text-foreground ${spinning ? "animate-pulse" : ""}`}>
                    {items[currentIndex]?.media?.title || items[currentIndex]?.title || "..."}
                  </div>
                  {picked && !spinning && (
                    <p className="text-sm text-accent mt-2">This is the one!</p>
                  )}
                </div>
              ) : (
                <p className="text-muted">Press the button to pick something random</p>
              )}
            </div>

            <Button onClick={spin} disabled={spinning} variant="primary" size="lg" className="w-full">
              {spinning ? "Picking..." : picked ? "Spin Again" : "Pick for Me!"}
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
