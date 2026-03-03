"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import postService from "@/services/postService";
import { PostWithDetails, CommentWithDetails } from "@/types/posts";
import PostCard from "@/components/posts/PostCard";
import CommentCard from "@/components/posts/CommentCard";
export default function PostDetailPage() {
  const { postId } = useParams();
  const [post, setPost] = useState<PostWithDetails | null>(null);
  const [comments, setComments] = useState<CommentWithDetails[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    const fetchPostData = async () => {
      if (!postId) return;
      try {
        setLoading(true);
        const postData = await postService.getPostById(Number(postId));
        setPost(postData);
        const commentsData = await postService.getPostComments(Number(postId));
        setComments(commentsData);
      } catch (err) {
        setError("Failed to load post details.");
      } finally {
        setLoading(false);
      }
    };
    fetchPostData();
  }, [postId]);
  if (loading) return <div>Loading post...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!post) return <div>Post not found.</div>;
  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <PostCard post={post} />
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Comments</h2>
        <div className="space-y-4">
          {comments.map((comment) => (
            <CommentCard key={comment.id} comment={comment} />
          ))}
          {comments.length === 0 && <p>No comments yet.</p>}
        </div>
      </div>
    </div>
  );
}