import api from "./client";

import type { ConfigType } from "@/types/configType";

export async function getConfigTypes() {
    const { data } = await api.get<ConfigType[]>("/config-types");

    return data;
}

export async function getConfigType(id: number) {
    const { data } = await api.get<ConfigType>(`/config-types/${id}`);

    return data;
}



export async function createConfigType(payload: Partial<ConfigType>) {
    const { data } = await api.post<ConfigType>(
        "/config-types",
        payload,
    );

    return data;
}

export async function updateConfigType(
    id: number,
    payload: Partial<ConfigType>,
) {
    const { data } = await api.put<ConfigType>(
        `/config-types/${id}`,
        payload,
    );

    return data;
}

export async function deleteConfigType(id: number) {
    await api.delete(`/config-types/${id}`);
}