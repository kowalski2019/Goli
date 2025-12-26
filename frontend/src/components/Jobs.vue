<template>
  <div class="space-y-6 pb-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold text-white dark:text-gray-100 mb-1">Jobs</h2>
        <p class="text-gray-300 dark:text-gray-400 text-sm">Monitor and manage your CI/CD jobs</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="btn btn-primary px-6 py-2.5 flex items-center gap-2 hover:shadow-lg transition-shadow"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Create Job
      </button>
    </div>

    <!-- Jobs Table -->
    <div class="card hover:shadow-lg transition-shadow duration-200">
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
      <div v-else-if="jobs.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="p-4 bg-gray-100 dark:bg-gray-700 rounded-full mb-4">
          <svg class="w-12 h-12 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">No jobs found</h3>
        <p class="text-gray-500 dark:text-gray-400 text-sm mb-4">Create a new job to get started!</p>
        <button @click="showCreateModal = true" class="btn btn-primary">
          Create Your First Job
        </button>
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
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Completed
              </th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
            <tr
              v-for="job in jobs"
              :key="job.id"
              class="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-150 cursor-pointer"
              @click="selectJob(job)"
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
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm text-gray-500 dark:text-gray-400">
                  {{ formatDate(job.completed_at) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                <div class="flex items-center justify-end gap-2" @click.stop>
                  <button
                    @click.stop="viewJobDetails(job.id)"
                    class="text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1.5"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    Details
                  </button>
                  <button
                    @click.stop="viewJobLogs(job.id)"
                    class="text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300 hover:bg-primary-50 dark:hover:bg-primary-900/30 px-3 py-1.5 rounded-lg transition-colors flex items-center gap-1.5"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    Logs
                  </button>
                  <button
                    v-if="job.status === 'pending' || job.status === 'running'"
                    @click.stop="cancelJob(job.id)"
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

    <!-- Job Details Modal -->
    <JobDetailsModal
      v-if="selectedJob"
      :job="selectedJob"
      @close="selectedJob = null"
      @view-logs="handleViewLogs"
    />

    <!-- Job Logs Modal -->
    <JobLogsModal
      v-if="logsJobId"
      :job-id="logsJobId"
      @close="logsJobId = null"
    />

    <!-- Create Job Modal -->
    <CreateJobModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @created="handleJobCreated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getJobs, cancelJob as cancelJobAPI } from '../api/client'
import JobDetailsModal from './JobDetailsModal.vue'
import JobLogsModal from './JobLogsModal.vue'
import CreateJobModal from './CreateJobModal.vue'
import StatusBadge from './StatusBadge.vue'

const emit = defineEmits(['view-logs'])

const loading = ref(false)
const jobs = ref([])
const selectedJob = ref(null)
const logsJobId = ref(null)
const showCreateModal = ref(false)
const isCancelling = ref(null)

function formatDate(dateString) {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function loadJobs() {
  loading.value = true
  try {
    const jobsData = await getJobs({ limit: 50 })
    jobs.value = Array.isArray(jobsData) ? jobsData : []
  } catch (error) {
    console.error('Error loading jobs:', error)
    jobs.value = []
  } finally {
    loading.value = false
  }
}

function selectJob(job) {
  selectedJob.value = job
}

function viewJobDetails(jobId) {
  const job = jobs.value.find(j => j.id === jobId)
  if (job) {
    selectedJob.value = job
  }
}

function viewJobLogs(jobId) {
  logsJobId.value = jobId
}

function handleViewLogs(jobId) {
  emit('view-logs', jobId)
  selectedJob.value = null
}

function handleJobCreated() {
  showCreateModal.value = false
  loadJobs()
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
