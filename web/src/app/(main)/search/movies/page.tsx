"use client";
import { useState } from "react";
import tmdbService from "@/services/tmdbService";
import { TMDBSearchItem } from "@/types/tmdb";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import MovieResultCard from "@/components/search/MovieResultCard";

export default function MovieSearchPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<TMDBSearchItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const response = await tmdbService.searchMovies(query);
      setResults(response.results);
      if (response.results.length === 0) {
        setError("No movies found.");
      }
    } catch (err) {
      setError("Failed to search for movies.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Search Movies</h1>
      <form onSubmit={handleSearch} className="flex gap-2 mb-8">
        <Input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for a movie..."
          className="flex-grow"
        />
        <Button type="submit" disabled={loading}>
          {loading ? "Searching..." : "Search"}
        </Button>
      </form>

      {error && <p className="text-red-500 text-center">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {results.map((movie) => (
          <MovieResultCard key={movie.id} movie={movie} />
        ))}
      </div>
    </div>
  );
}