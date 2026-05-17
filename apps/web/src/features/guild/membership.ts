import { findGuildBySlug } from "./guildMaster";

const currentGuildSlugKey = "lang-war.guild.current-slug";
const everJoinedGuildKey = "lang-war.guild.ever-joined";

export function getSelectedGuildSlug(): string | null {
  try {
    return window.localStorage.getItem(currentGuildSlugKey);
  } catch {
    return null;
  }
}

export function hasSelectedGuildMembership(): boolean {
  return findGuildBySlug(getSelectedGuildSlug()) !== null;
}

export function hasEverJoinedGuild(): boolean {
  try {
    return window.localStorage.getItem(everJoinedGuildKey) === "true";
  } catch {
    return false;
  }
}

export function selectGuildMembership(slug: string): void {
  if (!findGuildBySlug(slug)) {
    return;
  }

  try {
    window.localStorage.setItem(currentGuildSlugKey, slug);
    window.localStorage.setItem(everJoinedGuildKey, "true");
  } catch {
    // The selection can still proceed visually even if storage is unavailable.
  }
}
