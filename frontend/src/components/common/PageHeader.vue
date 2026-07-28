<template>
    <div class="page-header">

        <div class="page-info">

            <h1>{{ title }}</h1>

            <p v-if="subtitle">
                {{ subtitle }}
            </p>

        </div>

        <div class="page-actions">

            <slot name="actions" />

            <Button
                v-if="buttonLabel"
                :label="buttonLabel"
                icon="pi pi-plus"
                @click="$emit('create')"
            />

        </div>

    </div>
</template>

<script setup lang="ts">
import Button from "primevue/button";

defineProps<{
    title: string;
    subtitle?: string;
    buttonLabel?: string;
}>();

defineEmits<{
    (e: "create"): void;
}>();
</script>

<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    margin-bottom: 2rem;
    gap: 1rem;
}

.page-info h1 {
    margin: 0;
    font-size: 2rem;
    font-weight: 700;
}

.page-info p {
    margin-top: .35rem;
    color: var(--p-text-muted-color);
}

.page-actions {
    display: flex;
    align-items: center;
    gap: .75rem;
}

@media (max-width: 768px) {

    .page-header {
        flex-direction: column;
        align-items: stretch;
    }

    .page-actions {
        width: 100%;
        justify-content: stretch;
    }

    .page-actions :deep(.p-button) {
        width: 100%;
    }

    .page-info h1 {
        font-size: 1.6rem;
    }
}
</style>