"use client";

import { useEffect, useState } from "react";
import franchiseService from "@/services/franchiseService";
import { Franchise } from "@/types/franchise";
import FranchiseCard from "@/components/franchise/FranchiseCard";

const CATEGORIES = [
  { key: "", label: "All" },
  { key: "movies", label: "Movies" },
  { key: "series", label: "Series" },
  { key: "anime", label: "Anime" },
  { key: "games", label: "Games" },
  { key: "books", label: "Books" },
];

export default function FranchisesPage() {
  const [franchises, setFranchises] = useState<Franchise[]>([]);
  const [loading, setLoading] = useState(true);
  const [category, setCategory] = useState("");

  useEffect(() => {
    let cancelled = false;
    const fetchData = async () => {
      try {
        const data = await franchiseService.list(50, 0, category || undefined);
        if (!cancelled) setFranchises(data || []);
      } catch {
        if (!cancelled) setFranchises([]);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    setLoading(true);
    fetchData();
    return () => { cancelled = true; };
  }, [category]);

  return (
    <div className="container mx-auto px-4 py-8 max-w-5xl">
      <h1 className="text-3xl font-bold text-foreground mb-2">Franchises & Sagas</h1>
      <p className="text-muted mb-6">Explore media universes and track your completion progress across entire franchises.</p>

      <div className="flex gap-2 mb-8 overflow-x-auto pb-2">
        {CATEGORIES.map((cat) => (
          <button
            key={cat.key}
            onClick={() => setCategory(cat.key)}
            className={`px-4 py-2 rounded-full text-sm font-medium transition-all whitespace-nowrap ${
              category === cat.key
                ? "bg-primary text-white shadow-sm"
                : "bg-card border border-border text-muted hover:text-foreground hover:border-primary/30"
            }`}
          >
            {cat.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="bg-card rounded-xl border border-border h-52 animate-pulse" />
          ))}
        </div>
      ) : franchises.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-4xl mb-4">🎬</p>
          <p className="text-muted">
            {category
              ? `No ${category} franchises found yet. They'll appear as you add media to your library.`
              : "No franchises discovered yet. Browse and add media to start building your franchise map."}
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {franchises.map((f) => (
            <FranchiseCard key={f.id} franchise={f} />
          ))}
        </div>
      )}
    </div>
  );
}
