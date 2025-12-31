import { Book } from "@/types/books";
import Image from "next/image";

interface BookResultCardProps {
  book: Book;
}

export default function BookResultCard({ book }: BookResultCardProps) {
  const thumbnailUrl = book.volumeInfo.imageLinks?.thumbnail;
  const imageUrl = thumbnailUrl
    ? thumbnailUrl.replace(/^http:/, "https:")
    : "/placeholder.svg";

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div className="relative h-64">
        <Image
          src={imageUrl}
          alt={book.volumeInfo.title}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="font-bold text-lg truncate">{book.volumeInfo.title}</h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {book.volumeInfo.authors?.join(", ")}
        </p>
      </div>
    </div>
  );
}