import axios from "axios";

const baseURL = import.meta.env.VITE_API_URL || "/api";

const client = axios.create({
    baseURL,
    headers: {
        "Content-Type": "application/json",
    },
});

client.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");

    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
});

client.interceptors.response.use((response) => {
    const refreshed = response.headers["x-refresh-token"];

    if (typeof refreshed === "string" && refreshed) {
        localStorage.setItem("token", refreshed);
    }

    return response;
});

export default client;