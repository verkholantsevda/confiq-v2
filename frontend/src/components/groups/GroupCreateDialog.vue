<template>
    <Dialog
        v-model:visible="visibleModel"
        modal
        :style="{ width: '62rem' }"
        :breakpoints="{
            '1200px': '80vw',
            '768px': '95vw',
            '560px': '100vw'
        }"
        :maximizable="true"
        :header="t('dialog.groups.name_dialog')"
    >
        <div class="form">

            <div class="field">
                <label>{{ t("dialog.groups.name") }}</label>

                <InputText
                    v-model="form.name"
                    fluid
                    :placeholder="t('dialog.groups.name_placeholder')"
                />

                <small class="hint">
                    <label>{{ t("dialog.groups.name_hint") }}</label>
                </small>
            </div>

            <div class="field">
                <label>{{ t("dialog.groups.description") }}</label>

                <Textarea
                    v-model="form.description"
                    rows="3"
                    :placeholder="t('dialog.groups.description_placeholder')"
                    autoResize
                    fluid
                />

                <small class="hint">
                    <label>{{ t("dialog.groups.description_hint") }}</label>
                </small>
            </div>

            <Divider />

            <div class="field">

                <label>{{ t("dialog.groups.endpoints") }}</label>

                <small class="hint">
                    <label>{{ t("dialog.groups.endpoints_hint") }}</label>
                </small>

                <div
                    v-if="loadingEndpoints"
                    class="loading"
                >
                    Загрузка...
                </div>

                <div
                    v-else
                    class="endpoints"
                >

                    <div
                        v-for="endpoint in endpoints"
                        :key="endpoint.id"
                        class="endpoint-row"
                    >
                        <Checkbox
                            v-model="form.endpoint_ids"
                            :input-id="`endpoint-${endpoint.id}`"
                            :value="endpoint.id"
                        />

                        <label :for="`endpoint-${endpoint.id}`">
                            <strong>{{ endpoint.name }}</strong>
                            <span class="address">
                                {{ endpoint.address }}:{{ endpoint.port }}
                            </span>
                        </label>
                    </div>

                </div>

            </div>

            <Divider />

            <div class="actions">

                <Button
                    :label="t('dialog.groups.create')"
                    icon="pi pi-check"
                    :loading="saving"
                    @click="create"
                />

                <Button
                    :label="t('common.cancel')"
                    icon="pi pi-times"
                    severity="secondary"
                    outlined
                    @click="visibleModel = false"
                />

            </div>

        </div>

    </Dialog>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { computed, onMounted, reactive, ref } from "vue";

import Dialog from "primevue/dialog";
import Button from "primevue/button";
import Divider from "primevue/divider";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import Checkbox from "primevue/checkbox";

// заменить своими API
import { createGroup } from "@/api/groups";
import { getEndpointsAll } from "@/api/endpoints";
const { t } = useI18n({
    useScope: "global",
});
const props = defineProps<{
    visible: boolean;
}>();

const emit = defineEmits([
    "update:visible",
    "created",
]);

const visibleModel = computed({
    get: () => props.visible,
    set: value => emit("update:visible", value),
});

const saving = ref(false);
const loadingEndpoints = ref(false);

const endpoints = ref<any[]>([]);

const form = reactive({
    name: "",
    description: "",
    endpoint_ids: [] as number[],
});

async function loadEndpoints() {
    loadingEndpoints.value = true;

    try {
        endpoints.value = await getEndpointsAll();
    } finally {
        loadingEndpoints.value = false;
    }
}

onMounted(loadEndpoints);

async function create() {

    saving.value = true;

    try {

        await createGroup({
            name: form.name,
            description: form.description,
            endpoint_ids: form.endpoint_ids,
        });

        emit("created");

        visibleModel.value = false;

        form.name = "";
        form.description = "";
        form.endpoint_ids = [];

    } finally {

        saving.value = false;

    }

}
</script>

<style scoped>
.form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
}

.field {
    display: flex;
    flex-direction: column;
    gap: .5rem;
}

.field label {
    font-weight: 600;
}

.hint {
    color: var(--p-text-muted-color);
    font-size: .85rem;
}

.endpoints {
    display: flex;
    flex-direction: column;
    gap: .75rem;
    max-height: 260px;
    overflow-y: auto;
    padding: .25rem;
}

.endpoint-row {
    display: flex;
    align-items: flex-start;
    gap: .75rem;
}

.endpoint-row label {
    display: flex;
    flex-direction: column;
    cursor: pointer;
}

.address {
    font-size: .85rem;
    color: var(--text-color-secondary);
}

.loading {
    padding: 1rem 0;
    color: var(--p-text-muted-color);
}

.actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
}

:deep(.p-dialog-header-actions .p-dialog-close-button) {
    border-radius: 0;
}

:deep(.p-dialog-header-actions .p-dialog-close-button .p-icon) {
    width: 1rem;
    height: 1rem;
}

:deep(.p-dialog-footer) {
    display: flex;
    gap: .75rem;
}

@media (max-width: 600px) {
    :deep(.p-dialog-footer) {
        flex-direction: column-reverse;
    }

    :deep(.p-dialog-footer .p-button) {
        width: 100%;
    }

    :deep(.p-dialog-content) {
        padding: 1rem;
    }

    :deep(.p-dialog-header) {
        padding: 1rem;
    }

    :deep(.p-dialog-footer) {
        padding: 1rem;
    }
}

@media (max-width: 640px) {

    .endpoint-list {
        grid-template-columns: 1fr;
    }

    .actions {
        flex-direction: column;
    }

    .actions :deep(.p-button) {
        width: 100%;
    }

}
</style>