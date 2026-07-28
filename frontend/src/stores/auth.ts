import { defineStore } from "pinia";
import { me } from "@/api/auth";

export interface User {
    id: number;
    username: string;
    config_limit: number;
    group_id: number | null;
    is_admin: boolean;
}

export const useAuthStore = defineStore("auth", {
    state: () => ({
        token: localStorage.getItem("token") || "",
        user: (() => {
            const raw = localStorage.getItem("user");
            return raw ? (JSON.parse(raw) as User) : null;
        })(),
        initialized: false,
    }),

    getters: {
        isAuthenticated: (state) => state.token.length > 0,
        isAdmin: (state) => state.user?.is_admin ?? false,
    },

    actions: {
        setToken(token: string) {
            this.token = token;
            localStorage.setItem("token", token);
        },

        setUser(user: User) {
            this.user = user;
            localStorage.setItem("user", JSON.stringify(user));
        },

        async restoreSession() {
            if (!this.token) {
                this.initialized = true;
                return;
            }

            try {
                const user = await me();
                this.setUser(user);
            } catch {
                this.logout();
            } finally {
                this.initialized = true;
            }
        },

        logout() {
            this.token = "";
            this.user = null;
            this.initialized = true;
            localStorage.removeItem("token");
            localStorage.removeItem("user");
        },
    },
});