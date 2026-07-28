import api from "./client";
import type { User } from "@/types/user";

export async function getMe() {
    const { data } = await api.get<User>("/me");
    return data;
}

export async function changePassword(
    currentPassword: string,
    newPassword: string,
) {
    await api.put("/me/password", {
        current_password: currentPassword,
        new_password: newPassword,
    });
}