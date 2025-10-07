export interface SteamAccount {
  userId: number;
  steamId: string;
  persona: string;
  avatar: string;
}

export interface AttachSteamRequest {
  steamId: string;
}

export interface SteamOwnedGames {
  response: {
    game_count: number;
    games: {
      appid: number;
      name: string;
      playtime_forever: number;
      img_icon_url: string;
      img_logo_url: string;
      has_community_visible_stats: boolean;
      playtime_windows_forever: number;
      playtime_mac_forever: number;
      playtime_linux_forever: number;
      rtime_last_played: number;
    }[];
  };
}

export interface SteamAchievements {
  playerstats: {
    steamID: string;
    gameName: string;
    achievements: {
      apiname: string;
      achieved: number;
      unlocktime: number;
    }[];
    success: boolean;
  };
}

export interface SteamGameSchema {
  game: {
    gameName: string;
    gameVersion: string;
    availableGameStats: {
      achievements: {
        name: string;
        defaultvalue: number;
        displayName: string;
        hidden: number;
        icon: string;
        icongray: string;
      }[];
    };
  };
}