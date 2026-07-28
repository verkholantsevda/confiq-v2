import api from "./client";
import type {
    Group,
    CreateGroupRequest,
    UpdateGroupRequest,
} from "@/types/group";

export async function getGroups() {
    const { data } = await api.get<Group[]>("/groups");
    return data;
}

export async function getGroup(id: number) {
    const { data } = await api.get<Group>(`/groups/${id}`);
    return data;
}

export async function createGroup(payload: CreateGroupRequest) {
    const { data } = await api.post<Group>("/groups", payload);
    return data;
}

export async function updateGroup(id: number, payload: UpdateGroupRequest) {
    const { data } = await api.put<Group>(`/groups/${id}`, payload);
    return data;
}

export async function deleteGroup(id: number) {
    await api.delete(`/groups/${id}`);
}