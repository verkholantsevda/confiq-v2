<template>
    <Dialog
        :visible="visible"
        modal
        :style="{ width: '62rem' }"
        :breakpoints="{
            '1200px': '80vw',
            '768px': '95vw',
            '560px': '100vw'
        }"
        :maximizable="true"
        :header="t('dialog.users.create')"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="layout">

            <!-- Левая часть -->
            <div class="form">

                <div class="field">
                    <label>{{ t("dialog.users.username.value") }}</label>

                    <InputText
                        v-model="form.username"
                        fluid
                    />

                    <small class="hint">
                        <label>{{ t("dialog.users.username.description") }}</label>
                    </small>
                </div>

                <div class="field">
                    <label>{{ t("dialog.users.password.value") }}</label>

                    <Password
                        v-model="form.password"
                        toggleMask
                        :feedback="false"
                        fluid
                    />

                    <small class="hint">
                        <label>{{ t("dialog.users.password.description") }}</label>
                    </small>
                </div>

                <div class="field">
                    <label>{{ t("dialog.users.limit_configurations") }}</label>

                    <InputNumber
                        v-model="form.config_limit"
                        showButtons
                        :min="1"
                        :max="100"
                        fluid
                    />

                    <small class="hint">
                        <label>{{ t("dialog.users.limit_configurations_hint") }}</label>
                    </small>
                </div>

                <div class="field">
                    <label>{{ t("dialog.users.group.value") }}</label>

                    <Select
                        v-model="form.group_id"
                        :options="groupOptions"
                        optionLabel="name"
                        optionValue="id"
                        fluid
                    />

                    <small class="hint">
                        <label>{{ t("dialog.users.group.description") }}</label>
                    </small>
                </div>

                <Divider />

                <div class="field administrator">
                    <Checkbox
                        v-model="form.is_admin"
                        binary
                        inputId="is-admin"
                    />

                    <label for="is-admin">
                        {{ t("dialog.users.admin.value") }}
                    </label>
                </div>

                <small class="hint">
                    <label>{{ t("dialog.users.admin.description") }}</label>
                </small>

            </div>

            <!-- Правая часть -->
            <div class="info">

                <Card>

                    <template #title>
                        <label>{{ t("dialog.users.info_title") }}</label>
                    </template>

                    <template #content>

                        <div class="info-section">

                            <label>{{ t("dialog.users.info_hint") }}</label>

                            <div class="info-row">
                                <strong> <label>{{ t("dialog.users.admin.value") }}</label></strong>
                                <span> <label>{{ t("dialog.users.admin.description") }}</label></span>
                            </div>

                            <div class="info-row">
                                <strong>{{ t("dialog.users.user.value") }}</strong>
                                <span>{{ t("dialog.users.user.description") }}</span>
                            </div>

                        </div>

                        <Divider />

                        <div class="info-section">

                            <h4>{{ t("dialog.users.info2_title") }}</h4>

                            <ul>
                                <li>{{ t("dialog.users.info2_hint") }}</li>
                                <li>{{ t("dialog.users.info2_hint2") }}</li>
                                <li>{{ t("dialog.users.info2_hint3") }}</li>
                            </ul>

                        </div>

                    </template>

                </Card>

            </div>

        </div>

        <template #footer>

            <Button
                :label="t('common.cancel')"
                severity="secondary"
                outlined
                @click="emit('update:visible', false)"
            />

            <Button
                :label="t('dialog.users.create')"
                icon="pi pi-check"
                :loading="saving"
                @click="create"
            />

        </template>

    </Dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";

import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import InputNumber from "primevue/inputnumber";
import Checkbox from "primevue/checkbox";
import Select from "primevue/select";
import Divider from "primevue/divider";
import Button from "primevue/button";
import Card from "primevue/card";

import type { Group } from "@/types/group";
import { createUser } from "@/api/users";

const { t } = useI18n();

const props = defineProps<{
    visible: boolean;
    groups: Group[];
}>();

const emit = defineEmits<{
    (e: "update:visible", value: boolean): void;
    (e: "created"): void;
}>();

const saving = ref(false);

const groupOptions = computed(() => [
    {
        id: null,
        name: "-- Без группы --",
    },
    ...props.groups,
]);

const form = reactive({
    username: "",
    password: "",
    config_limit: 5,
    group_id: null as number | null,
    is_admin: false,
});

async function create() {

    saving.value = true;

    try {

        await createUser({
            username: form.username,
            password: form.password,
            config_limit: form.config_limit,
            group_id: form.group_id,
            is_admin: form.is_admin,
        });

        form.username = "";
        form.password = "";
        form.config_limit = 5;
        form.group_id = null;
        form.is_admin = false;

        emit("created");
        emit("update:visible", false);

    } finally {

        saving.value = false;

    }

}
</script>

<style scoped>
.layout {
    display: grid;
    grid-template-columns: minmax(0, 2fr) minmax(280px, 1fr);
    gap: 2rem;
}

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
    font-size: .85rem;
}

.administrator {
    flex-direction: row;
    align-items: center;
    gap: .75rem;
}

.info {
    align-self: start;
}

.info-section {
    display: flex;
    flex-direction: column;
    gap: .75rem;
}

.info-row {
    display: flex;
    flex-direction: column;
    gap: .15rem;
}

.info h4 {
    margin: 0;
}

.info ul {
    margin: 0;
    padding-left: 1.2rem;
}

.info li {
    margin-bottom: .5rem;
}
:deep(.p-dialog-footer) {

    display: flex;

    gap: .75rem;

}

:deep(.p-dialog-header-actions .p-dialog-close-button) {
    border-radius: 0;
}

:deep(.p-dialog-header-actions .p-dialog-close-button .p-icon) {
    width: 1rem;
    height: 1rem;
}
@media (max-width: 1100px) {

    .layout {
        grid-template-columns: 1fr;
    }

    .info {
        margin-top: 1rem;
    }

}


@media (max-width: 600px) {

    .layout {
        gap: 1rem;
    }


    .form {
        gap: 1rem;
    }


    .administrator {
        align-items: flex-start;
    }


    .info-row {
        gap: .5rem;
    }

}
@media (max-width: 600px) {

    :deep(.p-dialog-footer) {

        flex-direction: column-reverse;

    }

    :deep(.p-dialog-footer .p-button) {

        width: 100%;

    }

}
@media (max-width: 600px) {

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
</style>