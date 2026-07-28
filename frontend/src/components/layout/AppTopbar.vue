<template>
    <header class="topbar">

        <div class="left">
            <Button
                icon="pi pi-bars"
                text
                rounded
                severity="secondary"
                class="menu-button"
                @click="layout.toggleSidebar()"            />
        </div>

        <div class="actions">

            <!-- Theme -->
            <Button
                :icon="isDark ? 'pi pi-moon' : 'pi pi-sun'"
                text
                rounded
                severity="secondary"
                @click="toggleTheme"
            />

            <!-- Language -->
            <SelectButton
                v-model="language"
                :options="languages"
                optionLabel="label"
                optionValue="value"
                size="small"
            />

            <Divider layout="vertical" />

            <span v-if="username" class="username">
                {{ username }}
            </span>

            <Button
                icon="pi pi-sign-out"
                text
                rounded
                severity="danger"
                @click="logout"
            />

        </div>

    </header>
</template>

<script setup lang="ts">
import { computed } from "vue";

import Button from "primevue/button";
import Divider from "primevue/divider";
import SelectButton from "primevue/selectbutton";

import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";

import { useAuthStore } from "@/stores/auth";
import { useThemeStore } from "@/stores/theme";
import { useLocaleStore } from "@/stores/locale";
import { useLayoutStore } from "@/stores/layout";

const layout = useLayoutStore();
const theme = useThemeStore();
const locale = useLocaleStore();

const { t } = useI18n({
    useScope: "global",
});


const isDark = computed(() => theme.theme === "dark");

const languages = [
    { label: "RU", value: "ru" },
    { label: "EN", value: "en" },
];

const language = computed({
    get: () => locale.locale,
    set: (value) => locale.setLocale(value),
});

function toggleTheme() {
    theme.toggleTheme();
}

const router = useRouter();
const auth = useAuthStore();

const username = computed(() => auth.user?.username ?? '');

function logout() {
    auth.logout();
    router.push("/login");
}
</script>

<style scoped>
.topbar {
    height: 64px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    padding: 0 20px;

    border-bottom: 1px solid var(--p-content-border-color);
    background: var(--p-surface-card);
}

.actions {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.username {
    font-weight: 500;
}
.left {
    display: flex;
    align-items: center;
}

.menu-button {
    display: none;
}
@media (max-width: 1024px) {

    .menu-button {
        display: inline-flex;
    }

}
</style>