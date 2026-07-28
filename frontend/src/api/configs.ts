import api from "./client";
import type { Configuration } from "@/types/config";

export async function getConfigurations() {
    const { data } = await api.get<Configuration[]>("/configs");
    return data;
}

export async function getConfiguration(id: number) {
    const { data } = await api.get<Configuration>(`/configs/${id}`);
    return data;
}

export async function getAllConfigurations() {
    const { data } = await api.get<Configuration[]>("/all-configs");
    return data;
}

export async function createConfiguration(payload: Partial<Configuration>) {
    const { data } = await api.post<Configuration>("/configs", payload);
    return data;
}

export async function updateConfiguration(
    id: number,
    payload: Partial<Configuration>,
) {
    const { data } = await api.put<Configuration>(
        `/configs/${id}`,
        payload,
    );

    return data;
}

export async function deleteConfiguration(id: number) {
    try {
        await api.delete(`/configs/${id}`);
    } catch (error: any) {
        console.error(
            "Delete config error:",
            error.response?.data,
            error.response?.status,
        );

        throw error;
    }
}