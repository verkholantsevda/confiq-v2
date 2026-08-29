import client from "./client";

export interface LoginRequest {
    username: string;
    password: string;
    totp_code?: string
}

export interface LoginResponse {
    token: string;
}

export async function login(data: LoginRequest) {
    const response = await client.post<LoginResponse>("/login", data);
    return response.data;
}

export interface MeResponse {
    id: number;
    username: string;
    config_limit: number;
    group_id: number | null;
    is_admin: boolean;
}

export async function me() {
    const response = await client.get<MeResponse>("/me");
    return response.data;
}