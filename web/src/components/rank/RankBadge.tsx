import { Rank, RankNames, RankThresholds, UserRank } from "@/types/rank";
import { Scroll, Books, Gavel, Eye, Crown, PenNib } from "@phosphor-icons/react";

interface RankBadgeProps {
    rank: UserRank;
    showProgress?: boolean;
}

export default function RankBadge({ rank, showProgress = true }: RankBadgeProps) {
    const rankName = RankNames[rank.currentRank];
    
    const getIcon = (r: Rank) => {
        switch (r) {
            case Rank.Scribe: return <PenNib size={24} weight="duotone" />;
            case Rank.Chronicler: return <Scroll size={24} weight="duotone" />;
            case Rank.Curator: return <Books size={24} weight="duotone" />;
            case Rank.Preserver: return <Gavel size={24} weight="duotone" />; // Or Shield/Archive
            case Rank.Warden: return <Eye size={24} weight="duotone" />;
            case Rank.Oracle: return <Crown size={24} weight="duotone" />;
            default: return <PenNib size={24} />;
        }
    };

    const currentThreshold = RankThresholds[rank.currentRank];
    const nextRank = (rank.currentRank + 1) as Rank;
    const nextThreshold = RankThresholds[nextRank];
    
    let progress = 100;
    let nextRankName = "Max Rank";
    
    if (nextThreshold) {
        const range = nextThreshold - currentThreshold;
        const gained = rank.currentElo - currentThreshold;
        progress = Math.min(100, Math.max(0, (gained / range) * 100));
        nextRankName = RankNames[nextRank];
    }

    return (
        <div className="flex flex-col items-center p-4 bg-card rounded-lg border border-border shadow-sm">
            <div className={`p-3 rounded-full bg-primary/10 text-primary mb-2 ring-2 ring-primary/20`}>
                {getIcon(rank.currentRank)}
            </div>
            <h3 className="font-bold text-lg text-foreground">{rankName}</h3>
            <p className="text-sm text-muted">Elo: {rank.currentElo}</p>

            {showProgress && nextThreshold && (
                <div className="w-full mt-3">
                    <div className="flex justify-between text-xs text-muted mb-1">
                        <span>{rank.currentElo}</span>
                        <span>{nextThreshold} ({nextRankName})</span>
                    </div>
                    <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
                        <div 
                            className="bg-primary h-full transition-all duration-500" 
                            style={{ width: `${progress}%` }}
                        />
                    </div>
                </div>
            )}
        </div>
    );
}
