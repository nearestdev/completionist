"use client";
import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import postService from "@/services/postService";
import Button from "@/components/ui/Button";
import { PostType } from "@/types/posts";
export default function NewPostPage() {
  const [content, setContent] = useState("");
  const [postType, setPostType] = useState<PostType>("general");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await postService.createPost({ content, postType });
      router.push("/posts");
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to create post");
    } finally {
      setLoading(false);
    }
  };
  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-6">Create a New Post</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && <p className="text-red-500">{error}</p>}
        <div>
          <label
            htmlFor="content"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Content
          </label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={6}
            className="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            required
          />
        </div>
        <div>
          <label
            htmlFor="postType"
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Post Type
          </label>
          <select
            id="postType"
            value={postType}
            onChange={(e) => setPostType(e.target.value as PostType)}
            className="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-200 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="general">General</option>
            <option value="discussion">Discussion</option>
            <option value="review">Review</option>
          </select>
        </div>
        <Button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create Post"}
        </Button>
      </form>
    </div>
  );
}