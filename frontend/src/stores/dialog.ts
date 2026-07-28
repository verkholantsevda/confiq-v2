import { defineStore } from "pinia";

export const useDialogsStore = defineStore("dialogs", {
  state: () => ({
    userCreate: false,
    endpointCreate: false,
    groupCreate: false,
    configCreate: false,
  }),

  actions: {
    openUserCreate() {
      this.userCreate = true;
    },
    openEndpointCreate() {
      this.endpointCreate = true;
    },
    openGroupCreate() {
      this.groupCreate = true;
    },
    openConfigCreate() {
      this.configCreate = true;
    },
  },
});