<template>
    <div class="page-grid">

        <!-- Левая колонка -->
        <div class="left-column">

            <!-- Основная информация -->
            <Card>

                <template #title>
                    <div class="header-row">
                      <div>
                        <div class="config-name">{{ config?.name }}</div>
                        <div class="config-type">{{ config?.config_type?.name }}</div>
                      </div>
                      <div class="header-actions">
                        <Button
                            icon="pi pi-pencil"
                            text
                            rounded
                            @click="editConfiguration"
                        />
                        <Button
                            icon="pi pi-download"
                            text
                            rounded
                            @click="downloadConfiguration"
                        />
                        <Button
                            icon="pi pi-trash"
                            severity="danger"
                            text
                            rounded
                            @click="removeConfiguration"
                        />
                      </div>
                    </div>
                </template>

                <template #content>

                    <div class="top-info">
                      <div>
                        <div class="label">{{ t("pages.configs_edit.endpoint") }}</div>
                        <div class="endpoint-name">{{ config?.endpoint?.name }}</div>
                        <div class="endpoint-address">{{ config?.endpoint?.address }}:{{ config?.endpoint?.port }}</div>
                      </div>
                      <div>
                        <div class="label">{{ t("pages.configs_edit.created") }}</div>
                        <div class="created-date">{{ formatDate(config?.created_at) }}</div>
                      </div>
                    </div>

                </template>

            </Card>

            <!-- Конфигурация -->
            <Card>

                <template #title>
                    <div class="section-header">
                      <span>{{ t("pages.configs_edit.configs") }}</span>
                      <Button icon="pi pi-copy" :label="t('pages.configs_edit.copy')" text size="small" />
                    </div>
                </template>

                <template #content>

                    <Textarea
                        :model-value="config?.config_content ?? ''"
                        :rows="22"
                        autoResize
                        class="w-full font-mono"
                        readonly
                        style="width:100%; min-width:100%;"
                    />

                </template>

            </Card>

            <!-- Детали подключения -->
            <Card>

                <template #title>
                    {{ t("pages.configs_edit.detail") }}
                </template>

                <template #content>

                    <div class="details-grid">

                        <div class="detail-card">
                            <div class="title">IPv4 {{ t("pages.configs_edit.address") }}</div>
                            <code>{{ config?.client_ipv4 }}</code>
                        </div>

                        <div class="detail-card">
                            <div class="title">IPv6 {{ t("pages.configs_edit.address") }}</div>
                            <code>{{ config?.client_ipv6 }}</code>
                        </div>

                        <div class="detail-card">
                            <div class="title">Public Key</div>
                            <code>{{ config?.public_key }}</code>
                        </div>

                        <div class="detail-card">
                            <div class="title">Peer Public Key</div>
                            <code>{{ config?.peer_public_key }}</code>
                        </div>

                    </div>

                </template>

            </Card>

        </div>

        <!-- Правая колонка -->
        <div class="right-column">

            <!-- QR -->
            <Card>

                <template #title>
                    QR Code
                </template>

                <template #content>

                    <div class="qr-block">

                        <img
                            v-if="qrCode"
                            :src="qrCode"
                            alt="QR Code"
                            class="qr-image"
                        />
                        <div v-else class="qr-placeholder">
                            QR CODE
                        </div>

                        <p class="qr-text">
                            {{ t("pages.configs_edit.QR_Code") }}
                            <strong>{{ config?.config_type?.name }}</strong>
                        </p>

                        <Button
                            icon="pi pi-download"
                            :label="t('pages.configs_edit.download')"
                            fluid
                        />

                    </div>

                </template>

            </Card>

            <!-- Донаты -->
            <Card class="donate-card">
                <template #title>
                    ❤️ {{ t("donate.title") }}
                </template>

                <template #content>
                    <p>
                        {{ t("donate.description") }}
                    </p>

                    <div class="donate-actions">
                        <Button
                            label="Youmoney"
                            icon="pi pi-heart"
                            as="a"
                            :href="DONATE_YOOMONEY_URL"
                            target="_blank"
                            rel="noopener noreferrer"
                        />
                        <Button
                            label="CloudTip"
                            icon="pi pi-credit-card"
                            severity="secondary"
                            as="a"
                            :href="DONATE_CLOUDTIP_URL"
                            target="_blank"
                            rel="noopener noreferrer"
                        />
                    </div>
                </template>
            </Card>

            <!-- Клиенты -->
            <Card>

                <template #title>
                    {{ t("pages.configs_edit.clients") }}
                </template>

                <template #content>

                    <div class="client-links">

                        <Button
                            icon="pi pi-desktop"
                            label="Windows"
                            fluid
                        />

                        <Button
                            icon="pi pi-desktop"
                            label="MacOS"
                            fluid
                        />

                        <Button
                            icon="pi pi-mobile"
                            label="Android"
                            fluid
                        />

                        <Button
                            icon="pi pi-mobile"
                            label="iOS"
                            fluid
                        />

                    </div>

                </template>

            </Card>

        </div>

    </div>

    <!-- Инструкция -->
    <Card class="instructions">

        <template #title>
            {{ t("pages.configs_edit.instruction") }}
        </template>

        <template #content>

            <div
                class="instructions-html"
                v-html="config?.config_type?.usage_instructions"
            />

        </template>

    </Card>

</template>

<script setup lang="ts">
import Textarea from 'primevue/textarea';
import { useI18n } from "vue-i18n";
import { computed, onMounted, ref, watchEffect } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import QRCode from 'qrcode';

