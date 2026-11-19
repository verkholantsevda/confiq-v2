<template>
    <PageHeader
        :title="t('menu.backups')"
        :subtitle="t('pages.backups.description')"
    >
        <template #actions>
            <Button
                :label="t('pages.backups.create')"
                icon="pi pi-download"
                :loading="loading"
                @click="create"
            />
            <Button
                :label="t('pages.backups.restore')"
                icon="pi pi-upload"
                severity="secondary"
                :disabled="restoring"
                @click="openRestore"
            />
        </template>
    </PageHeader>

    <div class="table-wrapper">
        <EntityTable
            :value="backups"
            :loading="loading"
            :empty-message="t('pages.backups.empty')"
            :search-placeholder="t('common.search')"
        >
            <Column field="name" :header="t('pages.backups.name')" sortable />
            <Column field="created_at" :header="t('common.created')" sortable>
                <template #body="{ data }">
                    {{ new Date(data.created_at).toLocaleString() }}
                </template>
            </Column>
            <Column field="size" :header="t('pages.backups.size')" sortable>
                <template #body="{ data }">
                    {{ formatSize(data.size) }}
                </template>
            </Column>
            <Column :header="t('common.actions')" style="width:100px">
                <template #body="{ data }">
                    <Button
                        icon="pi pi-undo"
                        text
                        rounded
                        :title="t('pages.backups.restoreFromBackup')"
                        :aria-label="t('pages.backups.restoreFromBackup')"
                        @click="confirmRestore(data.name)"
                    />
                    <Button
                        icon="pi pi-download"
                        text
                        rounded
                        @click="download(data.name)"
                    />
                    <Button
                        icon="pi pi-trash"
                        severity="danger"
                        text
                        rounded
                        @click="confirmDelete(data.name)"
                    />
                </template>
            </Column>
        </EntityTable>
    </div>

    <input
        ref="fileInput"
        type="file"
        accept=".db,.sqlite,.sqlite3"
        style="display: none"
        @change="onFileChange"
    />

    <Dialog
        v-model:visible="showDeleteDialog"
        modal
        :header="t('pages.backups.deleteTitle')"
        :style="{ width: '420px' }"
    >
        <p>{{ t('pages.backups.deleteConfirm') }}</p>
        <template #footer>
            <Button
                :label="t('common.cancel')"
                severity="secondary"
                @click="showDeleteDialog = false"
            />
            <Button
                :label="t('common.confirm')"
                severity="danger"
                :loading="loading"
                @click="doDelete"
            />
        </template>
    </Dialog>

    <Dialog
        v-model:visible="showRestoreDialog"
        modal
        :header="t('pages.backups.restoreTitle')"
        :style="{ width: '460px' }"
    >
        <p>{{ t('pages.backups.restoreConfirm', { name: pendingFile?.name ?? pendingRestoreName ?? '' }) }}</p>
        <p class="restore-warning">{{ t('pages.backups.restoreWarning') }}</p>
        <template #footer>
            <Button
                :label="t('common.cancel')"
                severity="secondary"
                @click="showRestoreDialog = false"
            />
            <Button
                :label="t('common.confirm')"
                severity="danger"
                :loading="restoring"
                @click="doRestore"
            />
        </template>
    </Dialog>

    <Dialog
        v-model:visible="showResultDialog"
        modal
        :header="t('pages.backups.resultTitle')"
        :style="{ width: '420px' }"
    >
        <p>{{ resultMessage }}</p>
        <template #footer>
            <Button :label="t('common.confirm')" @click="reloadPage" />
        </template>
    </Dialog>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

import Button from "primevue/button";
import Column from "primevue/column";
import Dialog from "primevue/dialog";

import PageHeader from "@/components/common/PageHeader.vue";
import EntityTable from "@/components/common/EntityTable.vue";

import { useBackups } from "@/composables/useBackups";
import {
    downloadBackup,
    restoreBackup,
    restoreFromBackup,
} from "@/api/backups";

const { t } = useI18n({ useScope: "global" });

const { backups, loading, load, create, remove } = useBackups();

const fileInput = ref<HTMLInputElement | null>(null);
const pendingFile = ref<File | null>(null);
const pendingName = ref<string | null>(null);
const pendingRestoreName = ref<string | null>(null);
const restoring = ref(false);

const showDeleteDialog = ref(false);
const showRestoreDialog = ref(false);
const showResultDialog = ref(false);
const resultMessage = ref("");

onMounted(load);

function formatSize(bytes: number) {
    if (bytes < 1024) {
        return `${bytes} B`;
    }

    if (bytes < 1024 * 1024) {
        return `${(bytes / 1024).toFixed(1)} KB`;
    }

    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function download(name: string) {
    void downloadBackup(name);
}

function confirmDelete(name: string) {
    pendingName.value = name;
    showDeleteDialog.value = true;
}

function confirmRestore(name: string) {
    pendingFile.value = null;
    pendingRestoreName.value = name;
    showRestoreDialog.value = true;
}

async function doDelete() {
    if (!pendingName.value) {
        return;
    }

    showDeleteDialog.value = false;

    try {
        await remove(pendingName.value);
    } finally {
        pendingName.value = null;
    }
}

function openRestore() {
    fileInput.value?.click();
}

function onFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (!file) {
        return;
    }

    pendingFile.value = file;
    pendingRestoreName.value = null;
    input.value = "";
    showRestoreDialog.value = true;
}

async function doRestore() {
    if (!pendingFile.value && !pendingRestoreName.value) {
        return;
    }

    restoring.value = true;
    showRestoreDialog.value = false;

    try {
        if (pendingFile.value) {
            await restoreBackup(pendingFile.value);
        } else if (pendingRestoreName.value) {
            await restoreFromBackup(pendingRestoreName.value);
        }

        resultMessage.value = t("pages.backups.restoreSuccess");
    } catch {
        resultMessage.value = t("pages.backups.restoreFailed");
    } finally {
        pendingFile.value = null;
        pendingRestoreName.value = null;
        restoring.value = false;
        showResultDialog.value = true;
    }
}

function reloadPage() {
    window.location.reload();
}
</script>

<style scoped>
.table-wrapper {
    width: 100%;
}

.restore-warning {
    margin-top: 0.75rem;
    color: var(--p-text-muted-color);
}

@media (max-width: 768px) {
    .table-wrapper {
        overflow-x: auto;
        -webkit-overflow-scrolling: touch;
    }

    .table-wrapper :deep(table) {
        min-width: 640px;
    }
}
</style>
