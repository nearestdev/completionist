export type AppRole = "guest" | "user" | "admin";

export function isGuestPublicPath(pathname: string): boolean {
  if (pathname === "/" || pathname === "/search") {
    return true;
  }

  return pathname.startsWith("/search/");
}

export function canAccessPath(role: AppRole, pathname: string): boolean {
  if (role === "guest") {
    return isGuestPublicPath(pathname);
  }

  return true;
}
