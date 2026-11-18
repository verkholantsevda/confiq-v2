<template>
    <PageHeader
        :title="t('user.profile.title')"
        :subtitle="t('user.profile.subtitle')"
    />

    <div class="profile-grid">

        <Card>
            <template #title>
                <label>{{ t("user.profile.title") }}</label>
            </template>

            <template #content>

                <div class="info-section">
                    <div class="label"><label>{{ t("user.profile.user") }}</label></div>
                    <div class="value">{{ me?.username ?? auth.user?.username }}</div>
                </div>

                <Divider />

                <div class="info-section">
                    <div><label>{{ t("user.profile.created") }}</label></div>
                    <div class="value">
                        {{ formatDate(me?.created_at) }}
                    </div>
                </div>

                <Divider />

                <div class="info-section">
                    <div class="label"><label>{{ t("user.profile.conf") }}</label></div>
                    <div class="value">
                        {{ usedConfigs }} / {{ me?.config_limit ?? configLimit }}
                    </div>
                </div>

            </template>
        </Card>

        <div class="right-column">

            <Card>
                <template #title>
                    <label>{{ t("user.profile.appearance.title") }}</label>
                </template>

                <template #content>

                    <div class="field">
                        <label>{{ t("user.profile.appearance.theme") }}</label>

                        <SelectButton
                            v-model="themeMode"
                            :options="themeOptions"
                            optionLabel="label"
                            optionValue="value"
                            :allowEmpty="false"
                        />
                    </div>

                    <div class="field">
                        <label>{{ t("user.profile.appearance.accent") }}</label>

                        <div class="accent-swatches">
                            <button
                                v-for="option in accents"
                                :key="option.id"
                                type="button"
                                class="accent-swatch"
                                :class="{ selected: draftAccent === option.id }"
                                :style="{ backgroundColor: option.color }"
                                :title="option.label"
                                :aria-label="option.label"
                                :aria-pressed="draftAccent === option.id"
                                @click="selectAccent(option.id)"
                            >
                                <i v-if="draftAccent === option.id" class="pi pi-check" />
                            </button>
                        </div>

                        <div class="custom-accent">
                            <ColorPicker
                                v-model="customAccent"
                                format="hex"
                                :aria-label="t('user.profile.appearance.custom')"
                            />
                            <span>{{ t("user.profile.appearance.custom") }}</span>
                        </div>

                        <small class="hint">
                            {{ t("user.profile.appearance.accent_hint") }}
                        </small>
                    </div>

                    <div class="field">
                        <label>{{ t("user.profile.appearance.surface") }}</label>

                        <div class="surface-options">
                            <button
                                v-for="option in surfaces"
                                :key="option.id"
                                type="button"
                                class="surface-option"
                                :class="{ selected: draftSurface === option.id }"
                                :aria-pressed="draftSurface === option.id"
                                @click="selectSurface(option.id)"
                            >
                                <span
                                    class="surface-dot"
                                    :style="{ backgroundColor: option.ground }"
                                ></span>
                                <span>{{ option.label }}</span>
                            </button>
                        </div>

                        <small class="hint">
                            {{ t("user.profile.appearance.surface_hint") }}
                        </small>
                    </div>

                    <Button
                        class="appearance-save"
                        :label="t('common.save')"
                        icon="pi pi-check"
                        size="small"
                        :disabled="!settingsChanged"
                        @click="saveAppearance"
                    />

                    <small v-if="settingsChanged" class="hint unsaved">
                        {{ t("user.profile.appearance.unsaved") }}
                    </small>

                </template>
            </Card>

            <Card>
                <template #title>
                    <label>{{ t("user.profile.password.title") }}</label>
                </template>

                <template #content>

                    <form
                        class="password-form"
                        @submit.prevent="changePassword"
                    >

                        <div class="field">
                            <label>{{ t("user.profile.password.current") }}</label>

                            <Password
                                v-model="form.currentPassword"
                                toggleMask
                                fluid
                                :feedback="false"
                            />
                        </div>

                        <div class="field">
                            <label>{{ t("user.profile.password.new") }}</label>

                            <Password
                                v-model="form.newPassword"
                                toggleMask
                                fluid
                                :feedback="false"
                            />

                            <small class="hint">
                                <label>{{ t("user.profile.password.new_hint") }}</label>
                            </small>
                        </div>

                        <div class="field">
                            <label>{{ t("user.profile.password.confirm") }}</label>

                            <Password
                                v-model="form.confirmPassword"
                                toggleMask
                                fluid
                                :feedback="false"
                            />
                        </div>

                        <Button
                            :label="t('user.profile.password.change')"
                            icon="pi pi-check"
                            type="submit"
                            :loading="loading"
                            :disabled="loading"
                            fluid
                        />

                    </form>

                </template>
            </Card>
            <Card v-if="totpAvailable">
                <template #title>
                    <label>Двухфакторная аутентификация</label>
                </template>

                <template #content>

                    <div v-if="!totpEnabled">

                        <p>
                            Включите защиту аккаунта через приложение-аутентификатор.
                        </p>

                        <Button
                            v-if="!totpSecret"
                            label="Создать ключ"
                            icon="pi pi-shield"
                            @click="setupTwoFactor"
                        />

                        <div v-else class="totp-setup">

                            <label>
                                Добавьте этот ключ в приложение:
                            </label>

                            <div class="qr-wrapper">
                                <QrcodeVue :value="totpUri" :size="180" level="M" />
                            </div>

                            <code>{{ totpSecret }}</code>

                            <InputText
                                v-model="totpCode"
                                placeholder="123456"
                            />

                            <Button
                                label="Подтвердить"
                                icon="pi pi-check"
                                :loading="totpLoading"
                                @click="confirmTwoFactor"
                            />

                        </div>

                    </div>


                    <div v-else>

                        <p>
                            Двухфакторная аутентификация включена
                        </p>

                        <Button
                            label="Отключить"
                            severity="danger"
                            icon="pi pi-times"
                            @click="removeTwoFactor"
                        />

                    </div>

                </template>
            </Card>

        </div>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, onUnmounted, computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import Card from "primevue/card";
