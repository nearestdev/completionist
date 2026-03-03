export interface LastFMAccount {
  userId: number;
  username: string;
  sessionKey: string;
  subscriber: number;
}

export interface RecentTracksResponse {
  recenttracks: {
    track: Track[];
  };
}

export interface Track {
  artist: {
    "#text": string;
  };
  name: string;
  album: {
    "#text": string;
  };
  image: {
    "#text": string;
    size: "small" | "medium" | "large" | "extralarge";
  }[];
  url: string;
  "@attr"?: {
    nowplaying: "true" | "false";
  };
}