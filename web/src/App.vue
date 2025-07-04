<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
</script>

<template>
  <div class="app-container">
    <Toolbar class="global-toolbar" style="border-radius: 3rem; padding: 1rem 1rem 1rem 1.5rem">
      <template #start>
        <div class="nav-links">
          <RouterLink to="/">
            <Button icon="pi pi-home" class="mr-2" severity="secondary" text />
          </RouterLink>
          <RouterLink to="/files/~/">
            <Button icon="pi pi-folder" class="mr-2" severity="secondary" text />
          </RouterLink>
        </div>
      </template>
    </Toolbar>
    <RouterView />
  </div>

</template>

