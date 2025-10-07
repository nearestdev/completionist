"use client";

import { useState } from "react";
import socialService from "@/services/socialService";
import { UserFollowResponse } from "@/types/social";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import UserCard from "@/components/users/UserCard";

export default function UserSearchPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [results, setResults] = useState<UserFollowResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchTerm.trim()) return;
    setLoading(true);
    setError(null);
    try {
      const searchResults = await socialService.searchUsers(searchTerm);
      setResults(searchResults);
      if (searchResults.length === 0) {
        setError("No users found.");
      }
    } catch (err) {
      setError("Failed to search for users.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Find Users</h1>
      <form onSubmit={handleSearch} className="flex gap-2 mb-8">
        <Input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Search for users..."
          className="flex-grow"
        />
        <Button type="submit" disabled={loading}>
          {loading ? "Searching..." : "Search"}
        </Button>
      </form>

      {error && <p className="text-red-500 text-center">{error}</p>}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {results.map((user) => (
          <UserCard key={user.userId} user={user} />
        ))}
      </div>
    </div>
  );
}