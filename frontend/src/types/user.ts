import type { Group } from "./group";

export interface CreateUserRequest {
    username: string;
    password: string;
    config_limit: number;
    group_id: number | null;
    is_admin: boolean;
}

export interface User {
    id: number;
    username: string;
    config_limit: number;
    group_id: number | null;
    is_admin: boolean;
    created_at: string;

    group?: Group | null;
    configs_count?: number;
}

export interface UpdateUserRequest {
    username: string;
    password: string;
    config_limit: number;
    group_id: number | null;
    is_admin: boolean;
    
}