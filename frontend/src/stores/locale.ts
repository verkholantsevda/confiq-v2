import { defineStore } from "pinia";
import { i18n } from "@/i18n";

export type Locale = "ru" | "en";

export const useLocaleStore = defineStore("locale", {
    state: () => ({
        locale: (localStorage.getItem("locale") as Locale) || "ru",
    }),

    actions: {
        setLocale(locale: Locale) {
            this.locale = locale;

            localStorage.setItem("locale", locale);

            i18n.global.locale.value = locale;
        },
    },
});