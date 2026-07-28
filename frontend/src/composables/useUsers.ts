import { ref } from "vue";

import { getUsers } from "@/api/users";
import type { User } from "@/types/user";

export function useUsers() {

    const users = ref<User[]>([]);
    const loading = ref(false);

    async function loadUsers() {

        loading.value = true;

        try {

            users.value = await getUsers();

        } finally {

            loading.value = false;

        }

    }

    return {
        users,
        loading,
        loadUsers,
    };

}