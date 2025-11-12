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
  pending: 'bg-yellow-100 text-yellow-800 border border-yellow-200',
  running: 'bg-blue-100 text-blue-800 border border-blue-200',
  completed: 'bg-green-100 text-green-800 border border-green-200',
  failed: 'bg-red-100 text-red-800 border border-red-200',
  cancelled: 'bg-gray-100 text-gray-800 border border-gray-200',
  success: 'bg-green-100 text-green-800 border border-green-200',
  error: 'bg-red-100 text-red-800 border border-red-200'
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

