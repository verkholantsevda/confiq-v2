import api from "./client";

export interface BackupFile {
    name: string;
    size: number;
    created_at: string;
}

export async function getBackups() {
    const { data } = await api.get<BackupFile[]>("/backups");

    return data;
}

export async function createBackup() {
    const { data } = await api.post<BackupFile>("/backups");

    return data;
}

export async function deleteBackup(name: string) {
    await api.delete(`/backups/${encodeURIComponent(name)}`);
}

export async function restoreBackup(file: File) {
    const formData = new FormData();

    formData.append("file", file);

    const { data } = await api.post<{ restarting: boolean }>(
        "/backups/restore",
        formData,
        { headers: { "Content-Type": "multipart/form-data" } },
    );

    return data;
}

export async function restoreFromBackup(name: string) {
    const { data } = await api.post<{ restarting: boolean }>(
        `/backups/${encodeURIComponent(name)}/restore`,
    );

    return data;
}

export async function downloadBackup(name: string) {
    const response = await api.get<Blob>(
        `/backups/${encodeURIComponent(name)}/download`,
        { responseType: "blob" },
    );

    const url = URL.createObjectURL(response.data);
    const link = document.createElement("a");

    link.href = url;
    link.download = name;
    document.body.appendChild(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
}
