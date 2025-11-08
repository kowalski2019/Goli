<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold text-white">Jobs</h2>
      <button @click="showCreateModal = true" class="btn btn-primary">
        + Create Job
      </button>
    </div>

    <!-- Jobs Table -->
    <div class="card">
      <div v-if="loading" class="text-center py-8 text-gray-500">
        Loading jobs...
      </div>
      <div v-else-if="jobs.length === 0" class="text-center py-8 text-gray-500">
        No jobs found
      </div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Started</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Completed</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr 
              v-for="job in jobs" 
              :key="job.id" 
              class="hover:bg-gray-50 cursor-pointer"
              @click="selectJob(job)"
            >
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">#{{ job.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ job.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getStatusBadgeClass(job.status)">
                  {{ job.status }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(job.started_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(job.completed_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <button 
                  @click.stop="viewJobDetails(job.id)"
                  class="text-primary-600 hover:text-primary-800 mr-3"
                >
                  Details
                </button>
                <button 
                  @click.stop="viewJobLogs(job.id)"
                  class="text-primary-600 hover:text-primary-800"
                >
                  Logs
                </button>
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
import { getJobs } from '../api/client'
import JobDetailsModal from './JobDetailsModal.vue'
import JobLogsModal from './JobLogsModal.vue'
import CreateJobModal from './CreateJobModal.vue'

const loading = ref(false)
const jobs = ref([])
const selectedJob = ref(null)
const logsJobId = ref(null)
const showCreateModal = ref(false)

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

function handleJobCreated() {
  showCreateModal.value = false
  loadJobs()
}

onMounted(() => {
  loadJobs()
  // Refresh every 5 seconds
  setInterval(loadJobs, 5000)
})
</script>

