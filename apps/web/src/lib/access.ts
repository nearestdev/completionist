export type AppRole = "guest" | "user" | "member" | "admin";

export function isGuestPublicPath(pathname: string): boolean {
  if (pathname === "/" || pathname === "/search") {
    return true;
  }

  return pathname.startsWith("/search/") || pathname.startsWith("/franchises");
}

export function canAccessPath(role: AppRole, pathname: string): boolean {
  if (role === "guest") {
    return isGuestPublicPath(pathname);
  }

  if (pathname.startsWith("/admin")) {
    return role === "admin";
  }

  return true;
}

export function isMemberOrAbove(role: AppRole): boolean {
  return role === "member" || role === "admin";
}
