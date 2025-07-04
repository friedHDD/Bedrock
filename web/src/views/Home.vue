<script setup>
import Button from "primevue/button"
import {RouterLink} from "vue-router";
</script>

<template>
  <div class="home-container">
    <div class="welcome-content">
      <h1 class="welcome-title">Welcome to Bedrock</h1>
      <p class="welcome-subtitle">A privacy-first file manager</p>
      <RouterLink to="/files">
        <Button label="Go to Files" icon="pi pi-arrow-right" icon-pos="right" />
      </RouterLink>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  height: calc(100vh - 120px);
  min-height: 400px;
}

.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.welcome-title {
  font-size: 3.5rem;
  font-weight: 700;
  letter-spacing: -0.05em;
}

.welcome-subtitle {
  font-size: 1.25rem;
}
</style>
