import { ref } from "vue";

import * as configTypesApi from "@/api/configTypes";
import type { ConfigType } from "@/types/configType";

export function useConfigTypes() {
    const configTypes = ref<ConfigType[]>([]);
    const loading = ref(false);

    async function loadConfigTypes() {
        loading.value = true;

        try {
            configTypes.value = await configTypesApi.getConfigTypes();
        } finally {
            loading.value = false;
        }
    }

    async function createConfigType(configType: Partial<ConfigType>) {
        const created = await configTypesApi.createConfigType(configType);

        configTypes.value.push(created);
    }

    async function updateConfigType(
        id: number,
        configType: Partial<ConfigType>,
    ) {
        const updated = await configTypesApi.updateConfigType(
            id,
            configType,
        );

        const index = configTypes.value.findIndex(
            (c) => c.id === id,
        );

        if (index !== -1) {
            configTypes.value[index] = updated;
        }
    }

    async function deleteConfigType(id: number) {
        await configTypesApi.deleteConfigType(id);

        configTypes.value = configTypes.value.filter(
            (c) => c.id !== id,
        );
    }

    return {
        configTypes,
        loading,

        loadConfigTypes,
        createConfigType,
        updateConfigType,
        deleteConfigType,
    };
}