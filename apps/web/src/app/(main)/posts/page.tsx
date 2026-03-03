"use client";

import { useEffect, useState } from "react";
import postService from "@/services/postService";
import { PostWithDetails } from "@/types/posts";
import PostCard from "@/components/posts/PostCard";
import Button from "@/components/ui/Button";

export default function PostsFeedPage() {
  const [posts, setPosts] = useState<PostWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const limit = 20;

  const fetchPosts = async (newOffset: number) => {
    try {
      setLoading(true);
      const newPosts = await postService.getPostsFeed(limit, newOffset);
      setPosts((prev) => (newOffset === 0 ? newPosts : [...prev, ...newPosts]));
      setOffset(newOffset + limit);
    } catch (err) {
      setError("Failed to load posts.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts(0);
  }, []);

  const handleLoadMore = () => {
    fetchPosts(offset);
  };

  if (loading && posts.length === 0) return <div>Loading posts...</div>;
  if (error) return <div className="text-red-500">{error}</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Activity Feed</h1>
      <div className="space-y-6">
        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>
      {posts.length > 0 && (
        <div className="mt-8 text-center">
          <Button onClick={handleLoadMore} disabled={loading}>
            {loading ? "Loading..." : "Load More"}
          </Button>
        </div>
      )}
      {posts.length === 0 && !loading && <p>No posts yet. Be the first to share something!</p>}
    </div>
  );
}