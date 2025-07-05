<script setup>
import { computed,ref,watch } from 'vue';
import axios from 'axios'

import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import Tabs from 'primevue/tabs';
import TabList from 'primevue/tablist';
import Tab from 'primevue/tab';
import TabPanels from 'primevue/tabpanels';
import TabPanel from 'primevue/tabpanel';
import Chip from 'primevue/chip'

const details = ref(null)
const isLoading = ref(false)
const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  file: {
    type: Object,
    default: null,
  },
  basePath: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:visible']);

const downloadLink = computed(() => {
  if (!props.file) return '#';
  const fullPath = props.basePath.endsWith('/')
    ? `${props.basePath}${props.file.name}`
    : `${props.basePath}/${props.file.name}`;
  return `/api/download?file=${encodeURIComponent(fullPath)}`;
});

const onHide = () => {
  emit('update:visible', false);
};

watch(() => props.file, async (newFile) => {
  if (newFile && props.visible) {
    isLoading.value = true;
    details.value = null;
    try {
      const fullPath = props.basePath.endsWith('/')
        ? `${props.basePath}${newFile.name}`
        : `${props.basePath}/${newFile.name}`;

      const response = await axios.get(`/api/file/details?file=${encodeURIComponent(fullPath)}`);
      details.value = response.data;
    } catch (e) {
      console.error("Failed to fetch file details:", e);
    } finally {
      isLoading.value = false;
    }
  }
  console.log(details)
}, { immediate: true });


</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="onHide"
    modal
    header="File Details"
    :style="{ width: '25rem' }"
  >
    <div class="card">
      <Tabs value="0">
        <TabList>
          <Tab value="0"><i class="pi pi-file"></i></Tab>
          <Tab value="1"><i class="pi pi-info-circle"></i></Tab>
        </TabList>
        <TabPanels>
          <TabPanel value="0">
            <div v-if="props.file" class="file-details-content">
              <div class="file-info">
                <span class="pi pi-file" style="font-size: 2rem"></span>
                <p>{{ props.file.name }}</p>
              </div>
              <!-- 点击下载按钮后，也关闭弹窗 -->
              <a :href="downloadLink" :download="props.file.name" @click="onHide">
                <Button icon="pi pi-download" />
              </a>
            </div>
          </TabPanel>
          <TabPanel value="1">
            <div v-if="isLoading"><i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i></div>
            <div v-else-if="details">
              <Chip :label="`${details.size}bytes`" icon="pi pi-file" />
              <p><strong>Last Modify:</strong> {{ new Date(details.lastModify).toLocaleString() }}</p>
              <p><strong>Permission:</strong> {{ details.permission }}</p>
              <p><strong>Full Path:</strong> {{ details.truePath }}</p>
            </div>
          </TabPanel>

        </TabPanels>
      </Tabs>
    </div>

  </Dialog>
</template>

<style scoped>
.file-details-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
  padding-top: 1rem;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.1rem;
}
</style>
