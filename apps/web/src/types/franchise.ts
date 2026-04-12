export interface Franchise {
  id: number;
  name: string;
  description?: string;
  coverImageUrl?: string;
  createdBy?: number;
  createdAt: string;
  updatedAt: string;
}

export interface FranchiseItemWithMedia {
  id: number;
  franchiseId: number;
  mediaItemId: string;
  position: number;
  relationship?: string;
  title: string;
  coverImageUrl?: string;
  itemType: string;
  isCompleted: boolean;
}

export interface FranchiseWithItems extends Franchise {
  items: FranchiseItemWithMedia[];
  totalItems: number;
  completed: number;
}
