"use client";
import { useState } from "react";
import rawgService from "@/services/rawgService";
import { RAWGGameSearch } from "@/types/rawg";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import GameResultCard from "@/components/search/GameResultCard";

export default function GameSearchPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<RAWGGameSearch["results"]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const response = await rawgService.searchGames(query);
      setResults(response.results);
      if (response.results.length === 0) {
        setError("No games found.");
      }
    } catch (err) {
      setError("Failed to search for games.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Search Games</h1>
      <form onSubmit={handleSearch} className="flex gap-2 mb-8">
        <Input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for a game..."
          className="flex-grow"
        />
        <Button type="submit" disabled={loading}>
          {loading ? "Searching..." : "Search"}
        </Button>
      </form>

      {error && <p className="text-red-500 text-center">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {results.map((game) => (
          <GameResultCard key={game.id} game={game} />
        ))}
      </div>
    </div>
  );
}