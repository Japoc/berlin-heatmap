<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

defineProps<{
  modelValue: boolean
  label: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

function toggle(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.checked)
}
</script>

<template>
  <label class="checkbox-wrapper">
    <input
        type="checkbox"
        :checked="modelValue"
        :disabled="disabled"
        @change="toggle"
    />
    <span>{{ label }}</span>
  </label>
</template>

<style scoped>
.checkbox-wrapper {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 14px;
  cursor: pointer;
  user-select: none;
}

.checkbox-wrapper input[type="checkbox"] {
  cursor: pointer;
}
.checkbox-wrapper input[type="checkbox"]:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}
</style>
