import { computed, ref } from "vue";

import { useConfigurations } from "@/composables/useConfigs";
import { getUsers } from "@/api/users";

export function useConfigsView() {
    const configurationsStore = useConfigurations();
    const showOnlyMine = ref(false);

    async function load() {
        if (showOnlyMine.value) {
            await configurationsStore.loadConfigurations();
        } else {
            await configurationsStore.loadAllConfigurations();
        }
    }

    const configs = computed(() => configurationsStore.configurations.value);

    const totalConfigs = computed(() => configs.value.length);

    const usedConfigs = computed(() => configs.value.length);

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
        showOnlyMine,
        load,
    };
}