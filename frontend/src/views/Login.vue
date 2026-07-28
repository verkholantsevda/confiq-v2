<template>
    <div class="login-container">

        <div class="toolbar">

            <Button
                text
                rounded
                severity="secondary"
                :icon="isDark ? 'pi pi-moon' : 'pi pi-sun'"
                @click="toggleTheme"
            />

            <SelectButton
                v-model="language"
                :options="languages"
                optionLabel="label"
                optionValue="value"
            />

        </div>

        <div class="login-card">

            <div class="logo">
                Confiq
            </div>

            <div class="subtitle">
                {{ t("auth.loginTitle") }}
            </div>

            <form @submit.prevent="onLogin">

                <div class="field">

                    <label>
                        {{ t("auth.username") }}
                    </label>

                    <input
                        v-model="username"
                        type="text"
                        autocomplete="username"
                        required
                    />

                </div>

                <div class="field">

                    <label>
                        {{ t("auth.password") }}
                    </label>

                    <input
                        v-model="password"
                        type="password"
                        autocomplete="current-password"
                        required
                    />

                </div>

                <Button
                    class="login-button"
                    type="submit"
                    :loading="loading"
                    :label="loading ? t('common.loading') : t('common.login')"
                    fluid
                />

                <p
                    v-if="error"
                    class="error"
                >
                    {{ error }}
                </p>

            </form>

            <div class="version">
                Confiq v0.1
            </div>

        </div>

    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";

import { login, me } from "@/api/auth";
import { useAuthStore } from "@/stores/auth";
import { computed } from "vue";
import Button from "primevue/button";
import SelectButton from "primevue/selectbutton";

import { useThemeStore } from "@/stores/theme";
import { useLocaleStore } from "@/stores/locale";
import { useI18n } from "vue-i18n";

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

const username = ref("");
const password = ref("");

const loading = ref(false);
const error = ref("");

async function onLogin() {
    error.value = "";
    loading.value = true;

    try {
        const result = await login({
            username: username.value,
            password: password.value,
        });

        auth.setToken(result.token);

        const user = await me();
        auth.setUser(user);

        if (user.is_admin) {
            await router.push("/admin");
        } else {
            await router.push("/user");
        }
    } catch (e) {
        console.error(e);
        error.value = "Не удалось выполнить вход";
    } finally {
        loading.value = false;
    }
}
</script>

<style scoped>

.login-container {

    min-height: 100vh;

    display: flex;
    justify-content: center;
    align-items: center;

    background: var(--p-surface-ground);

    padding: 2rem;

}

.toolbar {

    position: absolute;

    top: 20px;
    right: 20px;

    display: flex;
    align-items: center;
    gap: .75rem;

}

.login-card {

    width: 100%;
    max-width: 420px;

    background: var(--p-surface-card);

    border: 1px solid var(--p-content-border-color);

    border-radius: 16px;

    padding: 2rem;

    box-shadow: var(--p-overlay-modal-shadow);

}

.logo {

    font-size: 2rem;

    font-weight: 700;

    text-align: center;

    margin-bottom: .5rem;

}

.subtitle {

    text-align: center;

    color: var(--p-text-muted-color);

    margin-bottom: 2rem;

}

.field {

    margin-bottom: 1.25rem;

}

.field label {

    display: block;

    margin-bottom: .5rem;

    font-weight: 600;

}

.field input {

    width: 100%;

    padding: .8rem 1rem;

    border-radius: 10px;

    border: 1px solid var(--p-content-border-color);

    background: var(--p-surface-ground);

    color: var(--p-text-color);

    outline: none;

    transition: .2s;

    box-sizing: border-box;

}

.field input:focus {

    border-color: var(--p-primary-color);

}

.login-button {

    margin-top: .5rem;

}

.error {

    color: var(--p-red-500);

    margin-top: 1rem;

    text-align: center;

}

.version {

    margin-top: 2rem;

    text-align: center;

    font-size: .85rem;

    color: var(--p-text-muted-color);

}

@media (max-width: 640px) {

    .login-card {

        padding: 1.5rem;

    }

    .toolbar {

        top: 12px;
        right: 12px;

    }

}

</style>