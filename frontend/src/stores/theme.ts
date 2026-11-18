import { defineStore } from "pinia";
import {
    DEFAULT_ACCENT,
    DEFAULT_SURFACE_DARK,
    DEFAULT_SURFACE_LIGHT,
    getAccentPalette,
    getSurfaceColors,
    getSurfacePalette,
    normalizeAccent,
    normalizeSurface,
    type Accent,
    type Surface,
    type ThemeMode,
} from "@/themes";

export type Theme = ThemeMode;

function storedTheme(): Theme {
    return localStorage.getItem("theme") === "dark" ? "dark" : "light";
}

function storedAccent(): Accent {
    return normalizeAccent(localStorage.getItem("accent")) ?? DEFAULT_ACCENT;
}

function storedSurface(
    key: "surfaceLight" | "surfaceDark",
    mode: ThemeMode,
    fallback: Surface,
): Surface {
    return normalizeSurface(localStorage.getItem(key), mode) ?? fallback;
}

export const useThemeStore = defineStore("theme", {
    state: () => ({
        theme: storedTheme(),
        accent: storedAccent(),
        surfaceLight: storedSurface("surfaceLight", "light", DEFAULT_SURFACE_LIGHT),
        surfaceDark: storedSurface("surfaceDark", "dark", DEFAULT_SURFACE_DARK),
    }),

    getters: {
        surface: (state): Surface =>
            state.theme === "dark" ? state.surfaceDark : state.surfaceLight,
    },

    actions: {
        applyTheme() {
            const html = document.documentElement;

            if (this.theme === "dark") {
                html.classList.add("app-dark");
            } else {
                html.classList.remove("app-dark");
            }

            localStorage.setItem("theme", this.theme);
            this.applyAccent();
            this.applySurface();
        },

        setTheme(theme: Theme) {
            this.theme = theme;
            this.applyTheme();
        },

        toggleTheme() {
            this.theme = this.theme === "light" ? "dark" : "light";
            this.applyTheme();
        },

        applyAccent(accent?: Accent) {
            const normalized =
                normalizeAccent(accent ?? this.accent) ?? DEFAULT_ACCENT;
            const root = document.documentElement;

            for (const [shade, color] of Object.entries(
                getAccentPalette(normalized),
            )) {
                root.style.setProperty(`--p-primary-${shade}`, color);
            }
        },

        applySurface(surface?: Surface) {
            const normalized =
                normalizeSurface(surface ?? this.surface, this.theme) ?? this.surface;
            const root = document.documentElement;
            const palette = getSurfacePalette(normalized, this.theme);
            const colors = getSurfaceColors(normalized, this.theme);

            for (const [shade, color] of Object.entries(palette)) {
                root.style.setProperty(`--p-surface-${shade}`, color);
            }

            root.style.setProperty("--p-surface-ground", colors.ground);
            root.style.setProperty("--p-surface-card", colors.card);
        },

        previewAccent(accent: string) {
            const normalized = normalizeAccent(accent);

            if (normalized) {
                this.applyAccent(normalized);
            }
        },

        previewSurface(mode: Theme, surface: string) {
            const normalized = normalizeSurface(surface, mode);

            if (!normalized) return;

            if (mode === this.theme) {
                this.applySurface(normalized);
            }
        },

        resetAccent() {
            this.applyAccent();
        },

        resetSurface() {
            this.applySurface();
        },

        setAccent(accent: string) {
            const normalized = normalizeAccent(accent);

            if (!normalized) return;

            this.accent = normalized;
            localStorage.setItem("accent", normalized);
            this.applyAccent();
        },

        setSurface(mode: Theme, surface: string) {
            const normalized = normalizeSurface(surface, mode);

            if (!normalized) return;

            if (mode === "dark") {
                this.surfaceDark = normalized;
                localStorage.setItem("surfaceDark", normalized);
            } else {
                this.surfaceLight = normalized;
                localStorage.setItem("surfaceLight", normalized);
            }

            if (mode === this.theme) {
                this.applySurface();
            }
        },
    },
});
