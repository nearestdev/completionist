import Link from "next/link";
import { GithubLogo, TwitterLogo, DiscordLogo } from "@phosphor-icons/react";

export default function Footer() {
  return (
    <footer className="w-full py-8 mt-auto border-t border-border bg-card">
      <div className="container mx-auto px-6">
        <div className="flex flex-col md:flex-row justify-between items-center gap-6">
          
          {/* Logo & Copyright */}
          <div className="flex flex-col items-center md:items-start gap-2">
            <h3 className="text-lg font-heading font-bold text-foreground flex items-center gap-2">
              <span className="text-primary">Completionist</span>
            </h3>
            <p className="text-sm text-muted">
              © {new Date().getFullYear()} Completionist. All rights reserved.
            </p>
          </div>

          {/* Links */}
          <div className="flex flex-wrap justify-center gap-8 text-sm font-medium text-muted">
            <Link href="/about" className="hover:text-primary transition-colors">
              About
            </Link>
            <Link href="/privacy" className="hover:text-primary transition-colors">
              Privacy
            </Link>
            <Link href="/terms" className="hover:text-primary transition-colors">
              Terms
            </Link>
            <Link href="/contact" className="hover:text-primary transition-colors">
              Contact
            </Link>
          </div>

          {/* Socials */}
          <div className="flex items-center gap-4">
            <a 
              href="https://github.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-muted hover:text-primary transition-colors p-2 hover:bg-primary/10 rounded-full"
              aria-label="GitHub"
            >
              <GithubLogo size={20} weight="fill" />
            </a>
            <a 
              href="https://twitter.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-muted hover:text-primary transition-colors p-2 hover:bg-primary/10 rounded-full"
              aria-label="Twitter"
            >
              <TwitterLogo size={20} weight="fill" />
            </a>
            <a 
              href="https://discord.com" 
              target="_blank" 
              rel="noopener noreferrer"
              className="text-muted hover:text-primary transition-colors p-2 hover:bg-primary/10 rounded-full"
              aria-label="Discord"
            >
              <DiscordLogo size={20} weight="fill" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}