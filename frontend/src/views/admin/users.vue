<template>
    <PageHeader
        :title="t('menu.users')"
        :subtitle="t('pages.users.description')"
        :button-label="t('pages.users.create')"
        @create="createUser"
    />
    <div class="table-wrapper">
        <EntityTable
            :value="users"
            :loading="loading"
            :empty-message="t('common.noData')"
            :search-placeholder="t('common.search')"
        >
            <Column field="id" header="ID" sortable />
            <Column field="username" :header="t('pages.users.username')" sortable />
            <Column :header="t('pages.users.group')">
                <template #body="{ data }">
                    <span v-if="data.group">
                        {{ data.group.name }}
                    </span>
                    <span v-else class="text-color-secondary">—</span>
                </template>
            </Column>
            <Column field="role" :header="t('pages.users.role')">
                <template #body="{ data }">
                    <Tag
                        :severity="data.is_admin ? 'danger' : 'secondary'"
                        :value="data.is_admin ? t('pages.users.admin') : t('pages.users.user')"
                    />
                </template>
            </Column>
            <Column field="configs_count" :header="t('pages.users.configs')" sortable />
            <Column
                field="configurations"
                :header="t('pages.users.configs')"
                sortable
            />
            <Column field="created_at" :header="t('pages.users.created')" sortable>
                <template #body="{ data }">
                    {{ new Date(data.created_at).toLocaleString() }}
                </template>
            </Column>
            <Column :header="t('common.actions')" style="width:120px">
                <template #body="{ data }">
                    <Button icon="pi pi-pencil" text rounded @click="editUser(data)" />
                    <Button icon="pi pi-trash" severity="danger" text rounded @click="removeUser(data)" />
                </template>
            </Column>
        </EntityTable>
    </div>
    <UserCreateDialog
        v-model:visible="showCreateDialog"
        :groups="groups"
        @created="load"
    />

    <UserEditDialog
        v-model:visible="showEditDialog"
        :user="selectedUser"
        :groups="groups"
        @updated="load"
    />

    <UserDeleteDialog
        v-model:visible="showDeleteDialog"
        :user="selectedUser"
        @deleted="load"
    />
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Button from 'primevue/button';
import Column from 'primevue/column';
import Tag from 'primevue/tag';

import PageHeader from '@/components/common/PageHeader.vue';
import EntityTable from '@/components/common/EntityTable.vue';

import { useUsersGroups } from '@/composables/useUsersGroups';

import UserCreateDialog from '@/components/users/UserCreateDialog.vue';
import UserEditDialog from '@/components/users/UserEditDialog.vue';
import UserDeleteDialog from '@/components/users/UserDeleteDialog.vue';
import type { User } from '@/types/user';

import type { Group } from '@/types/group';
import { getGroups } from '@/api/groups';

const { t } = useI18n({ useScope: 'global' });

const { users, loading, load } = useUsersGroups();

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const showDeleteDialog = ref(false);
const selectedUser = ref<User | null>(null);

const groups = ref<Group[]>([]);

onMounted(async () => {
    await load();
    groups.value = await getGroups();
});

function createUser() {
    showCreateDialog.value = true;
}

function editUser(user: User) {
    selectedUser.value = { ...user };
    showEditDialog.value = true;
}

function removeUser(user: User) {
    selectedUser.value = user;
    showDeleteDialog.value = true;
}

</script>

<style scoped>
.table-wrapper {
    width: 100%;
}

@media (max-width: 768px) {
    .table-wrapper {
        overflow-x: auto;
        -webkit-overflow-scrolling: touch;
    }

    .table-wrapper :deep(table) {
        min-width: 900px;
    }

    .table-wrapper :deep(.p-datatable) {
        min-width: 900px;
    }
}
</style>