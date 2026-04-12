"use client";

import { useState } from "react";
import { Shuffle } from "@phosphor-icons/react";
import SpinnerOverlay from "./SpinnerOverlay";
import { useAuth } from "@/hooks/useAuth";

export default function RandomizerFAB() {
  const [open, setOpen] = useState(false);
  const { user } = useAuth();

  if (!user) return null;

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="fixed bottom-6 right-6 z-40 w-14 h-14 rounded-full bg-gradient-to-r from-primary to-accent text-white shadow-lg hover:scale-105 transition-transform flex items-center justify-center"
        title="I'm bored, pick something!"
      >
        <Shuffle size={24} weight="bold" />
      </button>
      {open && <SpinnerOverlay onClose={() => setOpen(false)} />}
    </>
  );
}
