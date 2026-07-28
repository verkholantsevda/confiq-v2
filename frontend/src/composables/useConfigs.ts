import { ref } from "vue";

import * as configsApi from "@/api/configs";
import type { Configuration } from "@/types/config";

export function useConfigurations() {
    const configurations = ref<Configuration[]>([]);
    const loading = ref(false);

    async function loadConfigurations() {
        loading.value = true;

        try {
            configurations.value =
                await configsApi.getConfigurations();
        } finally {
            loading.value = false;
        }
    }

    async function loadAllConfigurations() {
        loading.value = true;

        try {
            configurations.value =
                await configsApi.getAllConfigurations();
        } finally {
            loading.value = false;
        }
    }

    async function createConfiguration(
        configuration: Partial<Configuration>,
    ) {
        const created =
            await configsApi.createConfiguration(configuration);

        configurations.value.push(created);
    }

    async function updateConfiguration(
        id: number,
        configuration: Partial<Configuration>,
    ) {
        const updated =
            await configsApi.updateConfiguration(id, configuration);

        const index = configurations.value.findIndex(
            c => c.id === id,
        );

        if (index !== -1) {
            configurations.value[index] = updated;
        }
    }

    async function deleteConfiguration(id: number) {
        await configsApi.deleteConfiguration(id);

        configurations.value =
            configurations.value.filter(c => c.id !== id);
    }

    async function loadConfiguration(id: number) {
        loading.value = true;

        try {
            return await configsApi.getConfiguration(id);
        } finally {
            loading.value = false;
        }
    }

    return {
        configurations,
        loading,
        loadConfigurations,
        loadAllConfigurations,
        loadConfiguration,
        createConfiguration,
        updateConfiguration,
        deleteConfiguration,
    };
}

