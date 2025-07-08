<script setup>
import { computed } from 'vue';

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
});

const iframeSrc = computed(() => {
  const bookApiUrl = `/api/library/get/${props.id}`;
  return `/foliate-js/reader.html?url=${encodeURIComponent(bookApiUrl)}`;
});
</script>

<template>
  <div class="reader-page">
    <iframe
      :src="iframeSrc"
      class="reader-iframe"
      allowfullscreen
    ></iframe>
  </div>
</template>

<style scoped>
.reader-page {
  width: 100%;
  height: 87vh;
  overflow: hidden;
}

.reader-iframe {
  width: 100%;
  height: 100%;
  border: none;
}
</style>
