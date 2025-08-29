<script setup lang="ts">
defineProps<{
  text?: string
}>()
</script>

<template>
  <div class="spinner-overlay">
    <div class="spinner"></div>
    <p>{{ text || "Loading..." }}</p>
  </div>
</template>

<style scoped>
.spinner-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0,0,0,0.7); /* dark backdrop */
  padding: 20px 30px;
  border-radius: 10px;
  text-align: center;
  z-index: 9999;
  color: #fff;
  font-weight: 500;
  font-size: 14px;
}

.spinner {
  width: 36px;
  height: 36px;
  border: 4px solid rgba(255,255,255,0.3);
  border-top: 4px solid #ffffff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
