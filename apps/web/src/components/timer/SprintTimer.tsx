"use client";

import { useState, useEffect, useRef } from "react";
import { Play, Pause, Stop, Timer } from "@phosphor-icons/react";
import Button from "@/components/ui/Button";

interface SprintTimerProps {
  onComplete?: () => void;
  durationMinutes?: number;
}

export default function SprintTimer({ onComplete, durationMinutes = 25 }: SprintTimerProps) {
  const totalSeconds = durationMinutes * 60;
  const [remaining, setRemaining] = useState(totalSeconds);
  const [running, setRunning] = useState(false);
  const [finished, setFinished] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    if (running && remaining > 0) {
      intervalRef.current = setInterval(() => {
        setRemaining((prev) => {
          if (prev <= 1) {
            clearInterval(intervalRef.current!);
            setRunning(false);
            setFinished(true);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }
    return () => { if (intervalRef.current) clearInterval(intervalRef.current); };
  }, [running, remaining]);

  const minutes = Math.floor(remaining / 60);
  const seconds = remaining % 60;
  const progress = ((totalSeconds - remaining) / totalSeconds) * 100;

  const reset = () => {
    setRemaining(totalSeconds);
    setRunning(false);
    setFinished(false);
  };

  return (
    <div className="bg-card rounded-xl border border-border p-4 inline-flex flex-col items-center gap-3">
      <div className="flex items-center gap-2 text-foreground">
        <Timer size={18} />
        <span className="text-sm font-medium">Sprint Mode</span>
      </div>

      <div className="relative w-20 h-20">
        <svg className="w-full h-full -rotate-90" viewBox="0 0 36 36">
          <circle cx="18" cy="18" r="15" fill="none" stroke="var(--border)" strokeWidth="3" />
          <circle
            cx="18" cy="18" r="15" fill="none" stroke="var(--primary)" strokeWidth="3"
            strokeDasharray={`${progress} ${100 - progress}`}
            strokeLinecap="round"
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center text-sm font-mono font-bold text-foreground">
          {String(minutes).padStart(2, "0")}:{String(seconds).padStart(2, "0")}
        </div>
      </div>

      {finished ? (
        <div className="text-center">
          <p className="text-sm text-accent font-semibold mb-2">Time&apos;s up!</p>
          <div className="flex gap-2">
            <Button size="sm" variant="primary" onClick={() => { onComplete?.(); reset(); }}>
              Log Progress
            </Button>
            <Button size="sm" variant="ghost" onClick={reset}>Reset</Button>
          </div>
        </div>
      ) : (
        <div className="flex gap-2">
          <button onClick={() => setRunning(!running)} className="p-2 rounded-full hover:bg-primary/10 text-foreground">
            {running ? <Pause size={20} /> : <Play size={20} />}
          </button>
          <button onClick={reset} className="p-2 rounded-full hover:bg-primary/10 text-muted">
            <Stop size={20} />
          </button>
        </div>
      )}
    </div>
  );
}
