export interface TMDBSearchResult {
  page: number;
  results: TMDBSearchItem[];
  total_pages: number;
  total_results: number;
}

export interface TMDBSearchItem {
  id: number;
  title?: string;
  name?: string; 
  overview: string;
  poster_path: string;
  media_type: string;
  genre_ids: number[];
  release_date?: string;
  first_air_date?: string;
}

export type TMDBMovieSearch = TMDBSearchResult;
export type TMDBTVShowSearch = TMDBSearchResult;

export interface TMDBGenre {
  id: number;
  name: string;
}

export interface TMDBMovie {
  id: number;
  title: string;
  overview: string;
  poster_path: string;
  release_date: string;
  genres: TMDBGenre[];
}

export interface TMDBTVShow {
  id: number;
  name: string;
  overview: string;
  poster_path: string;
  first_air_date: string;
  genres: TMDBGenre[];
}