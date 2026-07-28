import { createApp } from "vue";
import { createPinia } from "pinia";
import PrimeVue from "primevue/config";
import Aura from "@primeuix/themes/aura";
import { useThemeStore } from "@/stores/theme";
import { useAuthStore } from "@/stores/auth";
import "primeicons/primeicons.css";
import "@/assets/theme.css";
import { i18n } from "@/i18n";
import App from "./App.vue";
import router from "./router";

async function bootstrap() {
    const app = createApp(App);

    const pinia = createPinia();
    app.use(i18n);
    app.use(pinia);

    const auth = useAuthStore(pinia);
    await auth.restoreSession();

    app.use(router);

    const theme = useThemeStore(pinia);
    theme.applyTheme();

    app.use(PrimeVue, {
        theme: {
            preset: Aura,
            options: {
                darkModeSelector: ".app-dark",
            },
        },
    });

    app.mount("#app");
}

bootstrap();