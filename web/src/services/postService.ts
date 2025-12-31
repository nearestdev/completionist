import api from "./api";
import {
  NewPost,
  Post,
  PostWithDetails,
  UpdatePost,
  NewComment,
  Comment,
  CommentWithDetails,
  UpdateComment,
} from "@/types/posts";
const postService = {
  createPost: async (data: NewPost): Promise<Post> => {
    const response = await api.post<Post>("/posts", data);
    return response.data;
  },
  getPostsFeed: async (limit: number = 20, offset: number = 0): Promise<PostWithDetails[]> => {
    const response = await api.get<PostWithDetails[]>(`/posts?limit=${limit}&offset=${offset}`);
    return response.data;
  },
  getPostById: async (postId: number): Promise<PostWithDetails> => {
    const response = await api.get<PostWithDetails>(`/posts/${postId}`);
    return response.data;
  },
  getUserPosts: async (username: string, limit: number = 20, offset: number = 0): Promise<PostWithDetails[]> => {
    const response = await api.get<PostWithDetails[]>(`/users/${username}/posts?limit=${limit}&offset=${offset}`);
    return response.data;
  },
  updatePost: async (postId: number, data: UpdatePost): Promise<Post> => {
    const response = await api.patch<Post>(`/posts/${postId}`, data);
    return response.data;
  },
  deletePost: async (postId: number): Promise<void> => {
    await api.delete(`/posts/${postId}`);
  },
  likePost: async (postId: number): Promise<void> => {
    await api.post(`/posts/${postId}/like`);
  },
  unlikePost: async (postId: number): Promise<void> => {
    await api.delete(`/posts/${postId}/unlike`);
  },
  createComment: async (postId: number, data: NewComment): Promise<Comment> => {
    const response = await api.post<Comment>(`/posts/${postId}/comments`, data);
    return response.data;
  },
  getPostComments: async (postId: number): Promise<CommentWithDetails[]> => {
    const response = await api.get<CommentWithDetails[]>(`/posts/${postId}/comments`);
    return response.data;
  },
  updateComment: async (commentId: number, data: UpdateComment): Promise<Comment> => {
    const response = await api.patch<Comment>(`/comments/${commentId}`, data);
    return response.data;
  },
  deleteComment: async (commentId: number): Promise<void> => {
    await api.delete(`/comments/${commentId}`);
  },
  likeComment: async (commentId: number): Promise<void> => {
    await api.post(`/comments/${commentId}/like`);
  },
  unlikeComment: async (commentId: number): Promise<void> => {
    await api.delete(`/comments/${commentId}/unlike`);
  },
};
export default postService;