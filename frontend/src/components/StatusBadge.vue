<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-semibold rounded-full',
      statusClasses[status] || statusClasses.pending
    ]"
  >
    <span
      v-if="showDot"
      :class="[
        'w-1.5 h-1.5 rounded-full',
        dotClasses[status] || dotClasses.pending
      ]"
    />
    <slot>{{ status }}</slot>
  </span>
</template>

<script setup>
const props = defineProps({
  status: {
    type: String,
    required: true
  },
  showDot: {
    type: Boolean,
    default: true
  }
})

const statusClasses = {
  pending: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300 border border-yellow-200 dark:border-yellow-800',
  running: 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300 border border-blue-200 dark:border-blue-800',
  completed: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 border border-green-200 dark:border-green-800',
  failed: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300 border border-red-200 dark:border-red-800',
  cancelled: 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 border border-gray-200 dark:border-gray-600',
  success: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300 border border-green-200 dark:border-green-800',
  error: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300 border border-red-200 dark:border-red-800'
}

const dotClasses = {
  pending: 'bg-yellow-600',
  running: 'bg-blue-600 animate-pulse',
  completed: 'bg-green-600',
  failed: 'bg-red-600',
  cancelled: 'bg-gray-600',
  success: 'bg-green-600',
  error: 'bg-red-600'
}
</script>

