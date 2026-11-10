<template>
    <PageHeader
        :title="t('menu.endpoints')"
        :subtitle="t('pages.endpoints.description')"
        :button-label="t('pages.endpoints.create')"
        @create="createEndpoint"
    />

    <EntityTable
        :value="endpoints"
        :loading="loading"
        :empty-message="t('common.noData')"
        :search-placeholder="t('common.search')"
    >
        <Column field="id" header="ID" sortable />
        <Column field="name" :header="t('pages.endpoints.name')" sortable />
        <Column field="address" :header="t('pages.endpoints.address')" sortable />
        <Column field="port" :header="t('pages.endpoints.port')" sortable />
        <Column :header="t('pages.endpoints.configs')" style="width:100px">
            <template #body="{ data }">
                <button
                  v-if="data.configs && data.configs.length"
                  type="button"
                  class="configs-count"
                  :title="t('dialog.endpoints.configs_title', { name: data.name })"
                  @click="openConfigs(data)"
                >
                  {{ data.configs.length }}
                </button>
                <span v-else class="text-color-secondary">—</span>
            </template>
        </Column>
        <Column :header="t('pages.endpoints.groups')">
            <template #body="{ data }">
                <span v-if="data.groups.length">
                    {{ data.groups.map((g: any) => g.name).join(", ") }}
                </span>
                <span v-else class="text-color-secondary">—</span>
            </template>
        </Column>
        <Column :header="t('common.actions')" style="width:120px">
            <template #body="{ data }">
                <Button icon="pi pi-pencil" text rounded @click="editEndpoint(data)" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="deleteEndpoint(data)" />
            </template>
        </Column>
    </EntityTable>
    <EndpointCreateDialog
        v-model:visible="showCreateDialog"
        @created="load"
    />
    <EndpointEditDialog
        v-model:visible="showEditDialog"
        :endpoint="selectedEndpoint"
        @updated="load"
    />

    <EndpointDeleteDialog
        v-model:visible="showDeleteDialog"
        :endpoint="selectedEndpoint"
        @deleted="load"
    />
    <EndpointConfigsDialog
        v-model:visible="showConfigsDialog"
        :endpoint="configsEndpoint"
    />
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import Button from 'primevue/button';
import Column from 'primevue/column';

import PageHeader from "@/components/common/PageHeader.vue";
import EntityTable from '@/components/common/EntityTable.vue';
import { onMounted, ref } from "vue";
import EndpointCreateDialog from "@/components/endpoints/EndpointCreateDialog.vue";
import EndpointEditDialog from "@/components/endpoints/EndpointEditDialog.vue";
import EndpointDeleteDialog from "@/components/endpoints/EndpointDeleteDialog.vue";
import EndpointConfigsDialog from "@/components/endpoints/EndpointConfigsDialog.vue";
import type { Endpoint } from "@/types/endpoint";
import { useEndpointGroups } from "@/composables/useEndpointGroups";
import { getAllConfigurations } from "@/api/configs";
    
const { t } = useI18n({
    useScope: "global",
});

const { endpoints, loading, load } = useEndpointGroups();
const configs = ref<any[]>([]);
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const showConfigsDialog = ref(false);
const selectedEndpoint = ref<Endpoint | null>(null);
const configsEndpoint = ref<Endpoint | null>(null);

onMounted(async () => {
  try {
    configs.value = await getAllConfigurations();
    await load(configs.value);
  } catch {
    configs.value = [];
    await load([]);
  }
});

function createEndpoint() {
    showCreateDialog.value = true;
}

function editEndpoint(endpoint: Endpoint) {
    selectedEndpoint.value = { ...endpoint };
    showEditDialog.value = true;
}

function deleteEndpoint(endpoint: Endpoint) {
    selectedEndpoint.value = endpoint;
    showDeleteDialog.value = true;
}

function openConfigs(endpoint: Endpoint) {
    configsEndpoint.value = { ...endpoint };
    showConfigsDialog.value = true;
}
</script>