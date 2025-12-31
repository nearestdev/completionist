import Link from "next/link";

export default function SearchHubPage() {
  const searchCategories = [
    { name: "Anime", href: "/search/anime" },
    { name: "Manga", href: "/search/manga" },
    { name: "Movies", href: "/search/movies" },
    { name: "TV Shows", href: "/search/tv" },
    { name: "Games", href: "/search/games" },
    { name: "Books", href: "/search/books" },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8 text-center">Search Media</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {searchCategories.map((category) => (
          <Link
            key={category.name}
            href={category.href}
            className="block p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 text-center"
          >
            <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
              {category.name}
            </h2>
          </Link>
        ))}
      </div>
    </div>
  );
}