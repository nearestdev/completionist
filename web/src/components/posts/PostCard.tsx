import { PostWithDetails } from "@/types/posts";
import Link from "next/link";
import Image from "next/image";
import { 
  HeartIcon, 
  ChatCircleIcon, 
  ShareNetworkIcon, 
  DotsThreeIcon 
} from "@phosphor-icons/react/dist/ssr";

interface PostCardProps {
  post: PostWithDetails;
}

export default function PostCard({ post }: PostCardProps) {
  return (
    <article className="bg-card rounded-xl p-5 border border-border shadow-sm mb-5 transition-all hover:shadow-md">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-full bg-indigo-100 border border-primary/20 flex items-center justify-center text-primary font-bold overflow-hidden">
            {post.username.charAt(0).toUpperCase()}
          </div>
          <div>
            <Link href={`/profile/${post.username}`}>
              <h4 className="font-heading font-semibold text-sm text-foreground hover:text-primary transition-colors">
                {post.username}
              </h4>
            </Link>
            <span className="text-xs text-muted block">
              {new Date(post.createdAt).toLocaleDateString()}
            </span>
          </div>
        </div>
        <button className="text-muted hover:text-foreground">
          <DotsThreeIcon size={24} weight="bold" />
        </button>
      </div>

      <div className="mb-4 text-foreground text-sm leading-relaxed whitespace-pre-wrap">
        {post.content}
      </div>

      {post.mediaTitle && (
        <div className="flex gap-4 bg-background p-3 rounded-lg border-l-4 border-accent mb-4">
          <div className="relative w-16 h-24 flex-shrink-0 bg-gray-200 rounded-md overflow-hidden">
            {post.mediaCoverImage ? (
              <Image 
                src={post.mediaCoverImage} 
                alt={post.mediaTitle}
                fill
                className="object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-xs text-muted">
                No Img
              </div>
            )}
          </div>
          <div className="flex flex-col justify-center">
            <span className="text-[10px] uppercase tracking-wider font-bold text-muted bg-gray-200 dark:bg-gray-700 px-2 py-0.5 rounded w-fit mb-1">
              Media
            </span>
            <h5 className="font-heading font-semibold text-foreground">
              {post.mediaTitle}
            </h5>
          </div>
        </div>
      )}

      <div className="flex items-center gap-6 pt-3 border-t border-border text-muted text-sm">
        <button className="flex items-center gap-2 hover:text-primary transition-colors group">
          <HeartIcon size={20} className={post.isLikedByUser ? "text-red-500 fill-current" : "group-hover:text-red-500"} weight={post.isLikedByUser ? "fill" : "regular"} />
          <span>{post.likesCount}</span>
        </button>
        
        <Link href={`/posts/${post.id}`} className="flex items-center gap-2 hover:text-primary transition-colors">
          <ChatCircleIcon size={20} />
          <span>{post.commentsCount}</span>
        </Link>
        
        <button className="flex items-center gap-2 hover:text-primary transition-colors ml-auto">
          <ShareNetworkIcon size={20} />
          <span>Share</span>
        </button>
      </div>
    </article>
  );
}