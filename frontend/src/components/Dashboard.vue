<template>
  <div class="space-y-6 pb-8">
    <!-- Header -->
    <div>
      <h2 class="text-3xl font-bold text-white dark:text-gray-100 mb-1">Dashboard</h2>
      <p class="text-gray-300 dark:text-gray-400 text-sm">Overview of your CI/CD pipeline activity</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="card hover:shadow-lg transition-shadow duration-200 border-l-4 border-primary-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Total Jobs</p>
            <p class="text-3xl font-bold text-gray-900 dark:text-white">{{ stats.total }}</p>
          </div>
          <div class="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
            <svg class="w-8 h-8 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card hover:shadow-lg transition-shadow duration-200 border-l-4 border-blue-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Running</p>
            <p class="text-3xl font-bold text-blue-600 dark:text-blue-400">{{ stats.running }}</p>
          </div>
          <div class="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
            <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card hover:shadow-lg transition-shadow duration-200 border-l-4 border-green-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Completed</p>
            <p class="text-3xl font-bold text-green-600 dark:text-green-400">{{ stats.completed }}</p>
          </div>
          <div class="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
            <svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
      </div>

      <div class="card hover:shadow-lg transition-shadow duration-200 border-l-4 border-red-500">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-1">Failed</p>
            <p class="text-3xl font-bold text-red-600 dark:text-red-400">{{ stats.failed }}</p>
          </div>
          <div class="p-3 bg-red-100 dark:bg-red-900 rounded-lg">
            <svg class="w-8 h-8 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Jobs -->
    <div class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white">Recent Jobs</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">Latest job executions</p>
        </div>
        <button
          @click="loadJobs()"
          :disabled="loading"
          class="btn btn-secondary text-sm px-4 py-2 flex items-center gap-2 disabled:opacity-50"
        >
          <svg
            :class="['w-4 h-4', loading && 'animate-spin']"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Refresh
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-12">
        <svg
          class="animate-spin h-8 w-8 text-primary-600 mb-4"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-500 dark:text-gray-400 text-sm">Loading jobs...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="recentJobs.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="p-4 bg-gray-100 dark:bg-gray-700 rounded-full mb-4">
          <svg class="w-12 h-12 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">No jobs found</h3>
        <p class="text-gray-500 dark:text-gray-400 text-sm">Jobs will appear here once you start running pipelines</p>
      </div>

      <!-- Jobs Table -->
      <div v-else class="overflow-x-auto scrollbar-thin">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">
                ID
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Started
              </th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            <tr
              v-for="job in recentJobs"
              :key="job.id"
              class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-150"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm font-medium text-gray-900 dark:text-gray-100">#{{ job.id }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <svg class="w-5 h-5 text-primary-600 dark:text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  <span class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ job.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <StatusBadge :status="job.status" />
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm text-gray-500 dark:text-gray-400">
                  {{ formatDate(job.started_at) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                <div class="flex items-center justify-end gap-2">
                  <button
                    @click="viewJobLogs(job.id)"
                    class="text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1.5"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    View Logs
                  </button>
                  <button
                    v-if="job.status === 'pending' || job.status === 'running'"
                    @click="cancelJob(job.id)"
                    :disabled="isCancelling === job.id"
                    class="text-red-600 dark:text-red-400 hover:text-red-800 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/30 px-3 py-1.5 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
                  >
                    <svg
                      v-if="isCancelling === job.id"
                      class="animate-spin h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <svg
                      v-else
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    {{ isCancelling === job.id ? 'Cancelling...' : 'Cancel' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getJobs, cancelJob as cancelJobAPI } from '../api/client'
import StatusBadge from './StatusBadge.vue'

const emit = defineEmits(['view-logs'])

const loading = ref(false)
const recentJobs = ref([])
const isCancelling = ref(null)
const stats = ref({
  total: 0,
  running: 0,
  completed: 0,
  failed: 0
})

function formatDate(dateString) {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function loadJobs() {
  loading.value = true
  try {
    const jobs = await getJobs({ limit: 10 })
    recentJobs.value = Array.isArray(jobs) ? jobs : []
    updateStats(recentJobs.value)
  } catch (error) {
    console.error('Error loading jobs:', error)
    recentJobs.value = []
    updateStats([])
  } finally {
    loading.value = false
  }
}

function updateStats(jobs) {
  const jobsArray = Array.isArray(jobs) ? jobs : []
  stats.value = {
    total: jobsArray.length,
    running: jobsArray.filter(j => j && j.status === 'running').length,
    completed: jobsArray.filter(j => j && j.status === 'completed').length,
    failed: jobsArray.filter(j => j && j.status === 'failed').length
  }
}

function viewJobLogs(jobId) {
  emit('view-logs', jobId)
}

async function cancelJob(jobId) {
  if (!confirm('Are you sure you want to cancel this job?')) {
    return
  }

  isCancelling.value = jobId
  try {
    await cancelJobAPI(jobId)
    await loadJobs()
  } catch (error) {
    alert(error.message || 'Failed to cancel job')
  } finally {
    isCancelling.value = null
  }
}

onMounted(() => {
  loadJobs()
  // Refresh every 5 seconds
  setInterval(loadJobs, 5000)
})
</script>
