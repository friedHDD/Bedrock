<script setup>
import { ref, watch } from 'vue'
import axios from 'axios'
import router from "@/router.js";

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';

const props = defineProps({
  path: {
    type: String,
    required: false,
    default: ''
  }
})

const files = ref([])
const getFolderLink = (fileName) => {
  if (props.path && !props.path.endsWith('/')) {
    return `/files/${props.path}/${fileName}/`
  }
  // 否则直接拼接 (处理 props.path 为空或'~/'的情况)
  return `/files/${props.path}${fileName}/`
}
watch(
  () => props.path,
  (newPath) => {
    // 如果路径为空，则使用 '~'
    if(newPath === '') {
      router.replace('/files/~/')//automatic go to ~/
      return
    }



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

  {{props.path}}

  <DataTable :value="files" tableStyle="min-width: 50rem">

    <Column field="name" header="Name" :sortable="true">
      <template #body="slotProps">
        <!-- 文件夹 -->
        <router-link :to="getFolderLink(slotProps.data.name)" v-if="slotProps.data.type === 'folder'">
          {{ slotProps.data.name }}
        </router-link>
        <!-- 文件 -->
        <span v-else>{{ slotProps.data.name }}</span>
      </template>
    </Column>
    <Column field="type" header="Type"></Column>
    <Column field="lastModify" header="Last Modify" :sortable="true"></Column>
    <Column field="permission" header="Permission"></Column>
  </DataTable>
</template>

<style scoped>

</style>
