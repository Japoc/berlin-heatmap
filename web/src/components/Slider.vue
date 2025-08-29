<script setup lang="ts">
import { ref, watch } from "vue"

// props
const props = defineProps<{
  modelValue: number
  label?: string
  min?: number
  max?: number
  steps?: number
}>()

// emits
const emit = defineEmits<{
  (e: "update:modelValue", value: number): void
}>()

// local state
const localValue = ref(props.modelValue)

// keep local in sync with prop
watch(() => props.modelValue, (newVal) => {
  if (newVal !== localValue.value) {
    localValue.value = newVal
  }
})

// update parent when slider changes
function onInput(e: Event) {
  const val = Number((e.target as HTMLInputElement).value)
  localValue.value = val
  emit("update:modelValue", val)
}
</script>

<template>
  <div class="slider-container">
    <label>
      {{ label }}: <strong>{{ localValue }}</strong>
    </label>
    <input
        type="range"
        :min="min ?? 10"
        :max="max ?? 300"
        :step="steps ?? 10"
        v-model="localValue"
        @input="onInput"
    />
  </div>
</template>

<style scoped>
.slider-container {
  text-align: center;
  width: 100%;
}
input[type="range"] {
  width: 100%;
}
</style>
