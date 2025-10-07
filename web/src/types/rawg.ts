export interface RAWGGameSearch {
  count: number;
  results: {
    id: number;
    name: string;
    released: string;
    background_image: string;
    rating: number;
    platforms: {
      platform: {
        id: number;
        name: string;
      };
    }[];
  }[];
}

export interface RAWGGameAchievements {
  count: number;
  results: {
    id: number;
    name: string;
    description: string;
    image: string;
    percent: number;
  }[];
}