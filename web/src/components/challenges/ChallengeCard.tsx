import { UserChallenge } from "@/types/challenge";
import { CheckCircleIcon, TrophyIcon } from "@phosphor-icons/react/dist/ssr";

interface ChallengeCardProps {
  challenge: UserChallenge;
}

export default function ChallengeCard({ challenge }: ChallengeCardProps) {
  const percentage = Math.min(
    100,
    Math.round((challenge.currentProgress / challenge.targetProgress) * 100)
  );

  return (
    <div 
      className="bg-card rounded-xl border border-border p-5 transition-all duration-300 hover:shadow-lg"
      style={{ boxShadow: 'var(--shadow-sm)' }}
    >      <div className="flex justify-between items-start mb-3">
        <div>
          <h3 className="font-bold font-heading text-foreground text-lg">
            {challenge.challengeTitle || "Unknown Challenge"}
          </h3>
          <p className="text-sm text-muted mt-1">
            {challenge.challengeDescription || "No description available."}
          </p>
        </div>
        <div className="flex items-center gap-1 text-amber-600 dark:text-amber-400 font-semibold text-sm bg-gradient-to-r from-amber-50 to-yellow-50 dark:from-amber-900/20 dark:to-yellow-900/20 px-3 py-1.5 rounded-full shrink-0 ml-2 border border-amber-200 dark:border-amber-800">
          <TrophyIcon weight="fill" size={16} />
          <span>{challenge.xpReward || 0} XP</span>
        </div>
      </div>

      <div className="mt-4">
        <div className="flex justify-between text-xs font-medium text-muted mb-2">
          <span>Progress</span>
          <span className="font-semibold text-foreground">
            {challenge.currentProgress} / {challenge.targetProgress}
          </span>
        </div>
        <div className="w-full bg-muted/20 rounded-full h-3 overflow-hidden">
          <div
            className={`h-3 rounded-full transition-all duration-700 ease-out ${
              challenge.isCompleted ? "bg-gradient-to-r from-accent to-emerald-400" : "bg-gradient-to-r from-primary to-primary-hover"
            }`}
            style={{ width: `${percentage}%` }}
          ></div>
        </div>
      </div>

      {challenge.isCompleted && (
        <div className="mt-4 flex items-center gap-2 text-accent text-sm font-bold bg-accent/10 px-3 py-2 rounded-lg">
          <CheckCircleIcon size={20} weight="fill" />
          <span>Completed!</span>
        </div>
      )}
    </div>
  );
}