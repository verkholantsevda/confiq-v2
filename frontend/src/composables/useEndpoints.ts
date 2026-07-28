import { ref } from "vue";

import * as endpointsApi from "@/api/endpoints";
import type { Endpoint, CreateEndpointRequest, UpdateEndpointRequest } from "@/types/endpoint";

export function useEndpoints() {
    const endpoints = ref<Endpoint[]>([]);
    const loading = ref(false);

    async function loadEndpoints() {
        loading.value = true;

        try {
            endpoints.value = await endpointsApi.getEndpointsAll();
        } finally {
            loading.value = false;
        }
    }

    async function createEndpoint(endpoint: CreateEndpointRequest) {
        const created = await endpointsApi.createEndpoint(endpoint);
        endpoints.value.push(created);
    }

    async function updateEndpoint(id: number, endpoint: UpdateEndpointRequest) {
        const updated = await endpointsApi.updateEndpoint(id, endpoint);

        const index = endpoints.value.findIndex(e => e.id === id);

        if (index !== -1) {
            endpoints.value[index] = updated;
        }
    }

    async function deleteEndpoint(id: number) {
        await endpointsApi.deleteEndpoint(id);

        endpoints.value = endpoints.value.filter(e => e.id !== id);
    }

    return {
        endpoints,
        loading,

        loadEndpoints,
        createEndpoint,
        updateEndpoint,
        deleteEndpoint,
    };
}