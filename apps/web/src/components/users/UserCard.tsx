import { UserFollowResponse } from "@/types/social";
import Link from "next/link";
import Button from "../ui/Button";
interface UserCardProps {
  user: UserFollowResponse;
}
export default function UserCard({ user }: UserCardProps) {
  return (
    <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-4 flex items-center justify-between">
      <div>
        <Link href={`/profile/${user.username}`}>
          <span className="font-bold text-lg text-indigo-600 hover:underline">{user.username}</span>
        </Link>
      </div>
      <Link href={`/profile/${user.username}`}>
         <Button>View Profile</Button>
      </Link>
    </div>
  );
}