<template>
    <PageHeader
        :title="t('menu.configtypes')"
        :subtitle="t('pages.configtypes.description')"
        :button-label="t('pages.configtypes.createbutton')"
        @create="createConfigType"
    />

    <EntityTable
        :value="configTypes"
        :loading="loading"
        :empty-message="t('common.noData')"
        :search-placeholder="t('common.search')"
    >
        <Column field="id" header="ID" sortable />
        <Column field="name" :header="t('pages.configtypes.name')" sortable />
        <Column field="description" :header="t('pages.configtypes.description_type')" sortable />
        <Column field="is_active" :header="t('pages.configtypes.active')" sortable>
            <template #body="{ data }">
                <ToggleSwitch
                    :model-value="data.is_active"
                    @update:model-value="value => onToggleStatus(data, value)"
                />
            </template>
        </Column>
        <Column :header="t('common.actions')" style="width:120px">
            <template #body="{ data }">
                <Button icon="pi pi-pencil" text rounded @click="editConfigType(data.id)" />
                <Button icon="pi pi-trash" severity="danger" text rounded />
            </template>
        </Column>
    </EntityTable>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";

import { onMounted } from "vue";
import Button from "primevue/button";
import Column from "primevue/column";
import EntityTable from "@/components/common/EntityTable.vue";
import { useConfigTypes } from "@/composables/useConfigTypes";

import PageHeader from "@/components/common/PageHeader.vue";

import ToggleSwitch from "primevue/toggleswitch";
import { useRouter } from "vue-router";
import { updateConfigType } from "@/api/configTypes";
const { t } = useI18n({
    useScope: "global",
});

const router = useRouter();

const {
    configTypes,
    loading,
    loadConfigTypes,
} = useConfigTypes();

onMounted(loadConfigTypes);

async function onToggleStatus(configType: any, value: boolean) {
    const previous = configType.is_active;

    configType.is_active = value;

    try {
        await updateConfigType(configType.id, {
            ...configType,
            is_active: value,
        });
    } catch (e) {
        configType.is_active = previous;
        console.error(e);
    }
}

function editConfigType(id: number) {
    router.push(`/admin/config-types/${id}`);
}

function createConfigType() {
    router.push("/admin/config-types/create");
}
</script>