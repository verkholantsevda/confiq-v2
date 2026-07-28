import { ref } from "vue";

import * as groupsApi from "@/api/groups";
import type { Group, CreateGroupRequest, UpdateGroupRequest } from "@/types/group";

export function useGroups() {
    const groups = ref<Group[]>([]);
    const loading = ref(false);

    async function loadGroups() {
        loading.value = true;

        try {
            groups.value = await groupsApi.getGroups();
        } finally {
            loading.value = false;
        }
    }

    async function createGroup(group: CreateGroupRequest) {
        const created = await groupsApi.createGroup(group);
        groups.value.push(created);
    }

    async function updateGroup(id: number, group: UpdateGroupRequest) {
        const updated = await groupsApi.updateGroup(id, group);

        const index = groups.value.findIndex(g => g.id === id);

        if (index !== -1) {
            groups.value[index] = updated;
        }
    }

    async function deleteGroup(id: number) {
        await groupsApi.deleteGroup(id);

        groups.value = groups.value.filter(g => g.id !== id);
    }

    return {
        groups,
        loading,

        loadGroups,
        createGroup,
        updateGroup,
        deleteGroup,
    };
}