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
import type { Endpoint } from "@/types/endpoint";
import { useEndpointGroups } from "@/composables/useEndpointGroups";
    
const { t } = useI18n({
    useScope: "global",
});

const { endpoints, loading, load } = useEndpointGroups();
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const selectedEndpoint = ref<Endpoint | null>(null);

onMounted(load);

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
</script>