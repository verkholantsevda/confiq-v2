<template>
    <aside class="sidebar">
        <div class="logo">
            Confiq
        </div>
        <nav>
            <Button
                v-for="item in menuItems"
                :key="item.key"
                :label="t(`menu.${item.key}`)"
                :icon="item.icon"
                :text="!isActive(item.to)"
                :severity="isActive(item.to) ? 'primary' : 'secondary'"
                fluid
                @click="navigate(item.to)"
            />
        </nav>

        <div class="footer">
            Confiq v0.1
        </div>
    </aside>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter, useRoute } from "vue-router";
import Button from "primevue/button";
import { useLayoutStore } from "@/stores/layout";
import { useAuthStore } from "@/stores/auth";

const { t } = useI18n({
    useScope: "global",
});

const router = useRouter();
const route = useRoute();
const layout = useLayoutStore();
const auth = useAuthStore();

const adminMenu = [
    { key: "dashboard", icon: "pi pi-home", to: "/admin" },
    { key: "users", icon: "pi pi-users", to: "/admin/users" },
    { key: "endpoints", icon: "pi pi-server", to: "/admin/endpoints" },
    { key: "groups", icon: "pi pi-sitemap", to: "/admin/groups" },
    { key: "configtypes", icon: "pi pi-tags", to: "/admin/config-types" },
    { key: "configs", icon: "pi pi-file", to: "/admin/configs" },
];

const userMenu = [
    { key: "dashboard", icon: "pi pi-home", to: "/user" },
    { key: "configs", icon: "pi pi-file", to: "/user/configs" },
];

const menuItems = computed(() =>
    auth.user?.is_admin ? adminMenu : userMenu,
);

function navigate(to: string) {
    router.push(to);
    layout.closeSidebar();
}

function isActive(to: string) {
    if (to === "/admin" || to === "/user") {
        return route.path === to;
    }

    return route.path.startsWith(to);
}
</script>

<style scoped>
.sidebar {
    width: 260px;
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--p-content-border-color);
    background: var(--p-surface-card);
}

.logo {
    padding: 24px;
    font-size: 1.3rem;
    font-weight: bold;
}

nav {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: .25rem;
    padding: 0 12px;
}

.footer {
    padding: 20px;
    opacity: .6;
    font-size: .9rem;
}

:deep(.p-button) {
    justify-content: flex-start;
}

:deep(.p-button-label) {
    text-align: left;
    flex: 1;
}

:deep(.p-button-icon) {
    margin-right: .75rem;
}
:deep(.p-button) {
    width: 100%;
    justify-content: flex-start;
    padding: .75rem 1rem;
    border-radius: 8px;
}
</style>