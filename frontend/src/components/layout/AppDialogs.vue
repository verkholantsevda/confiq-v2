

<template>
  <UserCreateDialog
    v-model:visible="dialogs.userCreate"
    :groups="groups"
    @created="loadGroups"
  />
  <GroupCreateDialog v-model:visible="dialogs.groupCreate" />
  <EndpointCreateDialog v-model:visible="dialogs.endpointCreate" />
  <ConfigCreateDialog v-model="dialogs.configCreate" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

import UserCreateDialog from '@/components/users/UserCreateDialog.vue';
import GroupCreateDialog from '@/components/groups/GroupCreateDialog.vue';
import EndpointCreateDialog from '@/components/endpoints/EndpointCreateDialog.vue';
import ConfigCreateDialog from '@/components/configs/ConfigCreateDialog.vue';

import { useDialogsStore } from '@/stores/dialog';
import { getGroups } from '@/api/groups';
import type { Group } from '@/types/group';

const dialogs = useDialogsStore();

const groups = ref<Group[]>([]);

async function loadGroups() {
  groups.value = await getGroups();
}

onMounted(loadGroups);
</script>