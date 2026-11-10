<template>
  <Dialog
    :visible="visible"
    :header="endpoint?.name"
    :modal="true"
    :closable="true"
    maximizable
    :breakpoints="{ '1200px': '95vw', '960px': '98vw', '640px': '100vw' }"
    :style="isMaximized ? { width: '100vw', height: '100vh', maxWidth: 'none' } : { width: '90vw', maxWidth: '1200px' }"
    @maximize="isMaximized = true"
    @unmaximize="isMaximized = false"
    @update:visible="emit('update:visible', $event)"
  >
    <div class="dialog-content">
      <div v-if="endpoint?.address" class="endpoint-address">
        {{ endpoint.address }}:{{ endpoint.port }}
      </div>
      <div class="toolbar">
        <span class="count">{{ t('dialog.endpoint_configs.count', { count: selected.length, total: filteredConfigs.length }) }}</span>
        <div v-if="selected.length" style="display:flex;gap:.5rem">
          <Button
            :label="t('dialog.endpoint_configs.delete_selected')"
            severity="danger"
            icon="pi pi-trash"
            @click="confirmDeleteVisible = true"
          />
          <Button
            :label="t('dialog.endpoint_configs.move_selected')"
            severity="primary"
            icon="pi pi-arrow-right"
            @click="moveDialogVisible = true"
          />
        </div>
      </div>

      <DataTable
        :value="filteredConfigs"
        v-model:selection="selected"
        :loading="loading"
        dataKey="id"
        stripedRows
        responsiveLayout="scroll"
        scrollable
        scrollHeight="calc(100vh - 220px)"
      >
        <Column selectionMode="multiple" style="width:3rem" />
        <Column field="username" :header="t('dialog.endpoint_configs.column_user')" sortable>
          <template #body="{ data }">
            {{ usernameById[data.user_id] || data.username || '-' }}
          </template>
        </Column>
        <Column field="name" :header="t('dialog.endpoint_configs.column_config')" sortable />
        <Column field="created_at" :header="t('dialog.endpoint_configs.column_created')" sortable>
          <template #body="{ data }">
            {{ formatDate(data.created_at) }}
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Подтверждение удаления -->
    <Dialog
      v-model:visible="confirmDeleteVisible"
      :header="t('dialog.endpoint_configs.confirm_delete')"
      modal
      :style="{ width: '30rem' }"
      @update:visible="confirmDeleteVisible = $event"
    >
      <template #footer>
        <Button
          :label="t('common.cancel')"
          severity="secondary"
          outlined
          @click="confirmDeleteVisible = false"
        />
        <Button
          :label="t('dialog.endpoint_configs.delete_selected')"
          severity="danger"
          icon="pi pi-trash"
          :loading="deleting"
          @click="deleteSelected"
        />
      </template>
    </Dialog>

    <!-- Перемещение на другой endpoint -->
    <Dialog
      v-model:visible="moveDialogVisible"
      :header="t('dialog.endpoint_configs.move_to')"
      modal
      :style="{ width: '30rem' }"
      @update:visible="moveDialogVisible = $event"
    >
      <div style="display:flex;flex-direction:column;gap:1rem">
        <div class="field">
          <label>{{ t('dialog.endpoint_configs.endpoint') }}</label>
          <select v-model="moveEndpointId" class="endpoint-select">
            <option v-for="e in endpoints" :key="e.id" :value="e.id">{{ e.name }} ({{ e.address }}:{{ e.port }})</option>
          </select>
        </div>
        <div class="hint">{{ t('dialog.endpoint_configs.confirm_move') }}</div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <Button
            :label="t('common.cancel')"
            severity="secondary"
            outlined
            @click="moveDialogVisible = false"
          />
          <Button
            :label="t('common.confirm')"
            severity="primary"
            :loading="moving"
            @click="confirmMove"
          />
        </div>
      </template>
    </Dialog>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from "vue-i18n";
import Dialog from 'primevue/dialog';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';

import { getAllConfigurations, deleteConfiguration, moveConfigurations } from '@/api/configs';
import { getUsers } from '@/api/users';
import { getEndpoints } from '@/api/endpoints';

const props = defineProps<{
  visible: boolean;
  endpoint: any | null;
}>();
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
}>();

const { t } = useI18n({ useScope: 'global' });

const loading = ref(false);
const deleting = ref(false);
const moving = ref(false);
const allConfigs = ref<any[]>([]);
const users = ref<any[]>([]);
const endpoints = ref<any[]>([]);
const moveDialogVisible = ref(false);
const confirmDeleteVisible = ref(false);
const moveEndpointId = ref<number | null>(null);

const filteredConfigs = computed(() => {
  if (!props.endpoint) return [];
  return allConfigs.value.filter((c: any) => c?.endpoint?.id === props.endpoint.id);
});

const selected = ref<any[]>([]);
const selectAll = ref(false);
const isMaximized = ref(false);

const formatDate = (value?: string) => {
  if (!value) return '-';
  return new Date(value).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
};

onMounted(async () => {
  loading.value = true;
  try {
    const [configsData, usersData, endpointsData] = await Promise.all([
      getAllConfigurations(),
      getUsers(),
      getEndpoints(),
    ]);
    allConfigs.value = configsData;
    users.value = usersData;
    endpoints.value = endpointsData;
  } catch {
    allConfigs.value = [];
    users.value = [];
    endpoints.value = [];
  } finally {
    loading.value = false;
  }
});

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      selected.value = [];
      selectAll.value = false;
      isMaximized.value = false;
      moveEndpointId.value = null;
    }
  }
);

const usernameById = computed(() => {
  const map: Record<number, string> = {};
  users.value.forEach((u: any) => { map[u.id] = u.username; });
  return map;
});

async function deleteSelected() {
  deleting.value = true;
  try {
    for (const c of selected.value) {
      await deleteConfiguration(c.id);
    }
    selected.value = [];
    confirmDeleteVisible.value = false;
    const data = await getAllConfigurations();
    allConfigs.value = data;
  } catch {
    // ошибки уже логируются в api
  } finally {
    deleting.value = false;
  }
}

async function confirmMove() {
  if (!moveEndpointId.value) return;
  moving.value = true;
  try {
    await moveConfigurations(selected.value.map((c: any) => c.id), moveEndpointId.value);
    moveDialogVisible.value = false;
    selected.value = [];
    const data = await getAllConfigurations();
    allConfigs.value = data;
  } catch {
    // ошибки уже логируются в api
  } finally {
    moving.value = false;
  }
}

watch(
  () => selected.value,
  () => {
    const ids = new Set(selected.value.map((s: any) => s.id));
    selectAll.value = filteredConfigs.value.length > 0 && filteredConfigs.value.every((c: any) => ids.has(c.id));
  }
);
</script>

<style scoped>
.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.9rem;
}

.count {
  font-weight: 600;
}

.endpoint-address {
  font-size: 0.9rem;
  color: var(--p-text-muted-color);
  margin: 0 0 0.75rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.field label {
  font-size: 0.9rem;
  font-weight: 600;
}

.endpoint-select {
  padding: 0.5rem;
  font-size: 0.95rem;
}

.hint {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

@media (max-width: 640px) {
  .dialog-content .toolbar {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}
</style>
