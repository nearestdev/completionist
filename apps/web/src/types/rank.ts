export enum Rank {
    Scribe = 1,
    Chronicler = 2,
    Curator = 3,
    Preserver = 4,
    Warden = 5,
    Oracle = 6
}

export interface UserRank {
    id: number;
    userId: number;
    seasonId?: number;
    currentRank: Rank;
    currentElo: number;
    peakRank: Rank;
    createdAt: string;
    updatedAt: string;
}

export const RankNames: Record<Rank, string> = {
    [Rank.Scribe]: "Scribe",
    [Rank.Chronicler]: "Chronicler",
    [Rank.Curator]: "Curator",
    [Rank.Preserver]: "Preserver",
    [Rank.Warden]: "Warden",
    [Rank.Oracle]: "Oracle"
};

export const RankThresholds: Record<Rank, number> = {
    [Rank.Scribe]: 0,
    [Rank.Chronicler]: 150,
    [Rank.Curator]: 400,
    [Rank.Preserver]: 800,
    [Rank.Warden]: 1500,
    [Rank.Oracle]: 2500
};
