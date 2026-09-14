<template>
  <Dialog
    :visible="visible"
    @hide="handleCancel"
    :header="t('dialog.endpoints.name_dialog_edit')"
    :modal="true"
    :closable="false"
    class="endpoint-create-dialog"
    maximizable
    @update:visible="emit('update:visible', $event)"
  >
    <div class="dialog-content">
      <form @submit.prevent="handleSave" class="form-left">
        <div class="field">
          <label for="name">{{ t('dialog.endpoints.name') }}</label>
          <InputText id="name" v-model="form.name" :placeholder="t('dialog.endpoints.name_placeholder')" required autofocus />
        </div>

        <div class="field">
          <label for="address">{{ t('dialog.endpoints.address') }}</label>
          <InputText id="address" v-model="form.address" :placeholder="t('dialog.endpoints.address_placeholder')" required />
          <small class="hint">{{ t('dialog.endpoints.address_hint') }}</small>
        </div>

        <div class="field">
          <label for="port">{{ t('dialog.endpoints.port') }}</label>
          <InputNumber id="port" v-model="form.port" :placeholder="t('dialog.endpoints.port_placeholder')" :min="1" :max="65535" showButtons :step="1" required />
          <small class="hint">{{ t('dialog.endpoints.port_hint') }}</small>
        </div>

        <div class="field">
          <label for="configTypes">{{ t('dialog.endpoints.configuration_types') }}</label>
          <MultiSelect
            id="configTypes"
            v-model="form.config_type_ids"
            :options="configTypes"
            optionLabel="name"
            optionValue="id"
            :placeholder="t('dialog.endpoints.config_types_placeholder')"
            display="chip"
            :filter="true"
            :disabled="loadingConfigTypes"
            required
          />
          <small class="hint">{{ t('dialog.endpoints.config_types_hint') }}</small>
        </div>

        <div class="field groups-field">
          <label>{{ t('dialog.endpoints.groups') }}</label>
          <div class="groups-checkboxes">
            <div v-for="group in groups" :key="group.id" class="group-checkbox">
              <Checkbox v-model="form.group_ids" :value="group.id" :inputId="'group-' + group.id" />
              <label :for="'group-' + group.id">{{ group.name }}</label>
            </div>
          </div>
          <small class="hint">{{ t('dialog.endpoints.groups_hint') }}</small>
        </div>

        <div class="dialog-footer">
          <Button :label="t('common.cancel')" icon="pi pi-times" class="p-button-text" @click="handleCancel" :disabled="loading" />
          <Button :label="t('dialog.endpoints.submit')" icon="pi pi-check" type="submit" :loading="loading" :disabled="loading || !isFormValid" />
        </div>
      </form>

      <div class="info-right">
        <div class="info-card">
          <h3>{{ t('dialog.endpoints.examples_title') }}</h3>
          <ul>
            <li>http://example.com</li>
            <li>https://api.example.com:8080</li>
            <li>192.168.1.100</li>
            <li>10.0.0.5:3000</li>
          </ul>
        </div>

        <div class="info-card">
          <h3>{{ t('dialog.endpoints.config_types_title') }}</h3>
          <p>{{ t('dialog.endpoints.config_types_description') }}</p>
          <ul>
            <li v-for="configType in configTypes" :key="configType.id">
              <strong>{{ configType.name }}</strong>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue';
import { useI18n } from "vue-i18n";
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import MultiSelect from 'primevue/multiselect';
import Checkbox from 'primevue/checkbox';
import Button from 'primevue/button';

import { updateEndpoint } from '@/api/endpoints';
import { getGroups } from '@/api/groups';
import { getConfigTypes } from '@/api/configTypes';
const { t } = useI18n({ useScope: 'global' });
interface Group {
  id: number;
  name: string;
}

interface ConfigType {
  id: number;
  name: string;
}

const props = defineProps<{
  visible: boolean;
  endpoint: any | null;
}>();
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'updated'): void;
}>();

const loading = ref(false);
const loadingGroups = ref(false);
const loadingConfigTypes = ref(false);

const groups = ref<Group[]>([]);
const configTypes = ref<ConfigType[]>([]);

const form = reactive({
  id: 0,
  name: '',
  address: '',
  port: 2408,
  group_ids: [] as number[],
  config_type_ids: [] as number[],
});

const resetForm = () => {
  form.id = 0;
  form.name = '';
  form.address = '';
  form.port = 2408;
  form.group_ids = [];
  form.config_type_ids = [];
};

const loadGroups = async () => {
  loadingGroups.value = true;
  try {
    const data = await getGroups();
    groups.value = data;
  } finally {
    loadingGroups.value = false;
  }
};

const loadConfigTypes = async () => {
  loadingConfigTypes.value = true;
  try {
    const data = await getConfigTypes();
    configTypes.value = data;
  } finally {
    loadingConfigTypes.value = false;
  }
};

watch(
  [() => props.visible, () => props.endpoint],
  ([visible, endpoint]) => {
    if (visible && endpoint) {
      loadGroups();
      loadConfigTypes();
      form.id = endpoint.id;
      form.name = endpoint.name;
      form.address = endpoint.address;
      form.port = endpoint.port;
      form.group_ids = endpoint.group_ids ?? [];
      form.config_type_ids = endpoint.config_type_ids ?? [];
    }
  }
);

const isFormValid = computed(() => {
  return (
    form.name.trim().length > 0 &&
    form.address.trim().length > 0 &&
    form.port >= 1 &&
    form.port <= 65535 &&
    form.config_type_ids.length > 0
  );
});

const handleSave = async () => {
  if (!isFormValid.value) return;
  loading.value = true;
  try {
    await updateEndpoint(form.id, {
      name: form.name.trim(),
      address: form.address.trim(),
      port: form.port,
      group_ids: form.group_ids,
      config_type_ids: form.config_type_ids,
    });
    emit('updated');
    emit('update:visible', false);
    resetForm();
  } catch (e) {
    // handle error here if needed
  } finally {
    loading.value = false;
  }
};

const handleCancel = () => {
  emit('update:visible', false);
};
</script>

<style scoped>
.endpoint-create-dialog .dialog-content {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
}

.form-left {
  flex: 1 1 300px;
  max-width: 500px;
  display: flex;
  flex-direction: column;
}

.field {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
}

label {
  font-weight: 600;
  margin-bottom: 0.3rem;
}

.hint {
  font-size: 0.875rem;
  color: var(--text-color-secondary, #6c757d);
  margin-top: 0.25rem;
}

.groups-field .groups-checkboxes {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid var(--surface-border, #dcdcdc);
  padding: 0.5rem;
  border-radius: 4px;
  background: var(--surface-card-100);
}

.group-checkbox {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  min-width: 120px;
}

.dialog-footer {
  margin-top: auto;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.info-right {
  flex: 1 1 280px;
  max-width: 350px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-card {
  background: var(--surface-100);
  border-radius: 6px;
  padding: 1rem 1.5rem;
  box-shadow: var(--surface-shadow, 0 2px 8px rgba(0,0,0,0.1));
}


.info-card h3 {
  margin-top: 0;
  margin-bottom: 0.75rem;
  font-weight: 700;
  font-size: 1.1rem;
}

.info-card ul {
  margin: 0;
  padding-left: 1.25rem;
  list-style-type: disc;
}

@media (max-width: 768px) {
  .endpoint-create-dialog .dialog-content {
    flex-direction: column;
  }
  .form-left,
  .info-right {
    max-width: 100%;
  }
}
</style>
