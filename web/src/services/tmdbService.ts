import api from "./api";
import { TMDBMovieSearch, TMDBTVShowSearch } from "@/types/tmdb";

const tmdbService = {
  searchMovies: async (query: string, page: number = 1): Promise<TMDBMovieSearch> => {
    const response = await api.get<TMDBMovieSearch>("/search/movies", {
      params: { q: query, page },
    });
    return response.data;
  },

  searchTvShows: async (query: string, page: number = 1): Promise<TMDBTVShowSearch> => {
    const response = await api.get<TMDBTVShowSearch>("/search/tv", {
      params: { q: query, page },
    });
    return response.data;
  },
};

export default tmdbService;