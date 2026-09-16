<template>
    <div class="dashboard-page">
    <PageHeader
        :title="t('menu.dashboard')"
        :subtitle="t('pages.dashboard.description')"
    />

    <!-- Statistics -->
    <div class="stats-grid">

        <Card class="stat-card">
            <template #title>👥 {{ t("menu.users") }}</template>
            <template #content>
                <div class="stat-value">
                    {{ users.length }}
                </div>
            </template>
        </Card>

        <Card class="stat-card">
            <template #title>🗂 {{ t("menu.groups") }}</template>
            <template #content>
                <div class="stat-value">
                    {{ groups.length }}
                </div>
            </template>
        </Card>

        <Card class="stat-card">
            <template #title>🖥 {{ t("menu.endpoints") }}</template>
            <template #content>
                <div class="stat-value">
                    {{ endpoints.length }}
                </div>
            </template>
        </Card>

        <Card class="stat-card">
            <template #title>📄 {{ t("menu.configs") }}</template>
            <template #content>
                <div class="stat-value">
                    {{ configurations.length }}
                </div>
            </template>
        </Card>

    </div>

    <div class="dashboard-grid">

        <!-- Quick actions -->
        <Card>

            <template #title>
                ⚡ {{ t("pages.dashboard.quickActions") }}
            </template>

            <template #content>
                <div class="actions-grid">

                    <div class="action-card" @click="dialogs.openUserCreate()">
                        <i class="pi pi-user-plus action-icon" />
                        <div class="action-title">{{ t("pages.users.create") }}</div>
                    </div>

                    <div class="action-card" @click="dialogs.openEndpointCreate()">
                        <i class="pi pi-server action-icon" />
                        <div class="action-title">{{ t("pages.endpoints.create") }}</div>
                    </div>

                    <div class="action-card" @click="dialogs.openConfigCreate()">
                        <i class="pi pi-file-plus action-icon" />
                        <div class="action-title">{{ t("pages.configs.create") }}</div>
                    </div>

                    <div class="action-card secondary" @click="dialogs.openGroupCreate()">
                        <i class="pi pi-users action-icon" />
                        <i class="pi pi-plus plus-overlay" />
                        <div class="action-title">{{ t("pages.groups.create") }}</div>
                    </div>

                    <div class="action-card secondary" @click="router.push('/admin/configs')">
                        <i class="pi pi-folder-open action-icon" />
                        <div class="action-title">{{ t("menu.configs") }}</div>
                    </div>

                    <div class="action-card secondary" @click="router.push('/admin/users')">
                        <i class="pi pi-users action-icon" />
                        <div class="action-title">{{ t("menu.users") }}</div>
                    </div>
                    <div class="action-card secondary" @click="router.push('/admin/profile')">
                        <i class="pi pi-user action-icon" />
                        <div class="action-title">{{t('user.quickactions.profile')}}</div>
                    </div>

                </div>

            </template>

        </Card>

        <!-- System -->
        <Card>

            <template #title>
                🖥 {{ t("pages.dashboard.systemStatus") }}
            </template>

            <template #content>

                <div class="status">
                    <li>В разработке</li>
                    <Divider />

                    <Tag severity="info" :value="t('common.inDevelopment')" />

                </div>

            </template>

        </Card>

    </div>

    <!-- Activity -->
    <Card>

        <template #title>
            🕒 {{ t("pages.dashboard.recentActivity") }}
        </template>

        <template #content>

            <ul v-if="activities.length" class="activity">
                <li v-for="activity in activities" :key="activity.id">
                    <div class="activity-message">
                        <template v-if="activity.action === 'auth.login'">
                            {{ activity.message }}
                        </template>
                        <template v-else>
                            <strong>{{ getActivityUser(activity) }}</strong>
                            {{ activity.message }}
                        </template>
                    </div>
                    <div class="activity-meta">
                        {{ activity.action }} · {{ formatActivityDate(activity.created_at) }}
                    </div>
                </li>
            </ul>

            <div v-else class="activity-empty">
                Нет активности
            </div>

        </template>

