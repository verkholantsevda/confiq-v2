<template>
    <Card>

        <template #content>

            <DataTable
                :value="value"
                :loading="loading"
                paginator
                :rows="10"
                :rowsPerPageOptions="[10, 20, 50, 100]"
                stripedRows
                removableSort
                responsiveLayout="scroll"
                dataKey="id"
                filterDisplay="row"
            >

                <template #header>

                    <div class="table-header">

                        <IconField>

                            <InputIcon class="pi pi-search" />

                            <InputText
                                v-model="search"
                                :placeholder="searchPlaceholder"
                            />

                        </IconField>

                        <slot name="header" />

                    </div>

                </template>

                <slot />

                <template #empty>

                    {{ emptyMessage }}

                </template>

            </DataTable>

        </template>

    </Card>
</template>

<script setup lang="ts">
import { ref } from "vue";

import Card from "primevue/card";
import DataTable from "primevue/datatable";
import InputText from "primevue/inputtext";
import IconField from "primevue/iconfield";
import InputIcon from "primevue/inputicon";

defineProps<{
    value: any[];
    loading?: boolean;

    emptyMessage?: string;
    searchPlaceholder?: string;
}>();

const search = ref("");
</script>

<style scoped>
.table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: .5rem;
}
</style>