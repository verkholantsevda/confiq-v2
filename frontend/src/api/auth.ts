import client from "./client";
import type { User } from "@/types/user";

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

export interface MeResponse extends User {
    totp_available: boolean;
}

export async function me() {
    const response = await client.get<MeResponse>("/me");
    return response.data;
}