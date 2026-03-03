export type PostType = "review" | "discussion" | "general";

export interface Post {
  id: number;
  userId: number;
  mediaItemId?: string;
  title?: string;
  content: string;
  postType: PostType;
  rating?: number;
  isSpoiler: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface NewPost {
  mediaItemId?: string;
  title?: string;
  content: string;
  postType: PostType;
  rating?: number;
  isSpoiler?: boolean;
}

export interface UpdatePost {
  title?: string;
  content?: string;
  rating?: number;
  isSpoiler?: boolean;
}

export interface Comment {
  id: number;
  userId: number;
  postId: number;
  parentCommentId?: number;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface NewComment {
  postId: number;
  parentCommentId?: number;
  content: string;
}

export interface UpdateComment {
  content: string;
}

export interface PostWithDetails {
  id: number;
  userId: number;
  username: string;
  mediaItemId?: string;
  mediaTitle?: string;
  mediaCoverImage?: string;
  title?: string;
  content: string;
  postType: PostType;
  rating?: number;
  isSpoiler: boolean;
  createdAt: string;
  updatedAt: string;
  likesCount: number;
  commentsCount: number;
  isLikedByUser: boolean;
}

export interface CommentWithDetails {
  id: number;
  userId: number;
  username: string;
  postId: number;
  parentCommentId?: number;
  content: string;
  createdAt: string;
  updatedAt: string;
  likesCount: number;
  isLikedByUser: boolean;
  replies?: CommentWithDetails[];
}