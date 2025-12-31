"use client";

import { useState } from "react";
import bookService from "@/services/bookService";
import { Book } from "@/types/books";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import BookResultCard from "@/components/search/BookResultCard";

export default function BookSearchPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Book[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setLoading(true);
    setError(null);
    try {
      const response = await bookService.searchBooks(query);
      setResults(response.items || []);
      if (!response.items || response.items.length === 0) {
        setError("No books found.");
      }
    } catch (err) {
      setError("Failed to search for books.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Search Books</h1>
      <form onSubmit={handleSearch} className="flex gap-2 mb-8">
        <Input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search for a book..."
          className="flex-grow"
        />
        <Button type="submit" disabled={loading}>
          {loading ? "Searching..." : "Search"}
        </Button>
      </form>

      {error && <p className="text-red-500 text-center">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {results.map((book) => (
          <BookResultCard key={book.id} book={book} />
        ))}
      </div>
    </div>
  );
}