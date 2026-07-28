<template>
    <div class="layout">
        <AppSidebar class="desktop-sidebar" />

        <Drawer v-model:visible="layout.sidebarVisible" position="left" class="mobile-drawer">
            <AppSidebar />
        </Drawer>

        <div class="main">
            <AppTopbar />

            <div class="content">
                <RouterView />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { RouterView } from "vue-router";

import AppSidebar from "@/components/layout/AppSidebar.vue";
import AppTopbar from "@/components/layout/AppTopbar.vue";

import Drawer from "primevue/drawer";
import { useLayoutStore } from "@/stores/layout";

const layout = useLayoutStore();
</script>

<style scoped>
.layout {
    display: flex;
    height: 100vh;
    background: var(--p-surface-ground);
}

.main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.content {
    flex: 1;
    overflow: auto;
    padding: 24px;
}

.desktop-sidebar {
    display: block;
}

.mobile-drawer {
    width: 280px;
}

@media (max-width: 1024px) {
    .desktop-sidebar {
        display: none;
    }

    .content {
        padding: 16px;
    }
}
</style>