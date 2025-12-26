<template>
  <div :class="['space-y-1.5', containerClass]">
    <label 
      v-if="label" 
      :for="id"
      :class="[
        'block text-sm font-medium',
        required ? 'text-gray-900 dark:text-white' : 'text-gray-700 dark:text-gray-300',
        labelClass
      ]"
    >
      {{ label }}
      <span v-if="required" class="text-red-500 dark:text-red-400 ml-0.5">*</span>
      <span v-if="hint" class="text-gray-500 dark:text-gray-400 font-normal ml-1">({{ hint }})</span>
    </label>
    <div class="relative">
      <slot />
      <div v-if="description" class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
        {{ description }}
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  hint: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  containerClass: {
    type: String,
    default: ''
  },
  labelClass: {
    type: String,
    default: ''
  }
})
</script>

