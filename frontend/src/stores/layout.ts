import { defineStore } from "pinia";

export const useLayoutStore = defineStore("layout", {
    state: () => ({
        sidebarVisible: false,
    }),

    actions: {
        openSidebar() {
            this.sidebarVisible = true;
        },

        closeSidebar() {
            this.sidebarVisible = false;
        },

        toggleSidebar() {
            this.sidebarVisible = !this.sidebarVisible;
        },
    },
});