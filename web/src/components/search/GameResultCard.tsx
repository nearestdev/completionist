import { RAWGGameSearch } from "@/types/rawg";
import Image from "next/image";

type Game = RAWGGameSearch["results"][0];

interface GameResultCardProps {
  game: Game;
}

export default function GameResultCard({ game }: GameResultCardProps) {
  const imageUrl = game.background_image || "/placeholder.svg";

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div className="relative h-48">
        <Image
          src={imageUrl}
          alt={game.name}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="font-bold text-lg truncate">{game.name}</h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">{game.released?.substring(0, 4)}</p>
      </div>
    </div>
  );
}