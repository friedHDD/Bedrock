<script setup>
import {ref, onMounted} from 'vue';
import axios from 'axios';
import Swal from 'sweetalert2';

import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Tag from 'primevue/tag';

const books = ref([]);
const isLoading = ref(true);

const fetchBooks = async () => {
  isLoading.value = true;
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
    isLoading.value = false;
  }
};

onMounted(() => {
  fetchBooks();
});

const readBook = (book) => {
  console.log('Reading book:', book.id);
  Swal.fire('Coming Soon!', `${book.id}`, 'info');
};
</script>

<template>
  <header class="library-header">
    <h1>Library</h1>
    <Button icon="pi pi-refresh" label="Refresh" @click="fetchBooks" :loading="isLoading"/>
  </header>

  <div class="card">
    <DataTable
      :value="books"
      :loading="isLoading"
      paginator
      :rows="10"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      tableStyle="min-width: 50rem"
      dataKey="id"
    >
      <template #empty>
        <div class="empty-state">
          <i class="pi pi-book" style="font-size: 2rem"></i>
          <p>Your library is empty. Add some books to get started!</p>
        </div>
      </template>

      <Column field="bookName" header="Title" sortable>
        <template #body="slotProps">
          <div class="book-title">
            <i class="pi pi-book"></i>
            <span>{{ slotProps.data.bookName }}</span>
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
  color: var(--p-text-muted-color);
}
</style>
