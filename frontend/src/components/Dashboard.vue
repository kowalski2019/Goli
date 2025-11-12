<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-white mb-6">Dashboard</h2>
    
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="card">
        <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Total Jobs</h3>
        <p class="mt-2 text-3xl font-bold text-gray-900">{{ stats.total }}</p>
      </div>
      <div class="card">
        <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Running</h3>
        <p class="mt-2 text-3xl font-bold text-blue-600">{{ stats.running }}</p>
      </div>
      <div class="card">
        <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Completed</h3>
        <p class="mt-2 text-3xl font-bold text-green-600">{{ stats.completed }}</p>
      </div>
      <div class="card">
        <h3 class="text-sm font-medium text-gray-500 uppercase tracking-wide">Failed</h3>
        <p class="mt-2 text-3xl font-bold text-red-600">{{ stats.failed }}</p>
      </div>
    </div>

    <!-- Recent Jobs -->
    <div class="card">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-xl font-bold text-gray-900">Recent Jobs</h3>
        <button @click="loadJobs()" class="btn btn-secondary text-sm">
          Refresh
        </button>
      </div>
      <div v-if="loading" class="text-center py-8 text-gray-500">
        Loading...
      </div>
      <div v-else-if="recentJobs.length === 0" class="text-center py-8 text-gray-500">
        No jobs found
      </div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Started</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="job in recentJobs" :key="job.id" class="hover:bg-gray-50">
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
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <button 
                  @click="viewJobLogs(job.id)"
                  class="text-primary-600 hover:text-primary-800 mr-3"
                >
                  View Logs
                </button>
                <button 
                  v-if="job.status === 'pending' || job.status === 'running'"
                  @click="cancelJob(job.id)"
                  :disabled="isCancelling === job.id"
                  class="text-red-600 hover:text-red-800 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ isCancelling === job.id ? 'Cancelling...' : 'Cancel' }}
                </button>
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

