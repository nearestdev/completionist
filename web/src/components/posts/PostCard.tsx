import { PostWithDetails } from "@/types/posts";
import Link from "next/link";
interface PostCardProps {
  post: PostWithDetails;
}
export default function PostCard({ post }: PostCardProps) {
  return (
    <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
      <div className="flex items-center mb-4">
        <div className="flex-shrink-0">
          <Link href={`/profile/${post.username}`}>
            <span className="font-bold text-indigo-600 hover:underline">{post.username}</span>
          </Link>
        </div>
      </div>
      {post.title && <h2 className="text-xl font-semibold mb-2">{post.title}</h2>}
      <p className="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{post.content}</p>
      <div className="mt-4 text-sm text-gray-500 dark:text-gray-400">
        <span>{post.likesCount} Likes</span>
        <span className="mx-2">·</span>
        <span>{post.commentsCount} Comments</span>
        <span className="mx-2">·</span>
        <time>{new Date(post.createdAt).toLocaleDateString()}</time>
      </div>
    </div>
  );
}