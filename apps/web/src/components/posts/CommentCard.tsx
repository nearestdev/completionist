import { CommentWithDetails } from "@/types/posts";
import Link from "next/link";
interface CommentCardProps {
  comment: CommentWithDetails;
}
export default function CommentCard({ comment }: CommentCardProps) {
  return (
    <div className="bg-white dark:bg-gray-800 shadow-sm rounded-lg p-4">
      <div className="flex items-center mb-2">
        <Link href={`/profile/${comment.username}`}>
          <span className="font-semibold text-indigo-600 hover:underline">
            {comment.username}
          </span>
        </Link>
        <span className="text-gray-500 dark:text-gray-400 text-sm ml-2">
          Â· {new Date(comment.createdAt).toLocaleDateString()}
        </span>
      </div>
      <p className="text-gray-700 dark:text-gray-300">{comment.content}</p>
      <div className="mt-2 text-sm text-gray-500 dark:text-gray-400">
        <span>{comment.likesCount} Likes</span>
      </div>
    </div>
  );
}