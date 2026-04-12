export interface CompletionReview {
  id: number;
  userId: number;
  mediaItemId: string;
  rating: number;
  tags?: string[];
  favoriteQuote?: string;
  reviewText?: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ReviewWithUser extends CompletionReview {
  username: string;
}

export interface NewReview {
  mediaItemId: string;
  rating: number;
  tags?: string[];
  favoriteQuote?: string;
  reviewText?: string;
  completedAt?: string;
}
