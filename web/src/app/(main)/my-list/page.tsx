"use client";

import { useEffect, useState } from "react";
import listService from "@/services/listService";
import { UserListItem } from "@/types/list";
import Link from "next/link";
import { 
  ListDashesIcon, 
  StarIcon, 
  CircleIcon, 
  CheckCircleIcon,
  PauseCircleIcon,
  XCircleIcon,
  ClockIcon,
  FunnelIcon,
  FilmStripIcon
} from "@phosphor-icons/react/dist/ssr";
import Button from "@/components/ui/Button";

export default function MyListPage() {
  const [items, setItems] = useState<UserListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchListItems = async () => {
      try {
        setLoading(true);
        const data = await listService.getMyListItems();
        setItems(data);
      } catch (err) {
        setError("Failed to load your list.");
      } finally {
        setLoading(false);
      }
    };
    fetchListItems();
  }, []);

  const getStatusConfig = (status: string) => {
    switch (status.toLowerCase()) {
      case "completed": return { color: "text-emerald-500 bg-emerald-500/10 border-emerald-200 dark:border-emerald-900", icon: CheckCircleIcon };
      case "planning": return { color: "text-blue-500 bg-blue-500/10 border-blue-200 dark:border-blue-900", icon: ClockIcon };
      case "current": return { color: "text-primary bg-primary/10 border-indigo-200 dark:border-indigo-900", icon: CircleIcon };
      case "paused": return { color: "text-amber-500 bg-amber-500/10 border-amber-200 dark:border-amber-900", icon: PauseCircleIcon };
      case "dropped": return { color: "text-red-500 bg-red-500/10 border-red-200 dark:border-red-900", icon: XCircleIcon };
      default: return { color: "text-gray-500 bg-gray-500/10 border-gray-200", icon: CircleIcon };
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>
  );

  if (error) return (
     <div className="flex items-center justify-center min-h-[40vh] text-red-500 font-medium">
      {error}
    </div>
  );

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl animate-in fade-in duration-500">
      <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-6 mb-10">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-indigo-500/10 rounded-xl text-indigo-600 dark:text-indigo-400">
            <ListDashesIcon size={32} weight="duotone" />
          </div>
          <div>
            <h1 className="text-3xl font-bold font-heading">My Collection</h1>
            <p className="text-muted">Track and manage your entertainment journey</p>
          </div>
        </div>

        <div className="flex gap-2 w-full md:w-auto overflow-x-auto pb-2 md:pb-0">
          <Button variant="outline" className="gap-2 whitespace-nowrap">
            <FunnelIcon size={18} />
            Filter Status
          </Button>
          <Link href="/search">
            <Button className="whitespace-nowrap shadow-lg shadow-primary/20">
              Add New Item
            </Button>
          </Link>
        </div>
      </div>

      {items.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item) => {
             const statusConfig = getStatusConfig(item.status);
             const StatusIcon = statusConfig.icon;

             return (
              <div
                key={item.id}
                className="group bg-card border border-border hover:border-primary/50 rounded-2xl p-5 hover:shadow-lg transition-all duration-300 flex flex-col justify-between h-[180px]"
              >
                <div className="flex justify-between items-start gap-4">
                   <div className="w-12 h-12 bg-muted/20 rounded-lg flex items-center justify-center text-muted-foreground shrink-0">
                      <FilmStripIcon size={24} weight="duotone" />
                   </div>
                   
                   <div className={`px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider border flex items-center gap-1.5 ${statusConfig.color}`}>
                     <StatusIcon size={14} weight="fill" />
                     {item.status}
                   </div>
                </div>

                <div>
                   <h2 className="text-lg font-bold text-foreground mb-1 line-clamp-2 leading-tight group-hover:text-primary transition-colors">
                     Media Item #{item.mediaItemId.substring(0, 8)}...
                   </h2>
                   <p className="text-xs text-muted">Added on {new Date(item.createdAt).toLocaleDateString()}</p>
                </div>

                <div className="pt-4 border-t border-border flex items-center justify-between">
                  <div className="flex items-center gap-1">
                    {item.rating ? (
                       <>
                        <StarIcon size={16} weight="fill" className="text-amber-500" />
                        <span className="font-bold text-foreground">{item.rating}</span>
                        <span className="text-muted text-xs">/ 5</span>
                       </>
                    ) : (
                      <span className="text-xs text-muted italic">Not rated</span>
                    )}
                  </div>
                  
                  <Link href={`/media/${item.itemType}/${item.mediaItemId}`} className="text-sm font-medium text-primary hover:text-primary-hover transition-colors">
                    View Details
                  </Link>
                </div>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="text-center py-24 bg-muted/5 rounded-3xl border border-dashed border-border group">
          <div className="w-20 h-20 bg-muted/20 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform duration-300">
             <ListDashesIcon size={40} className="text-muted opacity-50" weight="duotone" />
          </div>
          <h2 className="text-2xl font-bold text-foreground mb-2">Your collection is empty</h2>
          <p className="text-muted mb-8 max-w-sm mx-auto">
            Start building your personal library of movies, games, and more.
          </p>
          <Link href="/search">
            <Button size="lg" className="shadow-xl shadow-primary/20">
              Start Exploring
            </Button>
          </Link>
        </div>
      )}
    </div>
  );
}