<template>
    <PageHeader
        :title="t('pages.configtypes.create')"
    />

    <Card>
        <template #content>

            <div class="form">

                <div class="field">
                    <label>{{ t("pages.configtypes.name") }}</label>

                    <InputText
                        v-model="configType.name"
                        fluid
                    />
                </div>

                <div class="field">
                    <label>{{ t("pages.configtypes.description_type") }}</label>

                    <Textarea
                        v-model="configType.description"
                        rows="3"
                        autoResize
                        fluid
                    />
                </div>

                <div class="field">
                    <label>{{ t("pages.configtypes.active") }}</label>

                    <ToggleSwitch
                        v-model="configType.is_active"
                    />
                </div>

                <Divider />

                <div class="field">
                    <label>{{ t("pages.configtypes.template") }}</label>

                    <Textarea
                        v-model="configType.config_template"
                        rows="14"
                        class="monospace"
                        autoResize
                        fluid
                    />
                </div>

                <Divider />

                <div class="field">
                    <label>{{ t("pages.configtypes.instruction") }}</label>

                    <Textarea
                        v-model="configType.usage_instructions"
                        rows="12"
                        autoResize
                        fluid
                    />
                </div>

                <Divider />

                <div class="field">
                    <label>{{ t("pages.configtypes.links_clients") }} (JSON)</label>

                    <Textarea
                        v-model="configType.client_links"
                        rows="10"
                        class="monospace"
                        autoResize
                        fluid
                    />
                </div>

                <Divider />

                <div class="actions">

                    <Button
                        icon="pi pi-check"
                        :label="t('pages.configtypes.create')"
                        :loading="saving"
                        @click="save"
                    />

                    <Button
                        icon="pi pi-times"
                        :label="t('common.cancel')"
                        severity="secondary"
                        outlined
                        @click="router.back()"
                    />

                </div>

            </div>

        </template>
    </Card>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";

import Card from "primevue/card";
import Button from "primevue/button";
import Divider from "primevue/divider";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import ToggleSwitch from "primevue/toggleswitch";

import PageHeader from "@/components/common/PageHeader.vue";

import { createConfigType } from "@/api/configTypes";

const { t } = useI18n();

const router = useRouter();

const saving = ref(false);

const configType = reactive({
    id: 0,
    name: "",
    description: "",
    config_template: "",
    usage_instructions: "",
    client_links: "",
    is_active: false,
});

async function save() {
    saving.value = true;

    try {
        await createConfigType(configType);

        router.push("/admin/config-types");
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

.monospace :deep(textarea) {
    font-family: "JetBrains Mono", monospace;
    font-size: .9rem;
}

.actions {
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
}
</style>