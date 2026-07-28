import { createI18n } from "vue-i18n";

import ru from "./locales/ru";
import en from "./locales/en";

const savedLocale = localStorage.getItem("locale");

const locale =
    savedLocale === "ru" || savedLocale === "en"
        ? savedLocale
        : "ru";

export const i18n = createI18n({
    legacy: false,
    locale,
    fallbackLocale: "en",
    messages: {
        ru,
        en,
    },
});