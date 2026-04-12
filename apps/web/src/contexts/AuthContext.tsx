"use client";

import React, { createContext, useState, useEffect, ReactNode, useMemo, useCallback } from "react";
import authService from "@/services/authService";
import { User } from "@/types/user";
import { AppRole, canAccessPath } from "@/lib/access";

interface AuthContextType {
  user: User | null;
  token: string | null;
  role: AppRole;
  isAdmin: boolean;
  isMember: boolean;
  isGuest: boolean;
  canAccess: (pathname: string) => boolean;
  login: (token: string) => void;
  logout: () => void;
  loading: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    if (storedToken) {
      setToken(storedToken);
      return;
    }

    setLoading(false);
  }, []);

  useEffect(() => {
    if (!token) {
      return;
    }

    const fetchUser = async () => {
      try {
        const userData = await authService.getMe();
        setUser(userData);
      } catch (error) {
        console.error("Failed to fetch user", error);
        localStorage.removeItem("token");
        setToken(null);
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [token]);

  const login = (newToken: string) => {
    localStorage.setItem("token", newToken);
    setToken(newToken);
    setLoading(true);
  };

  const logout = () => {
    localStorage.removeItem("token");
    setToken(null);
    setUser(null);
    setLoading(false);
  };

  const role: AppRole = user?.role ?? "guest";
  const isAdmin = role === "admin";
  const isMember = role === "member" || role === "admin";
  const isGuest = role === "guest";

  const canAccess = useCallback(
    (pathname: string) => canAccessPath(role, pathname),
    [role]
  );

  const value = useMemo(
    () => ({
      user,
      token,
      role,
      isAdmin,
      isMember,
      isGuest,
      canAccess,
      login,
      logout,
      loading,
    }),
    [user, token, role, isAdmin, isMember, isGuest, canAccess, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