import Card from "primevue/card";
import Button from "primevue/button";

import PageHeader from "@/components/common/PageHeader.vue";

import type { Configuration } from "@/types/config";

import { useConfigurations } from "@/composables/useConfigs";
const { t } = useI18n({
    useScope: "global",
});
const {
    loading,
    loadConfiguration,
} = useConfigurations();

const config = ref<Configuration | null>(null);


const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

const {
    deleteConfiguration,
} = useConfigurations();

const DONATE_YOOMONEY_URL = import.meta.env.VITE_DONATE_YOOMONEY_URL ?? '#';
const DONATE_CLOUDTIP_URL = import.meta.env.VITE_DONATE_CLOUDTIP_URL ?? '#';



const qrCode = ref('');

watchEffect(async () => {
    if (!config.value?.config_content) {
        qrCode.value = '';
        return;
    }

    qrCode.value = await QRCode.toDataURL(config.value.config_content, {
        width: 260,
        margin: 1,
    });
});

onMounted(async () => {
    config.value = await loadConfiguration(
        Number(route.params.id)
    );
});

async function removeConfiguration() {
    if (!config.value) return;

    const confirmed = confirm('Удалить конфигурацию?');

    if (!confirmed) return;

    await deleteConfiguration(config.value.id);

    await router.push(
        auth.user?.is_admin
            ? '/admin/configs'
            : '/user/configs',
    );
}

function editConfiguration() {
    if (!config.value) return;

    router.push(`/configs/${config.value.id}/edit`);
}

function downloadConfiguration() {
    if (!config.value?.config_content) return;

    const blob = new Blob(
        [config.value.config_content],
        { type: 'text/plain;charset=utf-8' },
    );

    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');

    link.href = url;
    link.download = `${config.value.name}.conf`;
    link.click();

    URL.revokeObjectURL(url);
}

function formatDate(date?: string) {
    if (!date) return "-";

    return new Date(date).toLocaleString("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}
</script>
<style scoped>
.page-grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 340px;
    gap: 24px;
    align-items: start;
}

.left-column,
.right-column {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    min-width: 0;
}

.left-column {
    min-width: 0;
}

.header-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 1rem;
    flex-wrap: wrap;
}

.header-actions {
    display: flex;
    gap: .5rem;
    flex-shrink: 0;
}

.config-name {
    font-size: 1.6rem;
    font-weight: 700;
    word-break: break-word;
}

.config-type {
    margin-top: .35rem;
    display: inline-block;
    padding: .25rem .6rem;
    background: var(--p-highlight-background);
    border-radius: 8px;
}

.top-info {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 2rem;
    margin-top: 1rem;
}

.label {
    font-size: .8rem;
    color: var(--p-text-muted-color);
    margin-bottom: .35rem;
}

.endpoint-name,
.created-date {
    font-weight: 600;
}

.endpoint-address {
    margin-top: .25rem;
    color: var(--p-text-muted-color);
    font-family: monospace;
    word-break: break-all;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
}

.details-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
}

.detail-card {
    padding: .75rem;
    border: 1px solid var(--p-content-border-color);
    border-radius: 10px;
}

.detail-card .title {
    font-size: .78rem;
    color: var(--p-text-muted-color);
    margin-bottom: .35rem;
}

.detail-card code {
    display: block;
    font-family: monospace;
    font-size: .83rem;
    word-break: break-all;
}

.qr-block {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
}

.qr-placeholder {
    width: 240px;
    height: 240px;
    border: 2px dashed var(--p-content-border-color);
    display: flex;
    justify-content: center;
    align-items: center;
}

.qr-image {
    width: 100%;
    max-width: 260px;
    aspect-ratio: 1;
    object-fit: contain;
}

.qr-text {
    text-align: center;
    font-size: .85rem;
    color: var(--p-text-muted-color);
}

.client-links,
.donate-actions {
    display: flex;
    flex-direction: column;
    gap: .75rem;
}

.instructions {
    margin-top: 1.5rem;
}

.instructions-html {
    overflow-wrap: anywhere;
}

textarea {
    width: 100% !important;
    min-width: 0 !important;
}

/* ---------------- Ноутбуки ---------------- */

@media (max-width: 1400px) {

    .page-grid {
        grid-template-columns: minmax(0, 1fr) 300px;
    }

}

/* ---------------- Планшеты ---------------- */

@media (max-width: 1100px) {

    .page-grid {
        grid-template-columns: 1fr;
    }

    .right-column {
        order: -1;
    }

    .qr-image {
        max-width: 220px;
    }

}

/* ---------------- Телефоны ---------------- */

@media (max-width: 768px) {

    .config-name {
        font-size: 1.3rem;
    }

    .top-info {
        grid-template-columns: 1fr;
        gap: 1rem;
    }

    .details-grid {
        grid-template-columns: 1fr;
    }

    .header-row {
        flex-direction: column;
    }

    .header-actions {
        width: 100%;
        justify-content: flex-end;
    }

    .section-header {
        flex-direction: column;
        align-items: stretch;
    }

    .section-header .p-button {
        width: 100%;
    }

    .qr-image {
        max-width: 180px;
    }

}

/* ---------------- Очень маленькие экраны ---------------- */

@media (max-width: 480px) {

    .config-name {
        font-size: 1.15rem;
    }

    .config-type {
        font-size: .8rem;
    }

    .qr-image {
        max-width: 150px;
    }

}

</style>