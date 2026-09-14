<template>
    <Dialog
        :visible="visible"
        modal
        maximizable
        :style="{ width: '42rem' }"
        :header="t('dialog.groups.name_edit_dialog')"
        @update:visible="emit('update:visible', $event)"
    >
        <div class="form">

            <div class="field">
                <label>{{ t("dialog.groups.name_edit_dialog") }}</label>

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

            <div class="group-layout">
                <div class="group-stats">
                    <div class="stats-title"> <label>{{ t("dialog.groups.stats") }}</label></div>

                    <div class="stat-field">
                        <span>{{ t("dialog.groups.count_users") }}</span>
                        <InputText :value="String(usersCount)" readonly />
                    </div>

                    <div class="stat-field">
                        <span>{{ t("dialog.groups.count_endpoints") }}</span>
                        <InputText :value="String(endpoints.length)" readonly />
                    </div>

                    <div class="stat-field">
                        <span> <label>{{ t("dialog.groups.date_created") }}</label></span>
                        <InputText :value="formatDate((props.group as any)?.created_at)" readonly />
                    </div>

                    <Divider />

                    <div class="users-title"> <label>{{ t("dialog.groups.users_groups") }}</label></div>

                    <div class="users-list">
                        <div
                            v-for="user in groupUsers"
                            :key="user.id"
                            class="user-row"
                        >
                            <i class="pi pi-user user-icon" />
                            <span>{{ user.username }}</span>
                        </div>

                        <div v-if="!groupUsers.length && usersCount === 0" class="empty-users">
                            {{ t("dialog.groups.no_users") }}
                        </div>
                    </div>
                </div>
                <div class="group-endpoints">
                    <div class="field">
                         <label>{{ t("dialog.groups.endpoints") }}</label>

                        <small class="hint">
                            <label>{{ t("dialog.groups.endpoints_hint") }}</label>
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
import { getEndpointsAll } from "@/api/endpoints";
import { getUsers } from "@/api/users";

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
const usersCount = ref(0);
const groupUsers = ref<Array<{ id: number | string; username: string }>>([]);
const allUsers = ref<Array<{ id: number | string; username: string; group_id?: number | null }>>([]);

function updateGroupUsers(group: Group | null) {
    if (!group) {
        groupUsers.value = [];
        usersCount.value = 0;
        return;
    }

    const rawUsers = (group as any).users;
    const embeddedUsers = Array.isArray(rawUsers)
        ? rawUsers
        : rawUsers && typeof rawUsers === "object"
            ? Object.values(rawUsers)
            : [];

    if (embeddedUsers.length) {
        groupUsers.value = embeddedUsers.map((user: any, index: number) => ({
            id: user?.id ?? index,
            username: user?.username ?? user?.name ?? user?.login ?? user?.email ?? String(user?.id ?? index),
        }));
    } else {
        groupUsers.value = allUsers.value
            .filter((user) => Number(user.group_id) === Number(group.id))
            .map((user) => ({
                id: user.id,
                username: user.username,
            }));
    }

    if (groupUsers.value.length) {
        usersCount.value = groupUsers.value.length;
    } else if (typeof (group as any).users_count === "number") {
        usersCount.value = (group as any).users_count;
    } else if (Array.isArray((group as any).user_ids)) {
        usersCount.value = (group as any).user_ids.length;
        groupUsers.value = (group as any).user_ids.map((id: number | string) => ({
            id,
            username: String(id),
        }));
    } else {
        usersCount.value = 0;
    }
}

onMounted(async () => {
    endpoints.value = await getEndpointsAll();
    allUsers.value = await getUsers();

    if (props.visible && props.group) {
        updateGroupUsers(props.group);
    }
});

const form = reactive({
    id: 0,
    name: "",
    description: "",
    endpoint_ids: [] as number[],
});

watch(
  [() => props.group, () => props.visible, allUsers],
  ([group, visible]) => {
    if (!group || !visible) return;

    form.id = group.id;
    form.name = group.name;
    form.description = group.description ?? "";

    form.endpoint_ids.splice(0);
    if (group.endpoint_ids) {
        form.endpoint_ids.push(...group.endpoint_ids);
    }

    updateGroupUsers(group);
  },
  { immediate: true }
);

function formatDate(value?: string | Date | null) {
    if (!value) return "—";

    return new Intl.DateTimeFormat("ru-RU", {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
    }).format(new Date(value));
}

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

.group-layout {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    gap: 1.5rem;
    align-items: start;
}

.group-endpoints {
    min-width: 0;
}

.group-stats {
    display: flex;
    flex-direction: column;
    gap: .75rem;
}

.stats-title,
.users-title {
    font-size: 1rem;
    font-weight: 600;
}

.stat-field {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 7rem;
    align-items: center;
    gap: .75rem;
    padding: .75rem;
    border-radius: var(--p-content-border-radius);
    background: var(--p-content-hover-background);
}

.stat-field span {
    font-weight: 600;
}

.stat-field :deep(.p-inputtext) {
    width: 100%;
}

.users-list {
    display: flex;
    flex-direction: column;
    gap: .5rem;
    max-height: 180px;
    overflow-y: auto;
}

.user-row {
    display: flex;
    align-items: center;
    gap: .65rem;
    padding: .65rem .75rem;
    border-radius: var(--p-content-border-radius);
    background: var(--p-content-hover-background);
}

.user-icon {
    font-size: .9rem;
    color: var(--text-color-secondary);
}

.empty-users {
    padding: .65rem .75rem;
    color: var(--text-color-secondary);
    background: var(--p-content-hover-background);
    border-radius: var(--p-content-border-radius);
}
</style>