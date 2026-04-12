import { UserListItem } from "./list";

export interface Collection {
  id: number;
  userId: number;
  name: string;
  description?: string;
  isPrivate: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CollectionWithItems extends Collection {
  items: UserListItem[];
}

export interface NewCollection {
  name: string;
  description?: string;
  isPrivate?: boolean;
}
