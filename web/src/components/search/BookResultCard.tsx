import { Book } from "@/types/books";
import Image from "next/image";
import { mediaService } from "@/services/api";
import { useRouter } from "next/navigation";

interface BookResultCardProps {
  book: Book;
}

export default function BookResultCard({ book }: BookResultCardProps) {
  const router = useRouter();
  const thumbnailUrl = book.volumeInfo.imageLinks?.thumbnail;
  const imageUrl = thumbnailUrl
    ? thumbnailUrl.replace(/^http:/, "https:")
    : "/placeholder.svg";

  const handleClick = async () => {
    try {
      const response = await mediaService.getMediaDetails("GOOGLE_BOOKS", book.id, "book");
      router.push(`/media/book/${response.data.id}`);
    } catch (e) {
      console.error("Failed to fetch media details", e);
    }
  };

  return (
    <div 
      className="bg-card rounded-xl border border-border overflow-hidden transition-all duration-300 hover:shadow-lg group cursor-pointer"
      style={{ boxShadow: 'var(--shadow-sm)' }}
      onClick={handleClick}
    >
      <div className="relative h-64 overflow-hidden">
        <Image
          src={imageUrl}
          alt={book.volumeInfo.title}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </div>
      <div className="p-4">
        <span className="text-[10px] uppercase tracking-wider font-bold text-primary bg-primary/10 px-2 py-0.5 rounded mb-2 inline-block">
          Book
        </span>
        <h3 className="font-heading font-bold text-foreground text-lg truncate mb-1">{book.volumeInfo.title}</h3>
        <p className="text-sm text-muted">
          {book.volumeInfo.authors?.join(", ")}
        </p>
      </div>
    </div>
  );
}