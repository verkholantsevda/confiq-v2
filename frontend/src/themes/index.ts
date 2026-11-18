import Aura from "@primeuix/themes/aura";
import { palette } from "@primeuix/themes";

export type ThemeMode = "light" | "dark";

export const ACCENTS = [
    { id: "green", palette: "emerald", color: "#10b981" },
    { id: "blue", palette: "blue", color: "#3b82f6" },
    { id: "violet", palette: "violet", color: "#8b5cf6" },
    { id: "orange", palette: "orange", color: "#f97316" },
    { id: "rose", palette: "rose", color: "#f43f5e" },
    { id: "cyan", palette: "cyan", color: "#06b6d4" },
] as const;

export type Accent = string;

export const DEFAULT_ACCENT: Accent = "green";

type SurfaceOption = {
    id: string;
    palette: string;
    ground: string;
    card: string;
};

export const LIGHT_SURFACES: SurfaceOption[] = [
    { id: "white", palette: "neutral", ground: "#ffffff", card: "#f4f4f5" },
    { id: "cloud", palette: "slate", ground: "#f1f5f9", card: "#f8fafc" },
    { id: "fog", palette: "gray", ground: "#e5e7eb", card: "#f3f4f6" },
    { id: "cream", palette: "stone", ground: "#f5efe4", card: "#fcf8f0" },
    { id: "ice", palette: "slate", ground: "#e0f2fe", card: "#f0f9ff" },
];

export const DARK_SURFACES: SurfaceOption[] = [
    { id: "graphite", palette: "zinc", ground: "#09090b", card: "#18181b" },
    { id: "coal", palette: "neutral", ground: "#0a0a0a", card: "#171717" },
    { id: "midnight", palette: "slate", ground: "#020617", card: "#0f172a" },
    { id: "charcoal", palette: "stone", ground: "#0c0a09", card: "#1c1917" },
    { id: "anthracite", palette: "gray", ground: "#030712", card: "#111827" },
];

export type Surface = string;

export const DEFAULT_SURFACE_LIGHT: Surface = "cloud";
export const DEFAULT_SURFACE_DARK: Surface = "graphite";

const HEX_COLOR = /^#[0-9a-f]{6}$/i;

export function normalizeAccent(value: string | null | undefined): Accent | null {
    if (!value) return null;

    if (ACCENTS.some((accent) => accent.id === value)) {
        return value;
    }

    const normalized = value.startsWith("#") ? value : `#${value}`;

    if (HEX_COLOR.test(normalized)) {
        return normalized;
    }

    return null;
}

export function getAccentPalette(accent: Accent): Record<string, string> {
    const preset = ACCENTS.find((item) => item.id === accent);

    if (preset) {
        const source = (Aura.primitive as Record<string, Record<string, string>>)[
            preset.palette
        ];

        return { ...source };
    }

    return palette(accent) as Record<string, string>;
}

function surfacesFor(mode: ThemeMode): SurfaceOption[] {
    return mode === "dark" ? DARK_SURFACES : LIGHT_SURFACES;
}

export function normalizeSurface(
    value: string | null | undefined,
    mode: ThemeMode,
): Surface | null {
    if (!value) return null;

    return surfacesFor(mode).some((surface) => surface.id === value)
        ? value
        : null;
}

export function getSurfaceColors(
    surface: Surface,
    mode: ThemeMode,
): Pick<SurfaceOption, "ground" | "card"> {
    const options = surfacesFor(mode);
    const option =
        options.find((item) => item.id === surface) ?? options[0];

    return { ground: option.ground, card: option.card };
}

export function getSurfacePalette(
    surface: Surface,
    mode: ThemeMode,
): Record<string, string> {
    const options = surfacesFor(mode);
    const option =
        options.find((item) => item.id === surface) ?? options[0];
    const source = (Aura.primitive as Record<string, Record<string, string>>)[
        option.palette
    ];

    return { ...source };
}
