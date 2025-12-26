<template>
  <div class="space-y-6 pb-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold text-white dark:text-gray-100 mb-1">Job Logs - #{{ jobId }}</h2>
        <p class="text-gray-300 dark:text-gray-400 text-sm">View detailed logs for each pipeline step</p>
      </div>
      <button
        @click="$emit('close')"
        class="btn btn-secondary flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Back
      </button>
    </div>

    <div class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex h-[600px] border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <!-- Steps Sidebar -->
        <div class="w-64 border-r border-gray-200 dark:border-gray-700 overflow-y-auto bg-gray-50 dark:bg-gray-700/50 scrollbar-thin">
          <div class="p-4">
            <div class="flex items-center justify-between mb-4">
              <h4 class="text-sm font-semibold text-gray-900 dark:text-white">Steps</h4>
              <button
                v-if="selectedStep && (selectedStep.status === 'running' || autoRefresh)"
                @click="autoRefresh = !autoRefresh"
                :class="[
                  'px-2 py-1 text-xs rounded transition-colors',
                  autoRefresh
                    ? 'bg-primary-600 dark:bg-primary-500 text-white'
                    : 'bg-gray-200 dark:bg-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-500'
                ]"
              >
                {{ autoRefresh ? 'Auto: ON' : 'Auto: OFF' }}
              </button>
            </div>
            
            <div v-if="loadingSteps" class="flex items-center justify-center py-8">
              <svg
                class="animate-spin h-6 w-6 text-primary-600 dark:text-primary-400"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </div>
            
            <div v-else-if="steps.length === 0" class="text-center py-8 text-sm text-gray-500 dark:text-gray-400">
              No steps found
            </div>
            
            <div v-else class="space-y-2">
              <button
                v-for="step in steps"
                :key="step.id"
                @click="selectedStep = step"
                :class="[
                  'w-full text-left px-3 py-3 rounded-lg text-sm transition-all duration-150',
                  'border-2',
                  selectedStep?.id === step.id
                    ? 'bg-primary-50 dark:bg-primary-900/50 text-primary-900 dark:text-primary-200 border-primary-300 dark:border-primary-700 shadow-sm'
                    : 'bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500'
                ]"
              >
                <div class="flex items-center justify-between mb-1.5">
                  <div class="flex items-center gap-2">
                    <span class="flex items-center justify-center w-6 h-6 rounded-full bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 text-xs font-semibold">
                      {{ step.step_order }}
                    </span>
                    <span class="font-medium text-gray-900 dark:text-white">{{ step.step_name }}</span>
                  </div>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs text-gray-500 dark:text-gray-400">Step {{ step.step_order }}</span>
                  <StatusBadge :status="step.status" :show-dot="true" />
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Logs Content -->
        <div class="flex-1 flex flex-col overflow-hidden bg-white dark:bg-gray-800">
          <div v-if="selectedStep" class="p-4 border-b border-gray-200 dark:border-gray-700 bg-gradient-to-r from-gray-50 to-white dark:from-gray-700/50 dark:to-gray-800">
            <div class="flex items-center justify-between mb-3">
              <div>
                <h4 class="font-semibold text-gray-900 dark:text-white text-lg">{{ selectedStep.step_name }}</h4>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">Step {{ selectedStep.step_order }}</p>
              </div>
              <StatusBadge :status="selectedStep.status" />
            </div>
            
            <div v-if="selectedStep.error_message" class="mb-3">
              <Alert
                type="error"
                :message="selectedStep.error_message"
                class="mb-0"
              />
            </div>
            
            <div class="flex items-center gap-2">
              <button
                @click="copyLogs"
                class="btn btn-secondary text-sm px-3 py-1.5 flex items-center gap-1.5"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Copy
              </button>
              <button
                @click="downloadLogs"
                class="btn btn-secondary text-sm px-3 py-1.5 flex items-center gap-1.5"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Download
              </button>
              <button
                @click="scrollToBottom"
                class="btn btn-secondary text-sm px-3 py-1.5 flex items-center gap-1.5"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
                Scroll to Bottom
              </button>
            </div>
          </div>

          <div
            ref="logsContainer"
            class="flex-1 overflow-y-auto p-4 bg-gray-900 dark:bg-gray-950 scrollbar-thin"
          >
            <div v-if="loadingLogs" class="flex items-center justify-center h-full">
              <div class="text-center">
                <svg
                  class="animate-spin h-8 w-8 text-green-400 dark:text-green-300 mx-auto mb-2"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p class="text-gray-400 dark:text-gray-500">Loading logs...</p>
              </div>
            </div>
            <div
              v-else-if="selectedStep && selectedStep.logs"
              class="font-mono text-sm text-green-400 dark:text-green-300 whitespace-pre-wrap break-words"
            >
              {{ selectedStep.logs }}
            </div>
            <div
              v-else-if="selectedStep && !selectedStep.logs"
              class="flex items-center justify-center h-full text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <p>No logs available for this step</p>
              </div>
            </div>
            <div
              v-else
              class="flex items-center justify-center h-full text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <svg class="w-12 h-12 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <p>Select a step to view logs</p>
              </div>
            </div>
          </div>

          <!-- Job-level logs -->
          <div v-if="jobLogs" class="border-t border-gray-200 dark:border-gray-700 p-4 bg-gray-50 dark:bg-gray-700/50">
            <div class="flex items-center justify-between mb-2">
              <h4 class="text-sm font-semibold text-gray-900 dark:text-white">Job Logs</h4>
              <button
                @click="copyJobLogs"
                class="text-xs text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300"
              >
                Copy
              </button>
            </div>
            <div class="bg-gray-900 rounded-lg p-3 font-mono text-sm text-green-400 whitespace-pre-wrap break-words max-h-32 overflow-y-auto scrollbar-thin">
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
import StatusBadge from './StatusBadge.vue'
import Alert from './Alert.vue'

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
  const logs = selectedStep.value?.logs || ''
  if (!logs) {
    alert('No logs to copy')
    return
  }

  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(logs)
      alert('Logs copied to clipboard!')
      return
    }

    // Fallback
    const textArea = document.createElement('textarea')
    textArea.value = logs
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()

    try {
      document.execCommand('copy')
      alert('Logs copied to clipboard!')
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

async function copyJobLogs() {
  if (!jobLogs.value) {
    alert('No logs to copy')
    return
  }

  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(jobLogs.value)
      alert('Job logs copied to clipboard!')
    } else {
      alert('Clipboard API not available')
    }
  } catch (error) {
    alert('Failed to copy logs')
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
