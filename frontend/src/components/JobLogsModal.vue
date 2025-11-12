<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col">
      <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">
          Job Logs - #{{ jobId }}
        </h3>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="flex-1 overflow-hidden flex">
        <!-- Steps Sidebar -->
        <div class="w-64 border-r border-gray-200 overflow-y-auto bg-gray-50">
          <div class="p-4">
            <h4 class="text-sm font-semibold text-gray-700 mb-3">Steps</h4>
            <div v-if="loadingSteps" class="text-sm text-gray-500">Loading...</div>
            <div v-else-if="steps.length === 0" class="text-sm text-gray-500">No steps found</div>
            <div v-else class="space-y-2">
              <button
                v-for="step in steps"
                :key="step.id"
                @click="selectedStep = step"
                :class="[
                  'w-full text-left px-3 py-2 rounded-lg text-sm transition-colors',
                  selectedStep?.id === step.id 
                    ? 'bg-primary-100 text-primary-800 border border-primary-300' 
                    : 'bg-white hover:bg-gray-100 border border-gray-200'
                ]"
              >
                <div class="flex items-center justify-between mb-1">
                  <span class="font-medium">{{ step.step_name }}</span>
                  <span :class="getStepStatusClass(step.status)">
                    {{ step.status }}
                  </span>
                </div>
                <div class="text-xs text-gray-500">
                  Step {{ step.step_order }}
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Logs Content -->
        <div class="flex-1 flex flex-col overflow-hidden">
          <div v-if="selectedStep" class="p-4 border-b border-gray-200 bg-gray-50">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-semibold text-gray-900">{{ selectedStep.step_name }}</h4>
                <p class="text-sm text-gray-500">Step {{ selectedStep.step_order }}</p>
              </div>
              <span :class="getStepStatusClass(selectedStep.status)">
                {{ selectedStep.status }}
              </span>
            </div>
            <div v-if="selectedStep.error_message" class="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-800">
              <strong>Error:</strong> {{ selectedStep.error_message }}
            </div>
          </div>

          <div class="flex-1 overflow-y-auto p-4 bg-gray-900">
            <div v-if="loadingLogs" class="text-gray-400">Loading logs...</div>
            <div v-else-if="selectedStep && selectedStep.logs" class="font-mono text-sm text-green-400 whitespace-pre-wrap">
              {{ selectedStep.logs }}
            </div>
            <div v-else-if="selectedStep && !selectedStep.logs" class="text-gray-500">
              No logs available for this step
            </div>
            <div v-else class="text-gray-500">
              Select a step to view logs
            </div>
          </div>

          <!-- Job-level logs -->
          <div v-if="jobLogs" class="border-t border-gray-200 p-4 bg-gray-50">
            <h4 class="text-sm font-semibold text-gray-700 mb-2">Job Logs</h4>
            <div class="bg-gray-900 rounded p-3 font-mono text-sm text-green-400 whitespace-pre-wrap max-h-32 overflow-y-auto">
              {{ jobLogs }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { getJob } from '../api/client'

const props = defineProps({
  jobId: {
    type: Number,
    required: true
  },
  wsMessage: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close'])

const loadingSteps = ref(false)
const loadingLogs = ref(false)
const steps = ref([])
const selectedStep = ref(null)
const jobLogs = ref('')
let ws = null

function getStepStatusClass(status) {
  const classes = {
    pending: 'px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800',
    running: 'px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800',
    completed: 'px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800',
    failed: 'px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800'
  }
  return classes[status] || classes.pending
}

async function loadJobDetails() {
  loadingSteps.value = true
  try {
    const job = await getJob(props.jobId)
    if (job && job.steps && Array.isArray(job.steps) && job.steps.length > 0) {
      steps.value = job.steps.sort((a, b) => (a.step_order || 0) - (b.step_order || 0))
      if (steps.value.length > 0 && !selectedStep.value) {
        selectedStep.value = steps.value[0]
      }
    } else {
      steps.value = []
    }
    if (job && job.logs) {
      jobLogs.value = job.logs
    } else {
      jobLogs.value = ''
    }
  } catch (error) {
    console.error('Error loading job details:', error)
    steps.value = []
    jobLogs.value = ''
  } finally {
    loadingSteps.value = false
  }
}

watch(() => props.jobId, () => {
  loadJobDetails()
}, { immediate: true })

// Handle WebSocket log updates
watch(() => props.wsMessage, (message) => {
  if (!message || message.type !== 'log_update') return
  
  const data = message.data
  if (data.job_id === props.jobId) {
    // Update job logs
    if (data.logs) {
      jobLogs.value = data.logs
    }
    
    // Update step logs if this is for the selected step
    if (data.step_id && selectedStep.value && selectedStep.value.id === data.step_id) {
      selectedStep.value.logs = data.step_logs
      // Also update in steps array
      const stepIndex = steps.value.findIndex(s => s.id === data.step_id)
      if (stepIndex !== -1) {
        steps.value[stepIndex].logs = data.step_logs
      }
    }
  }
}, { deep: true })

// Setup WebSocket connection for real-time log updates
function setupWebSocket() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/ws`
  
  try {
    ws = new WebSocket(wsUrl)
    
    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        if (message.type === 'log_update' && message.data.job_id === props.jobId) {
          const data = message.data
          
          // Update job logs
          if (data.logs) {
            jobLogs.value = data.logs
          }
          
          // Update step logs if this is for the selected step
          if (data.step_id && selectedStep.value && selectedStep.value.id === data.step_id) {
            selectedStep.value.logs = data.step_logs
            // Also update in steps array
            const stepIndex = steps.value.findIndex(s => s.id === data.step_id)
            if (stepIndex !== -1) {
              steps.value[stepIndex].logs = data.step_logs
            }
          }
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
    
    ws.onclose = () => {
      // Auto-reconnect after 3 seconds
      setTimeout(setupWebSocket, 3000)
    }
  } catch (error) {
    console.error('Error setting up WebSocket:', error)
  }
}

// Auto-refresh logs for running jobs (fallback if WebSocket fails)
let refreshInterval = null
watch(selectedStep, (step) => {
  if (step && step.status === 'running') {
    // Use longer interval since we have WebSocket updates
    refreshInterval = setInterval(loadJobDetails, 30000) // 30 seconds
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
})

onMounted(() => {
  loadJobDetails()
  setupWebSocket()
})

onUnmounted(() => {
  if (ws) {
    ws.close()
  }
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

