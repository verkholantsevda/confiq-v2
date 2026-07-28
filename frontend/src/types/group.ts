import type { Endpoint } from "./endpoint";
export interface Group {
    id: number;

    name: string;
    description?: string;

    endpoint_ids: number[];
    endpoints?: Endpoint[];

    created_at?: string;
    updated_at?: string;
}

export interface CreateGroupRequest {
    name: string;
    description?: string;

    endpoint_ids: number[];
}

export interface UpdateGroupRequest {
    name: string;
    description?: string;

    endpoint_ids?: number[];
}