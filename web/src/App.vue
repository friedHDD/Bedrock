<script setup>
import {RouterLink, RouterView} from 'vue-router'
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
</script>

<template>
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

</template>


<style>
.global-toolbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  margin-top: 4px;
  z-index: 1000;
  padding: 0.75rem 1.5rem;
}

.app-container {
  padding: 5rem 2rem 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.swal2-container {
  z-index: 9999 !important;
}
</style>
