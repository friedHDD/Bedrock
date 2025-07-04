<script setup>
import {computed, ref, watch} from 'vue'
import axios from 'axios'

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Chip from 'primevue/chip';
import Breadcrumb from 'primevue/breadcrumb';

const props = defineProps({
  path: {
    type: String,
    required: false,
    default: ''
  }
})

const files = ref([])

const getFolderLink = (fileName) => {
  if (props.path && !props.path.endsWith('/')){
    return `/files${props.path}/${fileName}/`
  }
  return `/files${props.path}${fileName}/`
}

const breadcrumbItems = computed(() => {
  const segments = props.path.split('/').filter(Boolean); // 分割并移除空项
  let currentPath = '';
  return segments.map(segment => {
    currentPath = currentPath ? `${currentPath}/${segment}` : segment;
    return {
      label: segment,
    };
  });
})

watch(//events after the path changed
  () => props.path,
  (newPath) => {
    axios.get(`/api/list?folder=${newPath}`)
      .then(response => {
        files.value = response.data
      })
      .catch(error => {
        console.error('请求文件列表失败:', error)
        files.value = []
      })
  },
  { immediate: true }
)
</script>

<template>
  <div class="card flex justify-center">
    <Breadcrumb :model="breadcrumbItems" />
  </div>

  <DataTable :value="files" tableStyle="min-width: 50rem">
    <Column field="name" header="Name" :sortable="true">
      <template #body="slotProps">
        <router-link :to="getFolderLink(slotProps.data.name)" v-if="slotProps.data.type === 'folder'">
          <Button icon="pi pi-folder" :label="slotProps.data.name" severity="secondary" variant="text" />
        </router-link>
        <Button v-else icon="pi pi-file" :label="slotProps.data.name" variant="text" />
      </template>
    </Column>
    <Column field="lastModify" header="Last Modify" :sortable="true"></Column>
    <Column field="permission" header="Permission">
      <template #body="slotProps">
        <Chip :label="slotProps.data.permission" />
      </template>
    </Column>
  </DataTable>
</template>

<style scoped>

</style>
