import { ref, onMounted } from "vue";

export interface DonateConfig {
    yoomoneyUrl: string;
    cloudtipUrl: string;
}

const config = ref<DonateConfig>({
    yoomoneyUrl: import.meta.env.VITE_DONATE_YOOMONEY_URL ?? "#",
    cloudtipUrl: import.meta.env.VITE_DONATE_CLOUDTIP_URL ?? "#",
});

const loaded = ref(false);
const error = ref<string | null>(null);

export function useDonate() {
    const load = async () => {
        if (loaded.value) return;
        try {
            const res = await fetch("/config.json", { cache: "no-store" });
            if (res.ok) {
                const data = (await res.json()) as Partial<DonateConfig>;
                config.value = {
                    yoomoneyUrl: data.yoomoneyUrl ?? config.value.yoomoneyUrl,
                    cloudtipUrl: data.cloudtipUrl ?? config.value.cloudtipUrl,
                };
            } else {
                error.value = `config.json returned ${res.status}`;
            }
        } catch (e) {
            error.value = e instanceof Error ? e.message : String(e);
        } finally {
            loaded.value = true;
        }
    };

    onMounted(() => {
        void load();
    });

    return {
        config,
        loaded,
        error,
        reload: load,
    };
}
