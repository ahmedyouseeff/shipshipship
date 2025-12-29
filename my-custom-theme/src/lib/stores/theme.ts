import { writable } from "svelte/store";
import { browser } from "$app/environment";

export type Theme = "light" | "dark";
export type ThemePreference = "light" | "dark";

function createThemeStore() {
  const { subscribe, set, update } = writable<Theme>("light");

  let mediaQuery: MediaQueryList | null = null;
  let mediaQueryHandler: ((e: MediaQueryListEvent) => void) | null = null;

  return {
    subscribe,
    init: () => {
      if (!browser) return;

      // Always force light mode - ignore any stored preference
      const initialTheme: Theme = "light";

      // Always set localStorage to light if accessed
      localStorage.setItem("theme", "light");

      // Apply theme to document - ensure dark class is removed
      document.documentElement.classList.remove("dark");
      set(initialTheme);

      // No longer listening for system theme changes since system option is removed
    },
    toggle: () => {
      // Do nothing - always keep light mode
      // But still update the store to prevent UI inconsistencies
      if (browser) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
      }
      set("light");
    },
    setPreference: (preference: ThemePreference) => {
      if (!browser) return;

      // Always force light mode regardless of preference
      localStorage.setItem("theme", "light");
      document.documentElement.classList.remove("dark");
      set("light");
    },
    getPreference: (): ThemePreference => {
      if (!browser) return "light";
      // Always return light and ensure localStorage is set to light
      localStorage.setItem("theme", "light");
      return "light";
    },
    set: (theme: Theme) => {
      // Always force light mode regardless of what's passed
      if (browser) {
        document.documentElement.classList.remove("dark");
        localStorage.setItem("theme", "light");
      }
      set("light");
    },
  };
}

export const theme = createThemeStore();
