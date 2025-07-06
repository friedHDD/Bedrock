<script setup>
import {RouterLink, RouterView} from 'vue-router'
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
</script>

<template>
  <div class="app-container">
    <Toolbar class="global-toolbar" style="border-radius: 3rem; padding: 1rem 1rem 1rem 1.5rem">
      <template #start>
        <div class="nav-links">
          <RouterLink to="/">
            <Button icon="pi pi-home" class="mr-2" severity="secondary" text/>
          </RouterLink>
          <RouterLink to="/files/~/">
            <Button icon="pi pi-folder" class="mr-2" severity="secondary" text/>
          </RouterLink>
          <RouterLink to="/library">
            <Button icon="pi pi-book" class="mr-2" severity="secondary" text/>
          </RouterLink>
        </div>
      </template>
      <template #end>
        <div class="nav-links">

        </div>
      </template>
    </Toolbar>
    <div class="app-container">
      <RouterView/>
    </div>
  </div>

</template>


<style>
.app-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.swal2-container {
  z-index: 9999 !important;
}
</style>