import ColorPicker from "primevue/colorpicker";
import Divider from "primevue/divider";
import Button from "primevue/button";
import Password from "primevue/password";
import InputText from "primevue/inputtext";
import SelectButton from "primevue/selectbutton";
import QrcodeVue from "qrcode.vue";

import PageHeader from "@/components/common/PageHeader.vue";

import { useAuthStore } from "@/stores/auth";
import { useThemeStore, type Theme } from "@/stores/theme";
import {
    ACCENTS,
    DARK_SURFACES,
    LIGHT_SURFACES,
    getAccentPalette,
    normalizeAccent,
} from "@/themes";
import { useConfigsView } from "@/composables/useConfigsView";

import { getMe, changePassword as changePasswordApi } from "@/api/me";
import type { User } from "@/types/user";
import {
    getTOTPStatus,
    setupTOTP,
    enableTOTP,
    disableTOTP,
} from "@/api/totp";

const auth = useAuthStore();
const theme = useThemeStore();
const me = ref<User | null>(null);
const { t } = useI18n({
    useScope:"global"
});

const themeMode = computed({
    get: () => theme.theme,
    set: (value: Theme) => theme.setTheme(value),
});

const themeOptions = computed(() => [
    { label: t("user.profile.appearance.light"), value: "light" },
    { label: t("user.profile.appearance.dark"), value: "dark" },
]);

const accents = computed(() =>
    ACCENTS.map((accent) => ({
        ...accent,
        label: t(`user.profile.appearance.accents.${accent.id}`),
    })),
);

const draftAccent = ref(theme.accent);
const draftSurface = ref(theme.surface);

const accentChanged = computed(() => draftAccent.value !== theme.accent);
const surfaceChanged = computed(() => draftSurface.value !== theme.surface);
const settingsChanged = computed(
    () => accentChanged.value || surfaceChanged.value,
);

const customAccent = computed({
    get: () => getAccentPalette(draftAccent.value)["500"].replace("#", ""),
    set: (value: string) => selectAccent(value),
});

const surfaces = computed(() => {
    const options = theme.theme === "dark" ? DARK_SURFACES : LIGHT_SURFACES;

    return options.map((surface) => ({
        ...surface,
        label: t(`user.profile.appearance.surfaces.${surface.id}`),
    }));
});

function selectAccent(value: string) {
    const normalized = normalizeAccent(value) ?? value;

    draftAccent.value = normalized;
    theme.previewAccent(normalized);
}

function selectSurface(value: string) {
    draftSurface.value = value;
    theme.previewSurface(theme.theme, value);
}

function saveAppearance() {
    if (accentChanged.value) {
        theme.setAccent(draftAccent.value);
    }

    if (surfaceChanged.value) {
        theme.setSurface(theme.theme, draftSurface.value);
    }
}

watch(
    () => theme.theme,
    (mode) => {
        draftSurface.value = mode === "dark" ? theme.surfaceDark : theme.surfaceLight;
    },
);

onUnmounted(() => {
    if (accentChanged.value) {
        theme.resetAccent();
    }

    if (surfaceChanged.value) {
        theme.resetSurface();
    }
});
const {
    usedConfigs,
    configLimit,
    load,
} = useConfigsView();

const loading = ref(false);
const totpEnabled = ref(false);
const totpAvailable = ref(false);
const totpSecret = ref("");
const totpCode = ref("");
const totpLoading = ref(false);

const form = reactive({
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
});

const totpUri = computed(() => me.value && totpSecret.value ? `otpauth://totp/Confiq:${encodeURIComponent(me.value.username)}?secret=${totpSecret.value}&issuer=Confiq` : "");

