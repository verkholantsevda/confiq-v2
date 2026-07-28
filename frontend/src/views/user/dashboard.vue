<template>
    <PageHeader
        :title="t('menu.dashboard')"
        :subtitle="t('pages.dashboard.description')"
    />


    <!-- Статистика -->
    <div class="stats-grid">

        <Card>
            <template #title>
                {{ t("pages.configs.total") }}
            </template>

            <template #content>
                <div class="stat-value">
                    {{ totalConfigs }}
                </div>
            </template>
        </Card>


        <Card>
            <template #title>
                {{ t("pages.configs.used") }}
            </template>

            <template #content>
                <div class="stat-value">
                    {{ usedConfigs }} / {{ configLimit }}
                </div>
            </template>
        </Card>


        <Card>
            <template #title>
                {{ t("pages.configs.available") }}
            </template>

            <template #content>
                <div class="stat-value">
                    {{ availableConfigs }}
                </div>
            </template>
        </Card>

    </div>



    <!-- Использование -->
    <Card class="usage-card">

        <template #title>
            Использование
        </template>

        <template #content>

            <ProgressBar
                :value="usagePercent"
            />

            <div class="usage-text">
                {{ usedConfigs }} из {{ configLimit }}
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
                />

                <Button
                    label="CloudTip"
                    icon="pi pi-credit-card"
                    severity="secondary"
                    as="a"
                    :href="DONATE_CLOUDTIP_URL"
                    target="_blank"
                />

            </div>

        </template>

    </Card>



    <!-- Быстрые действия -->
    <Card class="actions-card">

        <template #title>
            Быстрые действия
        </template>


        <template #content>

            <div class="actions-grid">

                <Button
                    icon="pi pi-plus"
                    :label="t('user.quickactions.create')"
                />


                <Button
                    icon="pi pi-list"
                    :label="t('user.quickactions.allconf')"
                    @click="router.push('/user/configs')"
                />


                <Button
                    icon="pi pi-user"
                    :label="t('user.quickactions.profile')"
                    @click="router.push('/user/profile')"
                />

            </div>

        </template>

    </Card>



    <!-- Конфигурации -->

    <div class="configs-grid">

        <Card
            v-for="config in configs"
            :key="config.id"
        >

            <template #title>
                {{ config.name }}
            </template>

            <template #subtitle>
                {{ config.config_type?.name || '-' }}
            </template>

            <template #content>
                <div class="info-row">
                    <i class="pi pi-map-marker" />
                    {{ config.endpoint?.name || '-' }}
                </div>

                <div class="info-row">
                    <i class="pi pi-globe" />
                    {{ config.endpoint?.address }}:{{ config.endpoint?.port }}
                </div>

                <div class="info-row">
                    <i class="pi pi-clock" />
                    {{ formatDate(config.created_at) }}
                </div>

                <Divider />

                <div class="config-actions">
                    <Button
                        icon="pi pi-eye"
                        rounded
                        @click="openConfig(config.id)"
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
            </template>

        </Card>

    </div>

    <ConfigQrCodeDialog
        v-model="qrDialogVisible"
        :config="selectedConfig"
    />

</template>



<script setup lang="ts">

import { computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";


import Card from "primevue/card";
import Button from "primevue/button";
import Divider from "primevue/divider";
import ProgressBar from "primevue/progressbar";


import PageHeader from "@/components/common/PageHeader.vue";


import { useConfigsView } from "@/composables/useConfigsView";

import { ref } from 'vue';
import ConfigQrCodeDialog from '@/components/configs/ConfigQrCodeDialog.vue';


const router = useRouter();


const { t } = useI18n({
    useScope:"global"
});



const DONATE_YOOMONEY_URL =
    import.meta.env.VITE_DONATE_YOOMONEY_URL ?? '#';


const DONATE_CLOUDTIP_URL =
    import.meta.env.VITE_DONATE_CLOUDTIP_URL ?? '#';



const {
    configs,
    totalConfigs,
    usedConfigs,
    availableConfigs,
    configLimit,
    load,
} = useConfigsView();



const usagePercent = computed(() => {

    if (!configLimit.value)
        return 0;


    return Math.round(
        usedConfigs.value / configLimit.value * 100
    );

});



onMounted(load);

const qrDialogVisible = ref(false);
const selectedConfig = ref<any | null>(null);

function openConfig(id:number){

    router.push(`/user/configs/${id}`);

}

function showQrCode(config: any) {
    selectedConfig.value = config;
    qrDialogVisible.value = true;
}

function downloadConfig(config: any) {
    if (!config.config_content) {
        return;
    }

    const blob = new Blob(
        [config.config_content],
        { type: 'text/plain;charset=utf-8' }
    );

    const url = URL.createObjectURL(blob);

    const link = document.createElement('a');
    link.href = url;
    link.download = `${(config.name || 'config').replace(/[^a-zA-Z0-9._-]/g, '_')}.conf`;

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    URL.revokeObjectURL(url);
}



function formatDate(value:string){

    return new Date(value).toLocaleString(
        "ru-RU",
        {
            day:"2-digit",
            month:"2-digit",
            year:"numeric",
            hour:"2-digit",
            minute:"2-digit"
        }
    );

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


.usage-card,
.donate-card,
.actions-card {
    margin-bottom:1.5rem;
}


.usage-text {
    margin-top:.5rem;
    text-align:right;
}


.donate-actions,
.actions-grid {
    display:flex;
    gap:1rem;
    flex-wrap:wrap;
}


.configs-grid {
    display:grid;
    grid-template-columns:repeat(auto-fill,minmax(320px,1fr));
    gap:1rem;
}


.config-header {
    display:flex;
    justify-content:space-between;
    align-items:center;
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


.config-actions {
    display:flex;
    justify-content:center;
    gap:.75rem;
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
  font-size:.75rem;
  color:var(--p-text-muted-color);
  margin-bottom:.35rem;
}

.ip-value {
  font-family:monospace;
  font-size:.85rem;
  font-weight:600;
  word-break:break-all;
}


@media(max-width:768px){

.stats-grid{
    grid-template-columns:1fr;
}

}

</style>