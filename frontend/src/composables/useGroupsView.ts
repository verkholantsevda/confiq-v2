// src/composables/useGroupsView.ts

import { ref } from "vue";

import { getGroups } from "@/api/groups";
import { getUsers } from "@/api/users";
import { getGroupsEndpoints } from "@/api/groupsEndpoints";

import type { Group } from "@/types/group";
import type { User } from "@/types/user";
import type { GroupEndpoint } from "@/types/groupEndpoint";

export interface GroupView extends Group {
    users_count: number;
    endpoints_count: number;
}

export function useGroupsView() {
    const groups = ref<GroupView[]>([]);
    const loading = ref(false);

    async function load() {
        loading.value = true;

        try {
            const [groupsData, usersData, groupsEndpointsData] = await Promise.all([
                getGroups(),
                getUsers(),
                getGroupsEndpoints(),
            ]);

            groups.value = groupsData.map((group: Group) => ({
                ...group,

                users_count: usersData.filter(
                    (user: User) => user.group_id === group.id,
                ).length,

                endpoints_count: groupsEndpointsData.filter(
                    (relation: GroupEndpoint) =>
                        relation.group_id === group.id,
                ).length,
            }));
        } finally {
            loading.value = false;
        }
    }

    return {
        groups,
        loading,
        load,
    };
}