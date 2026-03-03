import api from "./api";
import { JikanAnimeSearch, JikanMangaSearch } from "@/types/jikan";

const jikanService = {
  searchAnime: async (query: string, page: number = 1): Promise<JikanAnimeSearch> => {
    const response = await api.get<JikanAnimeSearch>("/search/anime", {
      params: { q: query, page },
    });
    return response.data;
  },

  searchManga: async (query: string, page: number = 1): Promise<JikanMangaSearch> => {
    const response = await api.get<JikanMangaSearch>("/search/manga", {
      params: { q: query, page },
    });
    return response.data;
  },
};

export default jikanService;