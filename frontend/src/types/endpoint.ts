export interface Endpoint {
    id: number;

    name: string;
    address: string;
    port: number;

    created_at?: string;
}

export interface CreateEndpointRequest {
    name: string;
    address: string;
    port: number;

    group_ids: number[];
    config_type_ids: number[];
}

export interface UpdateEndpointRequest {
    name: string;
    address: string;
    port: number;

    group_ids: number[];
    config_type_ids: number[];
}