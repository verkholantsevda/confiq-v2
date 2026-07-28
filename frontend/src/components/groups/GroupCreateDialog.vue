<template>
    <Dialog
        v-model:visible="visibleModel"
        modal
        :header="t('dialog.groups.name_dialog')"
        :style="{ width: '650px' }"
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
                    class="endpoint-list"
                >

                    <div
                        v-for="endpoint in endpoints"
                        :key="endpoint.id"
                        class="endpoint-item"
                    >
                        <Checkbox
                            v-model="form.endpoint_ids"
                            :inputId="'ep-' + endpoint.id"
                            :value="endpoint.id"
                        />

                        <label :for="'ep-' + endpoint.id">
                            {{ endpoint.name }}
                        </label>

                    </div>

                </div>

            </div>

            <Divider />

            <div class="actions">

                <Button
                    label="Создать группу"
                    icon="pi pi-check"
                    :loading="saving"
                    @click="create"
                />

                <Button
                    label="Отмена"
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
import { getEndpoints } from "@/api/endpoints";
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
        endpoints.value = await getEndpoints();
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

.endpoint-list {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: .75rem 1rem;
    margin-top: .75rem;
}

.endpoint-item {
    display: flex;
    align-items: center;
    gap: .6rem;
    padding: .35rem 0;
}

.endpoint-item label {
    font-weight: 400;
    cursor: pointer;
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