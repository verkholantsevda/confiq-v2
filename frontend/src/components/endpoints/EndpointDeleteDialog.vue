<template>
    <Dialog
        :visible="visible"
        modal
        :style="{ width: '30rem' }"
        header="Delete Endpoint"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="delete-content">

            <i class="pi pi-exclamation-triangle warning-icon" />

            <div class="text">

                <p>
                    Are you sure you want to delete endpoint
                    <strong>{{ endpoint?.name }}</strong>?
                </p>

                <small class="hint">
                    This action cannot be undone.
                </small>

            </div>

        </div>

        <template #footer>

            <Button
                label="Cancel"
                severity="secondary"
                outlined
                @click="emit('update:visible', false)"
            />

            <Button
                label="Delete"
                icon="pi pi-trash"
                severity="danger"
                :loading="loading"
                @click="remove"
            />

        </template>
    </Dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";

import Dialog from "primevue/dialog";
import Button from "primevue/button";

import type { Endpoint } from "@/types/endpoint";
import { deleteEndpoint } from "@/api/endpoints";

const props = defineProps<{
    visible: boolean;
    endpoint: Endpoint | null;
}>();

const emit = defineEmits<{
    (e: "update:visible", value: boolean): void;
    (e: "deleted"): void;
}>();

const loading = ref(false);

async function remove() {
    if (!props.endpoint) return;

    loading.value = true;

    try {
        await deleteEndpoint(props.endpoint.id);

        emit("deleted");
        emit("update:visible", false);
    } finally {
        loading.value = false;
    }
}
</script>

<style scoped>
.delete-content {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: .5rem 0;
}

.warning-icon {
    font-size: 2rem;
    color: var(--orange-500);
}

.text {
    display: flex;
    flex-direction: column;
    gap: .35rem;
}

.hint {
    color: var(--text-color-secondary);
}
</style>