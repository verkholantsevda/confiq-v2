import api from "./client";

export interface GroupEndpoint {
    group_id: number;
    endpoint_id: number;
}

export async function getGroupsEndpoints() {
    const { data } = await api.get<GroupEndpoint[]>("/groups_endpoints");
    return data;
}