// src/composables/useConfigsView.ts

import { computed } from "vue";

import { useConfigurations } from "@/composables/useConfigs";

export function useConfigsView() {
    const configurationsStore = useConfigurations();
    
    async function load() {
        await configurationsStore.loadConfigurations();
    }

    const configs = computed(() => configurationsStore.configurations.value);

    const totalConfigs = computed(() => configs.value.length);

    const usedConfigs = computed(() => configs.value.length);

    // TODO: replace with the logged-in user's real config limit
    const configLimit = computed(() => 5);

    const availableConfigs = computed(
        () => Math.max(0, configLimit.value - usedConfigs.value),
    );

    const loading = computed(() => configurationsStore.loading.value);

    return {
        configs,
        loading,
        totalConfigs,
        usedConfigs,
        configLimit,
        availableConfigs,
        load,
    };
}