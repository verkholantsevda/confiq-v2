<template>
    <Dialog
        :visible="visible"
        modal
        :style="{ width: '42rem' }"
        :header="t('pages.groups.edit')"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="form">

            <div class="field">
                <label>{{ t("groups.name") }}</label>

                <InputText
                    v-model="form.name"
                    fluid
                />

                <small class="hint">
                    Уникальное имя для группы
                </small>
            </div>

            <div class="field">
                <label>{{ t("groups.description") }}</label>

                <Textarea
                    v-model="form.description"
                    rows="3"
                    autoResize
                    fluid
                />

                <small class="hint">
                    Необязательное описание группы
                </small>
            </div>

            <Divider />

            <div class="field">
                <label>Доступные Endpoints</label>

                <small class="hint">
                    Выберите endpoints доступные для этой группы
                </small>

                <div class="endpoints">

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

        </div>

        <template #footer>

            <Button
                label="Отмена"
                severity="secondary"
                outlined
                @click="emit('update:visible', false)"
            />

            <Button
                label="Сохранить"
                icon="pi pi-check"
                :loading="saving"
                @click="save"
            />

        </template>
    </Dialog>
</template>

<script setup lang="ts">
import { reactive, watch, ref, onMounted } from "vue";
import { useI18n } from "vue-i18n";

import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import Checkbox from "primevue/checkbox";
import Divider from "primevue/divider";
import Button from "primevue/button";

import type { Endpoint } from "@/types/endpoint";
import type { Group } from "@/types/group";

import { updateGroup } from "@/api/groups";
import { getEndpoints } from "@/api/endpoints";

const { t } = useI18n();

const props = defineProps<{
    visible: boolean;
    group: Group | null;
}>();

const emit = defineEmits<{
    (e: "update:visible", value: boolean): void;
    (e: "updated"): void;
}>();

const saving = ref(false);

const endpoints = ref<Endpoint[]>([]);

onMounted(async () => {
    endpoints.value = await getEndpoints();
});

const form = reactive({
    id: 0,
    name: "",
    description: "",
    endpoint_ids: [] as number[],
});

watch(
  [() => props.group, () => props.visible],
  ([group, visible]) => {
    if (!group || !visible) return;

    form.id = group.id;
    form.name = group.name;
    form.description = group.description ?? "";

    form.endpoint_ids.splice(0);
    if (group.endpoint_ids) {
      form.endpoint_ids.push(...group.endpoint_ids);
    }
  },
  { immediate: true }
);

async function save() {
    if (!form.id) return;

    saving.value = true;

    try {
        await updateGroup(form.id, {
            name: form.name,
            description: form.description,
            endpoint_ids: form.endpoint_ids,
        });

        form.endpoint_ids = [...form.endpoint_ids];
        emit("updated");
        emit("update:visible", false);
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
    color: var(--text-color-secondary);
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
</style>