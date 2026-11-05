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
                :allowEmpty="false"
                size="small"
            />

            <Divider layout="vertical" />

            <button
                v-if="username"
                type="button"
                class="username"
                :title="t('user.quickactions.profile')"
                :aria-label="t('user.quickactions.profile')"
                @click="toggleUserMenu"
            >
                <i class="pi pi-user"></i>
                <span>{{ username }}</span>
                <i class="pi pi-chevron-down"></i>
            </button>

            <Menu ref="userMenu" :model="menuItems" popup />

        </div>

    </header>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";

import Button from "primevue/button";
import Divider from "primevue/divider";
import Menu from "primevue/menu";
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

const userMenu = ref<InstanceType<typeof Menu> | null>(null);

const menuItems = computed(() => [
    {
        label: t("user.quickactions.profile"),
        icon: "pi pi-user",
        command: openProfile,
    },
    {
        label: t("common.logout"),
        icon: "pi pi-sign-out",
        command: logout,
    },
]);

function openProfile() {
    router.push(auth.isAdmin ? "/admin/profile" : "/user/profile");
}

function logout() {
    auth.logout();
    router.push("/login");
}

function toggleUserMenu(event: Event) {
    userMenu.value?.toggle(event);
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
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    font: inherit;
    font-weight: 500;
    color: inherit;
    background: none;
    border: 0;
    padding: 0;
    cursor: pointer;
}

.username:hover {
    color: var(--p-primary-color);
}

.username:focus-visible {
    outline: 2px solid var(--p-primary-color);
    outline-offset: 2px;
    border-radius: 4px;
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