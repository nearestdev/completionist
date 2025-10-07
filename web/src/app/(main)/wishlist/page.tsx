"use client";
import { useEffect, useState } from "react";
import listService from "@/services/listService";
import { WishlistItem } from "@/types/list";
import Link from "next/link";
export default function WishlistPage() {
  const [items, setItems] = useState<WishlistItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    const fetchWishlistItems = async () => {
      try {
        setLoading(true);
        const data = await listService.getMyWishlist();
        setItems(data);
      } catch (err) {
        setError("Failed to load your wishlist.");
      } finally {
        setLoading(false);
      }
    };
    fetchWishlistItems();
  }, []);
  if (loading) return <div>Loading your wishlist...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">My Wishlist</h1>
      {items.length > 0 ? (
        <div className="space-y-4">
          {items.map((item) => (
            <div
              key={item.id}
              className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-4"
            >
              <h2 className="text-xl font-semibold">
                Media ID: {item.mediaItemId}
              </h2>
              <p>Priority: {item.priority}</p>
            </div>
          ))}
        </div>
      ) : (
        <p>
          Your wishlist is empty. Go{" "}
          <Link href="/search" className="text-indigo-600 hover:underline">
            search
          </Link>{" "}
          for something to add!
        </p>
      )}
    </div>
  );
}