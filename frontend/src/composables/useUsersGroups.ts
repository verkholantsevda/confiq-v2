import { ref } from "vue";

import { getUsers } from "@/api/users";
import { getGroups } from "@/api/groups";
import { getConfigurations } from "@/api/configs";

import type { User } from "@/types/user";
import type { Group } from "@/types/group";

export interface UserWithGroup extends User {
    group: Group | null;
}

export function useUsersGroups() {
    const users = ref<UserWithGroup[]>([]);
    const loading = ref(false);

    async function load() {
        loading.value = true;

        try {
            const [usersData, groupsData, configsData] = await Promise.all([
                getUsers(),
                getGroups(),
                getConfigurations(),
            ]);

            const configsCount = new Map<number, number>();

            for (const config of configsData) {
                configsCount.set(
                    config.user_id,
                    (configsCount.get(config.user_id) ?? 0) + 1,
                );
            }

            users.value = usersData.map(user => ({
                ...user,
                group: user.group_id
                    ? groupsData.find(g => g.id === user.group_id) ?? null
                    : null,
                configs_count: configsCount.get(user.id) ?? 0,
            }));
        } finally {
            loading.value = false;
        }
    }

    return {
        users,
        loading,
        load,
    };
}