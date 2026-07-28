<template>
    <Dialog
        v-model:visible="visible"
        modal
        :header="t('groups.delete')"
        :style="{ width: '30rem' }"
    >
        <div class="content">
            <i class="pi pi-exclamation-triangle warning-icon" />

            <p>
                {{ t("groups.deleteConfirm") }}
            </p>

            <strong>{{ group?.name }}</strong>
        </div>

        <template #footer>
            <Button
                :label="t('common.cancel')"
                severity="secondary"
                outlined
                @click="close"
            />

            <Button
                :label="t('common.delete')"
                severity="danger"
                :loading="loading"
                @click="remove"
            />
        </template>
    </Dialog>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useI18n } from "vue-i18n";

import Dialog from "primevue/dialog";
import Button from "primevue/button";

import { deleteGroup } from "@/api/groups";
import type { Group } from "@/types/group";

const { t } = useI18n();

const emit = defineEmits<{
    (e: "deleted"): void;
}>();

const visible = ref(false);
const loading = ref(false);

const group = ref<Group | null>(null);

function open(item: Group) {
    group.value = item;
    visible.value = true;
}

function close() {
    visible.value = false;
    group.value = null;
}

async function remove() {
    if (!group.value) return;

    loading.value = true;

    try {
        await deleteGroup(group.value.id);

        emit("deleted");

        close();
    } finally {
        loading.value = false;
    }
}

defineExpose({
    open,
});
</script>

<style scoped>
.content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    text-align: center;
    padding: .5rem 0;
}

.warning-icon {
    font-size: 2.5rem;
    color: var(--p-red-500);
}
</style>