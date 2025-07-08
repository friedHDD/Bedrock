<script setup>
import {ref, onMounted} from 'vue';
import axios from 'axios';
import Swal from 'sweetalert2';
import router from "@/router.js";

import {
  DataTable,
  Column,
  Button,
  ButtonGroup,
  Tag,
} from 'primevue';


const books = ref([]);
const isRefreshing = ref(true);
const isScanning = ref(false);
const isCleaning = ref(false);

const refresh = async () => {
  isRefreshing.value = true;
  try {
    const response = await axios.get('/api/library/list');
    books.value = response.data.books || [];
  } catch (error) {
    const msg = error.response?.data?.message || 'An unknown error occurred.'
    await Swal.fire({
      title: 'Error',
      text: msg,
      icon: 'error',
      confirmButtonText: 'OK'
    });
  } finally {
    isRefreshing.value = false;
  }
};

const scan = async () => {
  isScanning.value = true;
  try {
    const response = await axios.get('/api/library/scan');
    const msg = response.data.message || [];
    await Swal.fire({
      title: 'Scan completed',
      text: msg,
      icon: 'success',
      confirmButtonText: 'Great!'
    });
  } catch (error) {
    const msg = error.response?.data?.message || 'An unknown error occurred.'
    await Swal.fire({
      title: 'Error',
      text: msg,
      icon: 'error',
      confirmButtonText: 'OK'
    });
  } finally {
    isScanning.value = false;
    await refresh()
  }
};

const clean = async () => {
  isCleaning.value = true;
  try {
    const response = await axios.get('/api/library/clean');
    const msg = response.data.message || [];
    await Swal.fire({
      title: 'Accepted',
      text: msg,
      icon: 'success',
      confirmButtonText: 'Great!'
    });
  } catch (error) {
    const msg = error.response?.data?.message || 'An unknown error occurred.'
    await Swal.fire({
      title: 'Error',
      text: msg,
      icon: 'error',
      confirmButtonText: 'OK'
    });
  } finally {
    isCleaning.value = false;
    await refresh()
  }

};

onMounted(() => {
  refresh();
});

const readBook = (book) => {
  router.push(`/library/${book.id}`)
  //window.open(`/library/${book.id}`, '_blank');
};

</script>

<template>
  <header class="library-header">
    <h1>Library</h1>
    <ButtonGroup>
      <Button icon="pi pi-refresh" label="Refresh" @click="refresh" :loading="isRefreshing"/>
      <Button icon="pi pi-wrench" label="Scan" @click="scan" :loading="isScanning"/>
      <Button icon="pi pi-eraser" label="Clean" @click="clean" :loading="isCleaning"/>
    </ButtonGroup>
  </header>

  <div class="card">
    <DataTable
      :value="books"
      :loading="isRefreshing"
      paginator
      :rows="10"
      :rowsPerPageOptions="[5, 10, 50, 100, 200]"
      tableStyle="min-width: 50rem"
      dataKey="id"
    >
      <template #empty>
        <div class="empty-state">
          <i class="pi pi-book" style="font-size: 2rem"></i>
          <p>Your library is empty. Add some books to get started!</p>
        </div>
      </template>

      <Column field="series" header="Series" sortable/>

      <Column field="bookName" header="Title" sortable>
        <template #body="slotProps">
          <div class="book-title">
            <i class="pi pi-book"></i>
            <span>{{ slotProps.data.bookName }}</span>
          </div>
        </template>
      </Column>

      <Column field="type" header="Type" sortable>
        <template #body="slotProps">
          <div class="book-title">
            <Tag severity="secondary" :value="slotProps.data.bookName.split('.').pop()"></Tag>
          </div>
        </template>
      </Column>

      <Column field="id" header="id">
        <template #body="slotProps">
          <Tag severity="secondary" :value="slotProps.data.id.slice(0,9)"></Tag>
        </template>
      </Column>

      <Column header="Actions" style="width: 10rem">
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            label="Read"
            @click="readBook(slotProps.data)"
            size="small"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<style scoped>
.library-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.library-header h1 {
  font-size: 2rem;
  font-weight: 600;
}

.book-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-weight: 500;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
}
</style>
