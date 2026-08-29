import api from "./client";

export interface TOTPStatus {
    enabled: boolean;
    available: boolean;
}
export async function getTOTPStatus(): Promise<TOTPStatus> {
    const response = await api.get<TOTPStatus>("/totp/status");
    return response.data;
}

export async function setupTOTP() {
    const response = await api.post("/totp/setup");
    return response.data;
}

export async function enableTOTP(code: string) {
    await api.post("/totp/enable", {
        code,
    });
}

export async function disableTOTP() {
    await api.delete("/totp");
}