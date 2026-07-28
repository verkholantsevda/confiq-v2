import api from "./client";
import type { User, CreateUserRequest, UpdateUserRequest } from "@/types/user";

export async function getUsers() {
    const { data } = await api.get<User[]>("/users");

    return data;
}

export async function getUser(id: number) {
    const { data } = await api.get<User>(`/users/${id}`);

    return data;
}

export async function createUser(payload: CreateUserRequest) {
    const { data } = await api.post<User>("/users", payload);
    return data;
}

export async function updateUser(id: number, payload: UpdateUserRequest) {
    const { data } = await api.put<User>(`/users/${id}`, payload);

    return data;
}

export async function deleteUser(id: number) {
    await api.delete(`/users/${id}`);
}