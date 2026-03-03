"use client";
import { Suspense, useEffect, useState } from "react";
import { useRouter, useSearchParams, usePathname } from "next/navigation";
import Link from "next/link";
import tmdbService from "@/services/tmdbService";
import { TMDBSearchItem } from "@/types/tmdb";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import MovieResultCard from "@/components/search/MovieResultCard";
import { FilmStripIcon, MagnifyingGlassIcon, SpinnerGapIcon, ArrowLeftIcon } from "@phosphor-icons/react/dist/ssr";

function MovieSearchPageContent() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const urlQuery = searchParams.get("q") || "";
  
  const [query, setQuery] = useState(urlQuery);
  const [results, setResults] = useState<TMDBSearchItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const performSearch = async (searchQuery: string) => {
    if (!searchQuery.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const response = await tmdbService.searchMovies(searchQuery);
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

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    router.push(`${pathname}?q=${encodeURIComponent(query)}`);
  };

  useEffect(() => {
    if (urlQuery) {
      performSearch(urlQuery);
    }
  }, [urlQuery]);

  return (
    <div className="container mx-auto px-4 py-12 max-w-7xl animate-in fade-in duration-500">
      <div className="text-center mb-12 space-y-4">
        <div className="inline-flex items-center justify-center p-4 bg-gradient-to-br from-red-500 to-rose-600 rounded-2xl mb-4 shadow-lg">
          <FilmStripIcon size={40} weight="fill" className="text-white" />
        </div>
        <h1 className="text-4xl md:text-5xl font-heading font-bold text-foreground tracking-tight">
          Search Movies
        </h1>
        <p className="text-xl text-muted max-w-2xl mx-auto">
          Discover trending films and cinematic masterpieces
        </p>
      </div>

      <div className="mb-8">
        <Link 
          href="/search"
          className="inline-flex items-center gap-2 text-muted hover:text-primary transition-colors group"
        >
          <ArrowLeftIcon size={20} weight="bold" className="group-hover:-translate-x-1 transition-transform" />
          <span className="font-medium">Back to Categories</span>
        </Link>
      </div>

      <form onSubmit={handleSearch} className="max-w-3xl mx-auto mb-12">
        <div className="relative">
          <div className="absolute left-4 top-1/2 -translate-y-1/2 text-muted">
            <MagnifyingGlassIcon size={24} weight="bold" />
          </div>
          <Input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search for a movie..."
            className="pl-14 pr-32 py-4 text-lg rounded-xl"
          />
          <div className="absolute right-2 top-1/2 -translate-y-1/2">
            <Button type="submit" disabled={loading} className="px-6 py-2">
              {loading ? (
                <span className="flex items-center gap-2">
                  <SpinnerGapIcon size={20} className="animate-spin" />
                  Searching...
                </span>
              ) : (
                "Search"
              )}
            </Button>
          </div>
        </div>
      </form>

      {error && (
        <div className="text-center py-12">
          <p className="text-muted text-lg">{error}</p>
        </div>
      )}

      {results.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
          {results.map((movie) => (
            <MovieResultCard key={movie.id} movie={movie} />
          ))}
        </div>
      )}
    </div>
  );
}

export default function MovieSearchPage() {
  return (
    <Suspense fallback={null}>
      <MovieSearchPageContent />
    </Suspense>
  );
}
