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
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 transition-all hover:shadow-md">
      <div className="flex justify-between items-start mb-2">
        <div>
          <h3 className="font-bold text-gray-900 dark:text-gray-100">
            {challenge.challengeTitle || "Unknown Challenge"}
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            {challenge.challengeDescription || "No description available."}
          </p>
        </div>
        <div className="flex items-center gap-1 text-amber-500 font-semibold text-sm bg-amber-50 dark:bg-amber-900/20 px-2 py-1 rounded shrink-0 ml-2">
          <TrophyIcon weight="fill" />
          <span>{challenge.xpReward || 0} XP</span>
        </div>
      </div>

      <div className="mt-4">
        <div className="flex justify-between text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">
          <span>Progress</span>
          <span>
            {challenge.currentProgress} / {challenge.targetProgress}
          </span>
        </div>
        <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5 overflow-hidden">
          <div
            className={`h-2.5 rounded-full transition-all duration-500 ${
              challenge.isCompleted ? "bg-green-500" : "bg-indigo-600"
            }`}
            style={{ width: `${percentage}%` }}
          ></div>
        </div>
      </div>

      {challenge.isCompleted && (
        <div className="mt-3 flex items-center gap-2 text-green-600 dark:text-green-400 text-sm font-bold">
          <CheckCircleIcon size={20} weight="fill" />
          <span>Completed</span>
        </div>
      )}
    </div>
  );
}