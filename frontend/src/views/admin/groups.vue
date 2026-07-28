<template>
    <PageHeader
        :title="t('menu.groups')"
        :subtitle="t('pages.groups.description')"
        :button-label="t('pages.groups.create')"
        @create="createGroup"
    />

    <EntityTable
        :value="groups"
        :loading="loading"
        :empty-message="t('common.noData')"
        :search-placeholder="t('common.search')"
    >
        <Column field="id" header="ID" sortable />
        <Column field="name" :header="t('pages.groups.name')" sortable />
        <Column field="description" :header="t('pages.groups.groups_descriptions')" sortable />
        <Column field="users_count" :header="t('pages.groups.users')" sortable />
        <Column field="endpoints_count" :header="t('pages.groups.endpoints')" sortable />
        <Column field="created_at" :header="t('common.created')" sortable />
        <Column :header="t('common.actions')" style="width:120px">
            <template #body="{ data }">
                <Button icon="pi pi-pencil" text rounded @click="editGroup(data)" />
                <Button icon="pi pi-trash" severity="danger" text rounded @click="deleteGroup(data)" />
            </template>
        </Column>
    </EntityTable>
    <GroupCreateDialog
        v-model:visible="showCreateDialog"
        @created="load"
    />
    <GroupEditDialog
        v-model:visible="showEditDialog"
        :group="selectedGroup"
        @updated="load"
    />

    <GroupDeleteDialog
        v-model:visible="showDeleteDialog"
        :group="selectedGroup"
        @deleted="load"
    />
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import Button from 'primevue/button';
import Column from 'primevue/column';
import PageHeader from '@/components/common/PageHeader.vue';
import EntityTable from '@/components/common/EntityTable.vue';
import { useGroupsView } from '@/composables/useGroupsView';
import { ref } from "vue";
import GroupCreateDialog from "@/components/groups/GroupCreateDialog.vue";
import type { Group } from "@/types/group";
import GroupEditDialog from "@/components/groups/GroupEditDialog.vue";
import GroupDeleteDialog from "@/components/groups/GroupDeleteDialog.vue";


const { t } = useI18n({
    useScope: 'global',
});

const { groups, loading, load } = useGroupsView();

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const selectedGroup = ref<Group | null>(null);

onMounted(load);

function createGroup() {

    showCreateDialog.value = true;

}

function editGroup(group: Group) {
    console.log(group);
    selectedGroup.value = { ...group };
    showEditDialog.value = true;
}

function deleteGroup(group: Group) {
    selectedGroup.value = group;
    showDeleteDialog.value = true;
}
</script>