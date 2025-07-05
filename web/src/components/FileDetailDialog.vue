<script setup>
import { computed,ref,watch } from 'vue';
import axios from 'axios'

import {Dialog,
  Button,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  Chip,
  Checkbox,
  Message,
  ButtonGroup
} from "primevue";

const details = ref(null)
const isLoading = ref(false)
const isConfirmedIPFS = ref(false);
const isAddingToIpfs = ref(false);
const ipfsCid = ref('');
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

async function addToIPFS() {
  isAddingToIpfs.value = true
  const fullPath = props.basePath.endsWith('/')
    ? `${props.basePath}${props.file.name}`
    : `${props.basePath}/${props.file.name}`;
  try {
    const response = await axios.get(`/api/ipfs/add?file=${encodeURIComponent(fullPath)}`);
    ipfsCid.value = response.data.cid;
  } catch (error) {
    ipfsCid.value = error.response?.data?.error || 'An unknown error occurred.';
  } finally {
    isAddingToIpfs.value = false;
  }
}

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
  ipfsCid.value = null
  console.log(details)
}, { immediate: true });


</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="onHide"
    modal
    header="File Details"
    :style="{ width: '30rem' }"
  >
    <div class="card">
      <Tabs value="0">
        <TabList>
          <Tab value="0"><i class="pi pi-file"></i></Tab>
          <Tab value="1"><i class="pi pi-info-circle"></i></Tab>
          <Tab value="2"><i class="pi pi-box"></i></Tab>
        </TabList>
        <TabPanels>
          <TabPanel value="0">
            <div v-if="props.file" class="file-details-content">
              <div class="file-info">
                <span class="pi pi-file" style="font-size: 2rem"></span>
                <p>{{ props.file.name }}</p>
              </div>
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
          <TabPanel value="2">
            <div v-if="props.file" class="file-details-content">
              <Message severity="warn">This will make the file public</Message>
              <div class="confirmation-section">
                <Checkbox v-model="isConfirmedIPFS" inputId="ipfs-confirm" binary />
                <label for="ipfs-confirm" class="ml-2"> I understand.</label>
              </div>

              <Button
                label="Add to IPFS"
                icon="pi pi-upload"
                :disabled="!isConfirmedIPFS||ipfsCid"
                :loading="isAddingToIpfs"
                @click="addToIPFS"
              />
              <div v-if="ipfsCid" class="file-details-content">
                <Chip :label="`cid:${ipfsCid.slice(0, 6)}...${ipfsCid.slice(-6)}`" icon="pi pi-file" />
                <ButtonGroup>
                  <Button icon="pi pi-copy" v-clipboard:copy="ipfsCid"/>
                  <a :href="`http://localhost:8080/ipfs/${ipfsCid}`" target="_blank">
                    <Button label="localhost" icon="pi pi-arrow-up-right" variant="text"/>
                  </a>
                  <a :href="`https://ipfs.io/ipfs/${ipfsCid}`" target="_blank">
                    <Button label="ipfs.io" icon="pi pi-arrow-up-right" variant="text"/>
                  </a>
                </ButtonGroup>

              </div>


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
