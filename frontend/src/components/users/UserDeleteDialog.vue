<template>
    <Dialog
        :visible="visible"
        modal
        :style="{ width: '30rem' }"
        :header="t('pages.users.delete')"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="content">
            <i class="pi pi-exclamation-triangle warning-icon"></i>

            <div class="message">
                <p>
                    {{ t("users.message_delete") }}
                    <strong>{{ user?.username }}</strong>?
                </p>

                <small>
                    {{ t("users.cant_message_delete") }}
                </small>
            </div>
        </div>

        <template #footer>
            <Button
                label="Отмена"
                severity="secondary"
                outlined
                @click="emit('update:visible', false)"
            />

            <Button
                label="Удалить"
                icon="pi pi-trash"
                severity="danger"
                :loading="deleting"
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

import { deleteUser } from "@/api/users";
import type { User } from "@/types/user";

const { t } = useI18n();

const props = defineProps<{
    visible: boolean;
    user: User | null;
}>();

const emit = defineEmits<{
    (e: "update:visible", value: boolean): void;
    (e: "deleted"): void;
}>();

const deleting = ref(false);

async function remove() {
    if (!props.user) return;

    deleting.value = true;

    try {
        await deleteUser(props.user.id);

        emit("deleted");
        emit("update:visible", false);
    } finally {
        deleting.value = false;
    }
}
</script>

<style scoped>
.content {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    padding-top: .5rem;
}

.warning-icon {
    font-size: 2rem;
    color: var(--p-red-500);
}

.message {
    display: flex;
    flex-direction: column;
    gap: .5rem;
}

.message p {
    margin: 0;
    font-size: 1rem;
}

.message small {
    color: var(--text-color-secondary);
}
</style>