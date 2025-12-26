<template>
  <label :class="['relative inline-flex items-center cursor-pointer', containerClass]">
    <input
      type="checkbox"
      :checked="modelValue"
      :disabled="disabled"
      class="sr-only peer"
      @change="$emit('update:modelValue', $event.target.checked)"
    />
    <div
      :class="[
        'w-11 h-6 rounded-full transition-colors duration-200 ease-in-out',
        'peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-offset-2',
        disabled
          ? 'bg-gray-200 dark:bg-gray-700 cursor-not-allowed'
          : modelValue
          ? 'bg-primary-600 dark:bg-primary-500 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-700'
          : 'bg-gray-300 dark:bg-gray-600 peer-focus:ring-gray-300 dark:peer-focus:ring-gray-700',
        error ? 'ring-2 ring-red-500 dark:ring-red-400' : ''
      ]"
    >
      <div
        :class="[
          'absolute top-[2px] left-[2px] bg-white dark:bg-gray-200 rounded-full transition-transform duration-200 ease-in-out',
          'h-5 w-5 border border-gray-300 dark:border-gray-600',
          modelValue ? 'translate-x-5 border-primary-600 dark:border-primary-400' : 'translate-x-0'
        ]"
      />
    </div>
    <span v-if="label" :class="['ml-3 text-sm font-medium', labelClass || 'text-gray-700 dark:text-gray-300']">
      {{ label }}
    </span>
  </label>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  },
  containerClass: {
    type: String,
    default: ''
  },
  labelClass: {
    type: String,
    default: 'text-gray-700'
  }
})

defineEmits(['update:modelValue'])
</script>

