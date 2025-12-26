<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-3xl max-h-[90vh] overflow-y-auto">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between sticky top-0 bg-white dark:bg-gray-800">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          Job Details - #{{ job.id }}
        </h3>
        <button @click="$emit('close')" class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-6 space-y-6">
        <!-- Job Info -->
        <div>
          <h4 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">Job Information</h4>
          <dl class="grid grid-cols-2 gap-4">
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Name</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ job.name }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Status</dt>
              <dd class="mt-1">
                <span :class="getStatusBadgeClass(job.status)">
                  {{ job.status }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Created At</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDate(job.created_at) }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Started At</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDate(job.started_at) || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Completed At</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDate(job.completed_at) || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Triggered By</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ job.triggered_by || '-' }}</dd>
            </div>
          </dl>
        </div>

        <!-- Steps -->
        <div v-if="job.steps && job.steps.length > 0">
          <h4 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">Pipeline Steps</h4>
          <div class="space-y-3">
            <div
              v-for="(step, index) in sortedSteps"
              :key="step.id"
              class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-gray-50 dark:bg-gray-700/50"
            >
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center space-x-3">
                  <span class="flex items-center justify-center w-8 h-8 rounded-full bg-primary-100 dark:bg-primary-900 text-primary-800 dark:text-primary-300 text-sm font-semibold">
                    {{ step.step_order }}
                  </span>
                  <span class="font-medium text-gray-900 dark:text-white">{{ step.step_name }}</span>
                </div>
                <span :class="getStepStatusClass(step.status)">
                  {{ step.status }}
                </span>
              </div>
              <div class="text-sm text-gray-500 dark:text-gray-400 space-y-1">
                <div v-if="step.started_at">
                  Started: {{ formatDate(step.started_at) }}
                </div>
                <div v-if="step.completed_at">
                  Completed: {{ formatDate(step.completed_at) }}
                </div>
                <div v-if="step.error_message" class="text-red-600 dark:text-red-400 mt-2">
                  <strong>Error:</strong> {{ step.error_message }}
                </div>
              </div>
              <button
                v-if="step.logs"
                @click="viewStepLogs(step)"
                class="mt-2 text-sm text-primary-600 hover:text-primary-800"
              >
                View Logs →
              </button>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="job.error_message" class="p-4 bg-red-50 border border-red-200 rounded-lg">
          <h4 class="text-sm font-semibold text-red-800 mb-1">Error</h4>
          <p class="text-sm text-red-700">{{ job.error_message }}</p>
        </div>
      </div>

      <div class="px-6 py-4 border-t border-gray-200 flex justify-end space-x-3">
        <button 
          v-if="job.status === 'pending' || job.status === 'running'"
          @click="cancelJob"
          :disabled="isCancelling"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isCancelling ? 'Cancelling...' : 'Cancel Job' }}
        </button>
        <button @click="$emit('close')" class="btn btn-secondary">
          Close
        </button>
        <button @click="viewLogs" class="btn btn-primary">
          View All Logs
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { cancelJob as cancelJobAPI, getJob } from '../api/client'

const props = defineProps({
  job: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'view-logs', 'job-updated'])

const isCancelling = ref(false)

const sortedSteps = computed(() => {
  if (!props.job || !props.job.steps || !Array.isArray(props.job.steps)) return []
  return [...props.job.steps].sort((a, b) => (a.step_order || 0) - (b.step_order || 0))
})

function formatDate(dateString) {
  if (!dateString) return null
  return new Date(dateString).toLocaleString()
}

function getStatusBadgeClass(status) {
  const classes = {
    pending: 'px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800',
    running: 'px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800',
    completed: 'px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800',
    failed: 'px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800',
    cancelled: 'px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-800'
  }
  return classes[status] || classes.pending
}

function getStepStatusClass(status) {
  return getStatusBadgeClass(status)
}

function viewStepLogs(step) {
  // Emit event to parent to show logs modal
  console.log('View logs for step:', step)
}

function viewLogs() {
  emit('view-logs', props.job.id)
}

async function cancelJob() {
  if (!confirm('Are you sure you want to cancel this job?')) {
    return
  }
  
  isCancelling.value = true
  try {
    await cancelJobAPI(props.job.id)
    // Reload job details
    const updatedJob = await getJob(props.job.id)
    emit('job-updated', updatedJob)
    // Close modal after a short delay
    setTimeout(() => {
      emit('close')
    }, 1000)
  } catch (error) {
    alert(error.message || 'Failed to cancel job')
  } finally {
    isCancelling.value = false
  }
}
</script>

