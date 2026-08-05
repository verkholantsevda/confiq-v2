<template>
    <Dialog
        :visible="visible"
        modal
        :style="{ width: '52rem' }"
        :header="t('pages.users.edit')"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="layout">

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
                        :min="1"
                        :max="100"
                        showButtons
                        fluid
                    />
                </div>

                <div class="field">
                    <label>{{ t("dialog.users.group.value") }}</label>

                    <Select
                        v-model="form.group_id"
                        :options="groups"
                        optionLabel="name"
                        optionValue="id"
                        fluid
                    >
                        <template #value="{ value }">
                            <span v-if="value === null">
                                --Без группы--
                            </span>

                            <span v-else>
                                {{
                                    groups.find(g => g.id === value)?.name
                                }}
                            </span>
                        </template>
                    </Select>

                    <small class="hint">
                        <label>{{ t("dialog.users.group.description") }}</label>
                    </small>
                </div>

                <Divider />

                <div class="field checkbox">
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

            <Card class="info">

                <template #title>
                     <label>{{ t("dialog.users.info_title") }}</label>
                </template>
                <template #content>

                    <div class="info-section">

                        <label>{{ t("dialog.users.info_hint") }}</label>

                        <div class="info-row">
                            <strong> <label>{{ t("dialog.users.admin.value") }}</label></strong>
                            <span> - </span>
                            <span> <label>{{ t("dialog.users.admin.description") }}</label></span>
                        </div>

                        <div class="info-row">
                            <strong>{{ t("dialog.users.user.value") }}</strong>
                            <span> - </span>
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

        <template #footer>

            <Button
                :label="t('common.cancel')"
                severity="secondary"
                outlined
                @click="emit('update:visible', false)"
            />

            <Button
                :label="t('common.save')"
                icon="pi pi-check"
                :loading="saving"
                @click="save"
            />

        </template>

    </Dialog>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from "vue";
import { useI18n } from "vue-i18n";

import Dialog from "primevue/dialog";
import Card from "primevue/card";
import Divider from "primevue/divider";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import InputNumber from "primevue/inputnumber";
import Select from "primevue/select";
import Checkbox from "primevue/checkbox";
import Button from "primevue/button";

import type { Group } from "@/types/group";
import type { User, UpdateUserRequest } from "@/types/user";

import { updateUser } from "@/api/users";

const { t } = useI18n();

const props = defineProps<{
    visible: boolean;
    user: User | null;
    groups: Group[];
}>();

const emit = defineEmits<{
    (e: "update:visible", value: boolean): void;
    (e: "updated"): void;
}>();

const saving = ref(false);

const form = reactive<UpdateUserRequest & { id: number }>({
    id: 0,
    username: "",
    password: "",
    config_limit: 5,
    group_id: null,
    is_admin: false,
});

watch(
    () => props.user,
    (user) => {
        if (!user) return;

        form.id = user.id;
        form.username = user.username;
        form.password = "";
        form.config_limit = user.config_limit;
        form.group_id = user.group_id ?? null;
        form.is_admin = user.is_admin;
    },
    {
        immediate: true,
    }
);

async function save() {
    if (!form.id) return;

    saving.value = true;

    try {

        const payload: UpdateUserRequest = {
            username: form.username,
            password: form.password,
            config_limit: form.config_limit,
            group_id: form.group_id,
            is_admin: form.is_admin,
        };

        if (form.password.trim() !== "") {
            payload.password = form.password;
        }

        await updateUser(form.id, payload);

        emit("updated");
        emit("update:visible", false);

    } finally {
        saving.value = false;
    }
}
</script>

<style scoped>
.layout {
    display: grid;
    grid-template-columns: 2fr 1fr;
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

.checkbox {
    flex-direction: row;
    align-items: center;
    gap: .75rem;
}

.hint {
    color: var(--text-color-secondary);
    font-size: .85rem;
}

.info {
    align-self: start;
}

.info h4 {
    margin: 0 0 .5rem;
}

.info ul {
    margin: 0;
    padding-left: 1rem;
}

.info li {
    margin-bottom: .5rem;
}
</style>