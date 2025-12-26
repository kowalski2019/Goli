<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] overflow-y-auto">
      <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between sticky top-0 bg-white dark:bg-gray-800">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          Pipeline Details - {{ pipeline.name }}
        </h3>
        <button @click="$emit('close')" class="text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-6 space-y-6">
        <!-- Pipeline Info -->
        <div>
          <h4 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">Pipeline Information</h4>
          <dl class="grid grid-cols-2 gap-4">
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">ID</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">#{{ pipeline.id }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Name</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ pipeline.name }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Description</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ pipeline.description || '-' }}</dd>
            </div>
            <div>
              <dt class="text-sm text-gray-500 dark:text-gray-400">Created At</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-100">{{ formatDate(pipeline.created_at) }}</dd>
            </div>
          </dl>
        </div>

        <!-- Pipeline Definition -->
        <div>
          <h4 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">Pipeline Definition</h4>
          <div class="bg-gray-900 dark:bg-gray-950 rounded-lg p-4 overflow-x-auto">
            <pre class="text-sm text-green-400 dark:text-green-300 font-mono whitespace-pre-wrap">{{ pipeline.definition }}</pre>
          </div>
        </div>

        <!-- Recent Jobs -->
        <div v-if="recentJobs.length > 0">
          <h4 class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-2">Recent Jobs</h4>
          <div class="space-y-2">
            <div
              v-for="job in recentJobs"
              :key="job.id"
              class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 bg-gray-50 dark:bg-gray-700/30"
            >
              <div class="flex items-center justify-between">
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">#{{ job.id }} - {{ job.name }}</div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">
                    Started: {{ formatDate(job.started_at) }}
                  </div>
                </div>
                <div class="flex items-center space-x-2">
                  <span :class="getStatusBadgeClass(job.status)">
                    {{ job.status }}
                  </span>
                  <button
                    @click="viewJobLogs(job.id)"
                    class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300"
                  >
                    View Logs
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end space-x-3">
        <button @click="$emit('close')" class="btn btn-secondary">
          Close
        </button>
        <button @click="runPipeline" class="btn btn-primary">
          Run Pipeline
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getJobs } from '../api/client'

const props = defineProps({
  pipeline: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'view-logs', 'run-pipeline'])

const recentJobs = ref([])

function formatDate(dateString) {
  if (!dateString) return '-'
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

function viewJobLogs(jobId) {
  emit('view-logs', jobId)
}

function runPipeline() {
  emit('run-pipeline', props.pipeline.id)
}

async function loadRecentJobs() {
  try {
    const jobs = await getJobs({ limit: 10 })
    // Filter jobs for this pipeline
    recentJobs.value = (Array.isArray(jobs) ? jobs : []).filter(job => 
      job.pipeline_id === props.pipeline.id
    ).slice(0, 5)
  } catch (error) {
    console.error('Error loading recent jobs:', error)
    recentJobs.value = []
  }
}

onMounted(() => {
  loadRecentJobs()
})
</script>

