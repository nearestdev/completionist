export interface TimeDonutEntry {
  itemType: string;
  count: number;
}

export interface GenreEntry {
  genre: string;
  count: number;
}

export interface UserStats {
  timeDonut: TimeDonutEntry[];
  genres: GenreEntry[];
}

export interface HeatmapDay {
  date: string;
  count: number;
}
