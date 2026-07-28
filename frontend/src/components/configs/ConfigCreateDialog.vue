<template>
  <Dialog
    v-model:visible="visible"
    modal :style="{ width: '520px' }"
    :header="t('dialog.config.name_dialog')"
  >
    <form class="form" @submit.prevent>
      <div class="field">
        <label for="name">{{ t("dialog.config.name") }}</label>
        <InputText id="name" v-model="form.name" />
      </div>
      <div class="field">
        <label for="endpoint">{{ t("dialog.config.endpoint") }}</label>
        <Select id="endpoint" v-model="form.endpoint_id" :options="endpoints" optionLabel="name" optionValue="id" />
      </div>
      <div class="field">
        <label for="configType">{{ t("dialog.config.configtype") }}</label>
        <Select
          id="configType"
          v-model="form.config_type_id"
          :options="configTypes"
          optionLabel="name"
          optionValue="id"
          :disabled="!form.endpoint_id"
        />
      </div>
      <Card class="limit-card">
        <div class="limit-row">
          Доступно конфигураций: {{ configLimit }}
        </div>
      </Card>
    </form>
    <template #footer>
      <Button :label="t('common.cancel')" text severity="secondary" @click="visible = false" />
      <Button
        :label="t('dialog.config.create')"
        icon="pi pi-check"
        @click="submitCreateConfig"
      />
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { useI18n } from "vue-i18n";
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Card from 'primevue/card';
import { getEndpoints, getEndpointConfigTypes } from '@/api/endpoints';
import { createConfiguration } from '@/api/configs';
const { t } = useI18n();
const props = defineProps<{
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e:'update:modelValue', value:boolean):void;
    (e:'created'):void;
}>();

const visible = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v),
});

const form = reactive({
    name: '',
    endpoint_id: null as number | null,
    config_type_id: null as number | null,
});

const endpoints = ref<any[]>([]);
const configTypes = ref<any[]>([]);
const configLimit = ref(0);

watch(
  () => visible.value,
  async (open) => {
    if (!open) return;

    endpoints.value = await getEndpoints();
    configTypes.value = [];
    form.endpoint_id = null;
    form.config_type_id = null;
  },
);

watch(
  () => form.endpoint_id,
  async (endpointId) => {
    form.config_type_id = null;

    if (!endpointId) {
      configTypes.value = [];
      return;
    }

    configTypes.value = await getEndpointConfigTypes(endpointId);
  },
);
const creating = ref(false);
async function submitCreateConfig() {
    if (
        creating.value ||
        !form.name ||
        !form.endpoint_id ||
        !form.config_type_id
    ) {
        return;
    }

    creating.value = true;

    try {
        await createConfiguration({
            name: form.name,
            endpoint_id: form.endpoint_id,
            config_type_id: form.config_type_id,
        });

        form.name = '';
        form.endpoint_id = null;
        form.config_type_id = null;
        configTypes.value = [];

        emit('created');
        visible.value = false;
    } catch (e) {
        console.error(e);
        alert((e as any)?.response?.data ?? 'Ошибка создания конфигурации');
    } finally {
        creating.value = false;
    }
}
</script>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.field {
  display: flex;
  flex-direction: column;
  gap: .35rem;
}
.limit-card {
  margin-top: .5rem;
}
.limit-row {
  display: flex;
  align-items: center;
  gap: .5rem;
}
</style>