<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold text-white">Job Logs - #{{ jobId }}</h2>
      <button @click="$emit('close')" class="btn btn-secondary">
        ← Back
      </button>
    </div>

    <div class="card">
      <div class="flex h-[600px]">
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
              <div class="flex items-center space-x-2">
                <span :class="getStepStatusClass(selectedStep.status)">
                  {{ selectedStep.status }}
                </span>
                <button
                  @click="copyLogs"
                  class="px-3 py-1 text-sm bg-primary-600 text-white rounded hover:bg-primary-700"
                >
                  Copy
                </button>
                <button
                  @click="downloadLogs"
                  class="px-3 py-1 text-sm bg-gray-600 text-white rounded hover:bg-gray-700"
                >
                  Download
                </button>
              </div>
            </div>
            <div v-if="selectedStep.error_message" class="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-800">
              <strong>Error:</strong> {{ selectedStep.error_message }}
            </div>
          </div>

          <div ref="logsContainer" class="flex-1 overflow-y-auto p-4 bg-gray-900">
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
import { ref, onMounted, watch, onUnmounted, nextTick } from 'vue'
import { getJob } from '../api/client'

const props = defineProps({
  jobId: {
    type: Number,
    required: true
  }
})

const emit = defineEmits(['close'])

const loadingSteps = ref(false)
const loadingLogs = ref(false)
const steps = ref([])
const selectedStep = ref(null)
const jobLogs = ref('')
const autoRefresh = ref(false)
const logsContainer = ref(null)
let refreshInterval = null

function getStepStatusClass(status) {
  const classes = {
    pending: 'px-2 py-1 text-xs font-semibold rounded-full bg-yellow-100 text-yellow-800',
    running: 'px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800',
    completed: 'px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800',
    failed: 'px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800',
    cancelled: 'px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-800'
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
      } else if (selectedStep.value) {
        // Update selected step with latest data
        const updatedStep = steps.value.find(s => s.id === selectedStep.value.id)
        if (updatedStep) {
          selectedStep.value = updatedStep
        }
      }
    } else {
      steps.value = []
    }
    if (job && job.logs) {
      jobLogs.value = job.logs
    } else {
      jobLogs.value = ''
    }
    
    // Auto-scroll to bottom if auto-refresh is on
    if (autoRefresh.value) {
      await nextTick()
      scrollToBottom()
    }
  } catch (error) {
    console.error('Error loading job details:', error)
    steps.value = []
    jobLogs.value = ''
  } finally {
    loadingSteps.value = false
  }
}

function scrollToBottom() {
  if (logsContainer.value) {
    logsContainer.value.scrollTop = logsContainer.value.scrollHeight
  }
}

async function copyLogs() {
  const logs = selectedStep.value?.logs || jobLogs.value || ''
  if (!logs) {
    alert('No logs to copy')
    return
  }

  try {
    // Try modern Clipboard API first
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(logs)
      alert('Logs copied to clipboard!')
      return
    }

    // Fallback to older method
    const textArea = document.createElement('textarea')
    textArea.value = logs
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    
    try {
      const successful = document.execCommand('copy')
      if (successful) {
        alert('Logs copied to clipboard!')
      } else {
        throw new Error('Copy command failed')
      }
    } catch (err) {
      alert('Failed to copy logs. Please select and copy manually.')
    } finally {
      document.body.removeChild(textArea)
    }
  } catch (error) {
    console.error('Error copying logs:', error)
    alert('Failed to copy logs. Please select and copy manually.')
  }
}

function downloadLogs() {
  const logs = selectedStep.value?.logs || jobLogs.value || ''
  if (logs) {
    const blob = new Blob([logs], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `job-${props.jobId}-${selectedStep.value?.step_name || 'logs'}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }
}

watch(() => props.jobId, () => {
  loadJobDetails()
}, { immediate: true })

watch(selectedStep, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

watch(autoRefresh, (enabled) => {
  if (enabled) {
    refreshInterval = setInterval(() => {
      loadJobDetails()
    }, 2000)
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
})

onMounted(() => {
  loadJobDetails()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

