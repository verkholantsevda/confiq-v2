import { defineStore } from "pinia";

export type Theme = "light" | "dark";

export const useThemeStore = defineStore("theme", {
    state: () => ({
        theme: (localStorage.getItem("theme") as Theme) || "light",
    }),

    actions: {
        applyTheme() {
            const html = document.documentElement;

            if (this.theme === "dark") {
                html.classList.add("app-dark");
            } else {
                html.classList.remove("app-dark");
            }

            localStorage.setItem("theme", this.theme);
        },

        setTheme(theme: Theme) {
            this.theme = theme;
            this.applyTheme();
        },

        toggleTheme() {
            this.theme = this.theme === "light" ? "dark" : "light";
            this.applyTheme();
        },
    },
});