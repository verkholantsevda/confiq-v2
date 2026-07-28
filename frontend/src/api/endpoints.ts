import api from "./client";
import type {

    Endpoint,

    CreateEndpointRequest,

    UpdateEndpointRequest,

} from "@/types/endpoint";

export async function getEndpoints() {
    const { data } = await api.get<Endpoint[]>("/endpoints");
    return data;
}

export async function getEndpointConfigTypes(endpointId: number) {
    const response = await api.get(`/endpoints/${endpointId}/config-types`);
    return response.data;
}

export async function getEndpointsAll() {
    const { data } = await api.get<Endpoint[]>("/endpoints-all");
    return data;
}

export async function getEndpoint(id: number) {
    const { data } = await api.get<Endpoint>(`/endpoints/${id}`);
    return data;
}

export async function createEndpoint(payload: CreateEndpointRequest) {

    const { data } = await api.post<Endpoint>("/endpoints-all", payload);

    return data;

}

export async function updateEndpoint(
    id: number,
    payload: UpdateEndpointRequest
) {
    const { data } = await api.put<Endpoint>(
        `/endpoints-all/${id}`,
        payload
    );
    return data;
}

export async function deleteEndpoint(id: number) {
    await api.delete(`/endpoints-all/${id}`);
}