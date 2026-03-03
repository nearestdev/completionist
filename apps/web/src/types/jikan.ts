export interface JikanSearchResponse<T> {
  data: T[];
  pagination: {
    has_next_page: boolean;
  };
}

export interface JikanAnime {
  mal_id: number;
  title: string;
  images: {
    jpg: { image_url: string };
    webp: { image_url: string };
  };
  synopsis: string;
  type: string;
  genres: { name: string }[];
  aired: {
    from: string | null;
  };
  score: number | null;
  scored_by: number | null;
}

export interface JikanManga {
  mal_id: number;
  title: string;
  images: {
    jpg: { image_url: string };
    webp: { image_url: string };
  };
  synopsis: string;
  genres: { name: string }[];
  published: {
    from: string | null;
  };
  score: number | null;
  scored_by: number | null;
}

export type JikanAnimeSearch = JikanSearchResponse<JikanAnime>;
export type JikanMangaSearch = JikanSearchResponse<JikanManga>;