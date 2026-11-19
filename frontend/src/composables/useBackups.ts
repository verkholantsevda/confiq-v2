import { ref } from "vue";

import {
    getBackups,
    createBackup,
    deleteBackup,
    type BackupFile,
} from "@/api/backups";

export function useBackups() {
    const backups = ref<BackupFile[]>([]);
    const loading = ref(false);

    async function load() {
        loading.value = true;

        try {
            backups.value = await getBackups();
        } finally {
            loading.value = false;
        }
    }

    async function create() {
        loading.value = true;

        try {
            await createBackup();
            await load();
        } finally {
            loading.value = false;
        }
    }

    async function remove(name: string) {
        loading.value = true;

        try {
            await deleteBackup(name);
            await load();
        } finally {
            loading.value = false;
        }
    }

    return {
        backups,
        loading,
        load,
        create,
        remove,
    };
}