</Card>

    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from 'vue-router';

import Card from "primevue/card";
import Button from "primevue/button";
import PageHeader from "@/components/common/PageHeader.vue";
import { useUsers } from "@/composables/useUsers";
import { useGroups } from "@/composables/useGroups";
import { useEndpoints } from "@/composables/useEndpoints";
import { useConfigurations } from "@/composables/useConfigs";
import { useDialogsStore } from "@/stores/dialog";
import { getActivity, type Activity } from "@/api/users";

const activities = ref<Activity[]>([]);
const dialogs = useDialogsStore();
const router = useRouter();
const { t } = useI18n({
    useScope: "global",
});

const {
    users,
    loadUsers,
} = useUsers();

const {
    groups,
    loadGroups,
} = useGroups();

const {
    endpoints,
    loadEndpoints,
} = useEndpoints();

const {
    configurations,
    loadConfigurations,
} = useConfigurations();

async function loadActivity() {
    activities.value = await getActivity(10);
}

function formatActivityDate(value?: string | null) {
    if (!value) {
        return "";
    }

    const normalized = value.replace(/\.(\d{3})\d+Z$/, ".$1Z");
    const date = new Date(normalized);

    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return date.toLocaleString();
}

function getActivityUser(activity: Activity) {
    if (activity.action === "auth.login") {
        return "";
    }

    if (activity.user_id == null) {
        return "Система:";
    }

    const user = users.value.find((item) => item.id === activity.user_id);

    return user ? `${user.username}:` : `Пользователь #${activity.user_id}:`;
}

onMounted(async () => {
    await loadUsers();
    await loadGroups();
    await loadEndpoints();
    await loadConfigurations();
    await loadActivity();
});
</script>

<style scoped>

.stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 1.5rem;
}

.stat-card {
    cursor: pointer;
    transition: .2s;
}

.stat-card:hover {
    transform: translateY(-2px);
}

.stat-value {
    font-size: 2rem;
    font-weight: 700;
    margin-top: .5rem;
}

.dashboard-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
}

.actions-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
}

.action-card {
    border: 1px solid var(--p-content-border-color);
    border-radius: 14px;
    padding: 1.2rem;
    cursor: pointer;

    display: flex;
    flex-direction: column;
    gap: .5rem;

    transition: all .2s ease;
    background: var(--p-surface-card);
}

.action-card:hover {
    transform: translateY(-3px);
    border-color: var(--p-primary-color);
    box-shadow: var(--p-overlay-modal-shadow);
}

.action-card.secondary {
    background: var(--p-surface-ground);
}

.action-icon {
    font-size: 2rem;
    color: var(--p-primary-color);
}

.action-card {
    position: relative;
}

.plus-overlay {
    position: absolute;
    top: 0.9rem;
    left: 2.2rem;
    font-size: 0.8rem;
    color: var(--p-primary-color);
    font-weight: bold;
}

.action-title {
    font-weight: 600;
    font-size: 1rem;
}

.action-subtitle {
    font-size: .85rem;
    color: var(--p-text-muted-color);
}

@media (max-width: 768px) {
    .actions-grid {
        grid-template-columns: 1fr;
    }
}

.status {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.activity {
    margin: 0;
    padding: 0;
    list-style: none;
}

.activity li {
    margin-bottom: .9rem;
    padding-bottom: .9rem;
    border-bottom: 1px solid var(--p-content-border-color);
}

.activity li:last-child {
    margin-bottom: 0;
    padding-bottom: 0;
    border-bottom: 0;
}

.activity-message {
    font-weight: 500;
}

.activity-meta {
    margin-top: .25rem;
    font-size: .8rem;
    color: var(--p-text-muted-color);
}

.activity-empty {
    color: var(--p-text-muted-color);
}

@media (max-width: 900px) {

    .stats-grid {
        grid-template-columns: repeat(2, 1fr);
    }

    .dashboard-grid {
        grid-template-columns: 1fr;
    }

}

@media (max-width: 640px) {

    .stats-grid {
        grid-template-columns: 1fr;
    }

}
</style>