"use client";

import { useEffect, useState } from "react";
import { ShieldCheckIcon, UsersIcon } from "@phosphor-icons/react/dist/ssr";
import { useAuth } from "@/hooks/useAuth";
import adminService from "@/services/adminService";
import { User } from "@/types/user";

export default function AdminPage() {
  const { isAdmin, isGuest } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isAdmin) {
      setLoading(false);
      return;
    }

    const fetchUsers = async () => {
      try {
        setLoading(true);
        const data = await adminService.listUsers();
        setUsers(data);
      } catch (err: any) {
        setError(err?.response?.data?.error || "Failed to load users.");
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, [isAdmin]);

  if (isGuest) {
    return null;
  }

  if (!isAdmin) {
    return (
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="rounded-2xl border border-red-500/20 bg-red-500/5 p-8 text-center">
          <h1 className="text-2xl font-bold text-foreground mb-2">403 - Forbidden</h1>
          <p className="text-muted">You do not have permission to access this area.</p>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-12 max-w-6xl">
        <div className="h-40 rounded-2xl border border-border bg-card/60 animate-pulse" />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl animate-in fade-in duration-500">
      <div className="flex items-center gap-3 mb-8">
        <div className="p-3 bg-primary/10 rounded-xl text-primary">
          <ShieldCheckIcon size={28} weight="duotone" />
        </div>
        <div>
          <h1 className="text-3xl font-bold font-heading">Admin</h1>
          <p className="text-muted">Manage users and roles</p>
        </div>
      </div>

      <div className="bg-card border border-border rounded-2xl overflow-hidden">
        <div className="p-4 border-b border-border flex items-center justify-between">
          <div className="flex items-center gap-2 text-foreground">
            <UsersIcon size={20} weight="bold" className="text-primary" />
            <span className="font-semibold">Users</span>
          </div>
          <span className="text-sm text-muted">{users.length} total</span>
        </div>

        {error ? (
          <div className="p-6 text-red-500">{error}</div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-muted">
                  <th className="px-4 py-3 font-medium">Username</th>
                  <th className="px-4 py-3 font-medium">Email</th>
                  <th className="px-4 py-3 font-medium">Role</th>
                  <th className="px-4 py-3 font-medium">Level</th>
                  <th className="px-4 py-3 font-medium">XP</th>
                  <th className="px-4 py-3 font-medium">Created</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id} className="border-b border-border/60 hover:bg-muted/10 transition-colors">
                    <td className="px-4 py-3 font-medium text-foreground">{user.username}</td>
                    <td className="px-4 py-3 text-muted">{user.email}</td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex rounded-full px-2.5 py-1 text-xs font-bold uppercase tracking-wide ${
                          user.role === "admin"
                            ? "bg-amber-500/15 text-amber-600"
                            : "bg-primary/10 text-primary"
                        }`}
                      >
                        {user.role}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-foreground">{user.level}</td>
                    <td className="px-4 py-3 text-foreground">{user.xp}</td>
                    <td className="px-4 py-3 text-muted">{new Date(user.createdAt).toLocaleDateString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
