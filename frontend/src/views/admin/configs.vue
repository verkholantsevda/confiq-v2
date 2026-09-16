<template>
    <PageHeader
        :title="t('menu.configs')"
        :subtitle="t('pages.configs.description')"
        :button-label="t('pages.configs.create')"
        @create="createConfig"
    />

    <div class="stats-grid">
        <Card>
            <template #title>{{ t("pages.configs.total") }}</template>
            <template #content>
                <div class="stat-value">{{ totalConfigs }}</div>
            </template>
        </Card>
    </div>

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

    <div class="configs-toolbar">
        <span class="p-input-icon-left search-field">
            <i class="pi pi-search" />
            <InputText
                v-model="search"
                :placeholder="t('common.search')"
                class="w-full"
            />
        </span>

        <div class="my-configs-toggle">
            <ToggleSwitch v-model="showOnlyMine" inputId="showOnlyMine" />
            <label for="showOnlyMine">Мои конфигурации</label>
        </div>
    </div>

    <div class="configs-grid">
        <Card v-for="config in filteredConfigs" :key="config.id">
            <template #title>
                {{ config.name }}
            </template>

            <template #subtitle>
                {{ config.config_type?.name ?? '-' }}
            </template>

            <template #content>
                <div class="config-info">
                    <div class="info-row"><i class="pi pi-map-marker"></i><span>{{ config.endpoint?.name ?? '-' }}</span></div>
                    <div class="info-row"><i class="pi pi-user"></i><span>{{ usernamesById.get(config.user_id) ?? '-' }}</span></div>
                    <div class="info-row"><i class="pi pi-globe"></i><span>{{ config.endpoint ? `${config.endpoint.address}:${config.endpoint.port}` : '-' }}</span></div>
                    <div class="info-row"><i class="pi pi-clock"></i><span>{{ new Date(config.created_at).toLocaleString('ru-RU') }}</span></div>
                    <Divider />
                    <div class="config-actions">
                        <Button
                            icon="pi pi-eye"
                            rounded
                            severity="contrast"
                            @click="viewConfig(config.id)"
                        />
                        <Button
                            icon="pi pi-qrcode"
                            rounded
                            severity="secondary"
                            @click="showQrCode(config)"
                        />
                        <Button
                            icon="pi pi-download"
                            rounded
                            @click="downloadConfig(config)"
                        />
                    </div>
                    <Divider />
                    <div class="ip-grid">
                      <div class="ip-card">
                        <div class="ip-title">IPv4</div>
                        <div class="ip-value">{{ config.client_ipv4 }}</div>
                      </div>
                      <div class="ip-card">
                        <div class="ip-title">IPv6</div>
                        <div class="ip-value">{{ config.client_ipv6 }}</div>
                      </div>
                    </div>
                </div>
            </template>
        </Card>
    </div>
    <ConfigCreateDialog
        v-model="createDialogVisible"
        @created="load"
    />
    <ConfigQrCodeDialog
        v-model="qrDialogVisible"
        :config="selectedConfig"
    />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useI18n } from "vue-i18n";
import { useRouter } from 'vue-router';

import Button from "primevue/button";
import ToggleSwitch from "primevue/toggleswitch";
import Divider from "primevue/divider";
import PageHeader from "@/components/common/PageHeader.vue";
import Card from "primevue/card";
import InputText from 'primevue/inputtext';
import ConfigQrCodeDialog from "@/components/configs/ConfigQrCodeDialog.vue";

import { useConfigsView } from '@/composables/useAdminConfigsView';
import ConfigCreateDialog from "@/components/configs/ConfigCreateDialog.vue";
import { getUsers } from "@/api/users";
const DONATE_YOOMONEY_URL = import.meta.env.VITE_DONATE_YOOMONEY_URL ?? '#';
const DONATE_CLOUDTIP_URL = import.meta.env.VITE_DONATE_CLOUDTIP_URL ?? '#';

const { t } = useI18n({
    useScope: "global",
});

const router = useRouter();

const {
    configs,
    totalConfigs,
    showOnlyMine,
    load,
} = useConfigsView();

const search = ref('');
const users = ref<Array<{ id: number; username: string }>>([]);

const usernamesById = computed(() => {
    return new Map(users.value.map(user => [user.id, user.username]));
});

const filteredConfigs = computed(() => {
    const q = search.value.trim().toLowerCase();

    if (!q) return configs.value;

    return configs.value.filter(config => {
        const configName = config.name?.toLowerCase() ?? '';
        const endpointName = config.endpoint?.name?.toLowerCase() ?? '';
        const username = usernamesById.value.get(config.user_id)?.toLowerCase() ?? '';

        return (
            configName.includes(q) ||
            endpointName.includes(q) ||
            username.includes(q)
        );
    });
});

watch(showOnlyMine, async () => {
    await load();
});

onMounted(async () => {
    users.value = await getUsers();
    await load();
});

function viewConfig(id: number) {
    router.push({
        name: "admin-config-edit",
        params: {
            id,
        },
    });
}

const createDialogVisible = ref(false);

function createConfig() {

    createDialogVisible.value = true;

}

const qrDialogVisible = ref(false);
const selectedConfig = ref<any | null>(null);  
function showQrCode(config:any) {

    selectedConfig.value = config;
    qrDialogVisible.value = true;

}

function downloadConfig(config:any) {

    if (!config.config_content) {
        return;
    }


    const blob = new Blob(
        [config.config_content],
        {
            type: "text/plain;charset=utf-8"
        }
    );


    const url = URL.createObjectURL(blob);


    const link = document.createElement("a");

    link.href = url;
    link.download = `${config.name || 'config'}.conf`;

    document.body.appendChild(link);

    link.click();


    document.body.removeChild(link);

    URL.revokeObjectURL(url);
}

</script>

<style scoped>
.stats-grid {
    display:grid;
    grid-template-columns:repeat(3,1fr);
    gap:1rem;
    margin-bottom:1.5rem;
}

.stat-value {
    font-size:2rem;
    font-weight:700;
}

.donate-card {
    margin-bottom:1.5rem;
}

.donate-actions {
    display:flex;
    gap:1rem;
    margin-top:1rem;
    flex-wrap:wrap;
}

.configs-toolbar {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
}

.search-field {
    flex: 1;
}

.my-configs-toggle {
    display: flex;
    align-items: center;
    gap: .5rem;
    white-space: nowrap;
}

.configs-grid {
    display:grid;   
    grid-template-columns:repeat(auto-fill,minmax(320px,1fr));
    gap:1rem;
}

.config-info {
    display:flex;
    flex-direction:column;
    gap:1rem;
}


.info-row {
  display:flex;
  align-items:center;
  gap:.75rem;
}
.info-row i {
  color: var(--p-primary-color);
  width:18px;
}
.ip-grid {
  display:grid;
  grid-template-columns:1fr 1fr;
  gap:.75rem;
}
.ip-card {
  background: var(--p-content-background);
  border: 1px solid var(--p-content-border-color);
  border-radius: 10px;
  padding: .75rem;
}

.ip-title {
  font-size: .75rem;
  color: var(--p-text-muted-color);
  margin-bottom: .35rem;
}

.ip-value {
  font-family: monospace;
  font-size: .85rem;
  color: var(--p-text-color);
  font-weight: 600;
  word-break: break-all;
}

.config-actions {
    display:flex;
    justify-content:center;
    gap:.75rem;
}

@media (max-width:768px){
.stats-grid{grid-template-columns:1fr;}
.configs-toolbar{flex-direction:column;align-items:stretch;}
}
</style>