onMounted(async () => {
    me.value = await getMe();
    await load();

    const totpStatus = await getTOTPStatus();

    totpEnabled.value = totpStatus.enabled;
    totpAvailable.value = totpStatus.available;
});

function formatDate(value?: string) {
    if (!value) return "-";

    return new Date(value).toLocaleString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

async function changePassword() {

    if (!form.currentPassword) {
        alert("Введите текущий пароль");
        return;
    }

    if (form.newPassword.length < 3) {
        alert("Пароль должен содержать минимум 3 символа");
        return;
    }

    if (form.newPassword !== form.confirmPassword) {
        alert("Пароли не совпадают");
        return;
    }

    loading.value = true;

    try {

        await changePasswordApi(form.currentPassword, form.newPassword);

        form.currentPassword = "";
        form.newPassword = "";
        form.confirmPassword = "";

        alert("Пароль успешно изменен");

    } catch (e) {
        console.error(e);
        alert((e as any)?.response?.data ?? "Не удалось изменить пароль");
    } finally {
        loading.value = false;
    }

}

async function setupTwoFactor() {

    try {
        const result = await setupTOTP();

        totpSecret.value = result.secret;

    } catch (e) {
        console.error(e);
        alert("Не удалось создать TOTP");
    }
}


async function confirmTwoFactor() {

    if (!totpCode.value) {
        alert("Введите код");
        return;
    }

    totpLoading.value = true;

    try {

        await enableTOTP(totpCode.value);

        totpEnabled.value = true;
        totpSecret.value = "";
        totpCode.value = "";

        alert("2FA включена");

    } catch (e) {
        console.error(e);
        alert("Неверный код");
    } finally {
        totpLoading.value = false;
    }
}


async function removeTwoFactor() {

    try {

        await disableTOTP();

        totpEnabled.value = false;

        alert("2FA отключена");

    } catch (e) {
        console.error(e);
        alert("Не удалось отключить 2FA");
    }
}
</script>

<style scoped>
.profile-grid {
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1.5rem;
    align-items: start;
}

.info-section {
    display: flex;
    flex-direction: column;
    gap: .35rem;
}

.label {
    font-size: .82rem;
    color: var(--p-text-muted-color);
}

.value {
    font-size: 1rem;
    font-weight: 600;
    word-break: break-word;
}

.password-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.field {
    display: flex;
    flex-direction: column;
    gap: .45rem;
}

.accent-swatches {
    display: flex;
    flex-wrap: wrap;
    gap: .6rem;
}

.accent-swatch {
    width: 28px;
    height: 28px;
    padding: 0;
    border: 0;
    border-radius: 50%;
    color: #fff;
    font-size: .75rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: transform .15s ease, box-shadow .15s ease;
}

.accent-swatch:hover {
    transform: scale(1.08);
}

.accent-swatch.selected {
    box-shadow: 0 0 0 2px var(--p-surface-card), 0 0 0 4px var(--p-text-color);
}

.accent-swatch:focus-visible {
    outline: 2px solid var(--p-primary-color);
    outline-offset: 3px;
}

.custom-accent {
    display: flex;
    align-items: center;
    gap: .75rem;
    margin-top: .75rem;
}

.surface-options {
    display: flex;
    flex-wrap: wrap;
    gap: .5rem;
}

.surface-option {
    display: inline-flex;
    align-items: center;
    gap: .5rem;
    padding: .4rem .65rem;
    border: 1px solid var(--p-content-border-color);
    border-radius: 8px;
    background: var(--p-surface-ground);
    color: var(--p-text-color);
    font: inherit;
    font-size: .85rem;
    cursor: pointer;
}

.surface-option:hover {
    border-color: var(--p-primary-color);
}

.surface-option.selected {
    border-color: var(--p-primary-color);
    box-shadow: inset 0 0 0 1px var(--p-primary-color);
}

.surface-option:focus-visible {
    outline: 2px solid var(--p-primary-color);
    outline-offset: 2px;
}

.surface-dot {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    border: 1px solid var(--p-content-border-color);
}

.appearance-save {
    align-self: flex-start;
    margin-top: 1rem;
}

.hint {
    color: var(--p-text-muted-color);
    font-size: .8rem;
}

.hint.unsaved {
    color: var(--p-primary-color);
}

@media (max-width: 900px) {
    .profile-grid {
        grid-template-columns: 1fr;
    }
}

.right-column {
    display:flex;
    flex-direction:column;
    gap:1.5rem;
}

.totp-setup {
    display:flex;
    flex-direction:column;
    gap:1rem;
}

.qr-wrapper {
    display:flex;
    justify-content:center;
    padding:.5rem 0;
}

.totp-setup code {
    padding:.75rem;
    border-radius:8px;
    background:var(--p-surface-ground);
    word-break:break-all;
}
</style>