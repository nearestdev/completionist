"use client";

import { useTheme } from "next-themes";
import { Sun, Moon } from "@phosphor-icons/react";
import { useSyncExternalStore } from "react";

export default function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();
  const mounted = useSyncExternalStore(
    () => () => {},
    () => true,
    () => false,
  );

  const isDark = theme === "dark";

  const toggleTheme = () => {
    setTheme(isDark ? "light" : "dark");
  };

  if (!mounted) {
    return (
      <button
        className="p-2.5 rounded-full bg-gray-200 dark:bg-gray-700 opacity-50 cursor-wait"
        aria-hidden="true"
      />
    );
  }

  return (
    <button
      onClick={toggleTheme}
      className={`
        p-2.5 rounded-full transition-all duration-300
        ${isDark 
          ? "bg-gray-800 text-yellow-400 hover:bg-gray-700 ring-1 ring-gray-700" 
          : "bg-indigo-50 text-indigo-600 hover:bg-indigo-100 ring-1 ring-indigo-100"
        }
      `}
      aria-label={isDark ? "Switch to Light Mode" : "Switch to Dark Mode"}
      title={isDark ? "Switch to Light Mode" : "Switch to Dark Mode"}
    >
      <div className={`transition-transform duration-500 ${isDark ? "rotate-90" : "rotate-0"}`}>
        {isDark ? <Moon size={20} weight="fill" /> : <Sun size={20} weight="fill" />}
      </div>
    </button>
  );
}
