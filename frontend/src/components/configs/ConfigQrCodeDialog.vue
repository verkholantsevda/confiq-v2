<template>
    <Dialog
        v-model:visible="visible"
        modal
        header="QR Code"
        :style="{width:'350px'}"
    >

        <div class="qr-container">

            <ProgressSpinner v-if="loading"/>

            <img
                v-if="qrCode"
                :src="qrCode"
                class="qr-image"
            />

        </div>

    </Dialog>
</template>


<script setup lang="ts">
import { ref, watch } from "vue";
import QRCode from "qrcode";

import Dialog from "primevue/dialog";


const props = defineProps<{
    modelValue: boolean;
    config: any | null;
}>();


const emit = defineEmits([
    "update:modelValue"
]);


const visible = ref(props.modelValue);

const qrCode = ref("");
const loading = ref(false);



watch(
    () => props.modelValue,
    async (value) => {

        visible.value = value;

        if (value && props.config) {
            await generateQr(props.config);
        }
    }
);


watch(
    visible,
    value => {
        emit("update:modelValue", value);
    }
);



async function generateQr(config:any) {

    if (!config.config_content) {
        qrCode.value = "";
        return;
    }


    loading.value = true;


    qrCode.value = await QRCode.toDataURL(
        config.config_content,
        {
            width:260,
            margin:1,
        }
    );


    loading.value = false;
}

</script>


<style scoped>

.qr-container {
    display:flex;
    justify-content:center;
    align-items:center;
    min-height:250px;
}

.qr-container img {
    width:250px;
    height:250px;
}

</style>