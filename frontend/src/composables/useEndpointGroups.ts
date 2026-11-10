import { ref } from "vue";

import { getEndpointsAll } from "@/api/endpoints";
import { getGroups } from "@/api/groups";
import { getGroupsEndpoints } from "@/api/groupsEndpoints";

import type { Endpoint } from "@/types/endpoint";
import type { Group } from "@/types/group";
import type { GroupEndpoint } from "@/types/groupEndpoint";

export interface EndpointWithGroups extends Endpoint {
    groups: Group[];
    configs?: any[];
}

export function useEndpointGroups() {

    const endpoints = ref<EndpointWithGroups[]>([]);
    const loading = ref(false);

    async function load(configsData: any[] = []) {
        loading.value = true;

        try {

            const [
                endpointsData,
                groupsData,
                relationsData,
            ] = await Promise.all([
                getEndpointsAll(),
                getGroups(),
                getGroupsEndpoints(),
            ]);

            endpoints.value = endpointsData.map(endpoint => {

                const groups = relationsData
                    .filter(r => r.endpoint_id === endpoint.id)
                    .map(r => groupsData.find(g => g.id === r.group_id))
                    .filter(Boolean) as Group[];

                const endpointConfigs = configsData.filter((c: any) => c?.endpoint?.id === endpoint.id);

                return {
                    ...endpoint,
                    groups,
                    configs: endpointConfigs,
                };
            });

        } finally {
            loading.value = false;
        }
    }

    return {
        endpoints,
        loading,
        load,
    };
}