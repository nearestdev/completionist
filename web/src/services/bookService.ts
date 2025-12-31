import api from "./api";
import { BookSearchResult } from "@/types/books";

const bookService = {
  searchBooks: async (
    query: string,
    page: number = 1
  ): Promise<BookSearchResult> => {
    const response = await api.get<BookSearchResult>("/search/books", {
      params: { q: query, page },
    });
    return response.data;
  },
};

export default bookService;