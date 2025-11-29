<template>
  <div class="space-y-6 pb-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold text-white mb-1">Pipelines</h2>
        <p class="text-gray-300 text-sm">Manage and execute your CI/CD pipelines</p>
      </div>
      <button
        @click="showUploadModal = true"
        class="btn btn-primary px-6 py-2.5 flex items-center gap-2 hover:shadow-lg transition-shadow"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Upload Pipeline
      </button>
    </div>

    <!-- Pipelines List -->
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
        <p class="text-gray-500 text-sm">Loading pipelines...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="pipelines.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="p-4 bg-gray-100 rounded-full mb-4">
          <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-1">No pipelines found</h3>
        <p class="text-gray-500 text-sm mb-4">Upload a YAML file to get started!</p>
        <button @click="showUploadModal = true" class="btn btn-primary">
          Upload Your First Pipeline
        </button>
      </div>

      <!-- Pipelines Table -->
      <div v-else class="overflow-x-auto scrollbar-thin">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                ID
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Name
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Description
              </th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Created
              </th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="pipeline in pipelines"
              :key="pipeline.id"
              class="hover:bg-gray-50 transition-colors duration-150"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm font-medium text-gray-900">#{{ pipeline.id }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span class="text-sm font-medium text-gray-900">{{ pipeline.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <span class="text-sm text-gray-500 line-clamp-1">
                  {{ pipeline.description || '-' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm text-gray-500">
                  {{ formatDate(pipeline.created_at) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                <div class="flex items-center justify-end gap-2">
                  <button
                    @click="runPipeline(pipeline.id, pipeline.name)"
                    :disabled="runningPipeline === pipeline.id"
                    class="btn btn-primary text-sm px-4 py-1.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
                  >
                    <svg
                      v-if="runningPipeline === pipeline.id"
                      class="animate-spin h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <svg
                      v-else
                      class="h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ runningPipeline === pipeline.id ? 'Running...' : 'Run' }}
                  </button>
                  <button
                    @click="viewPipeline(pipeline.id)"
                    class="btn btn-secondary text-sm px-4 py-1.5 flex items-center gap-1.5"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    View
                  </button>
                  <button
                    @click="confirmDeletePipeline(pipeline.id, pipeline.name)"
                    :disabled="deletingPipeline === pipeline.id"
                    class="btn bg-red-600 hover:bg-red-700 text-white text-sm px-4 py-1.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
                  >
                    <svg
                      v-if="deletingPipeline === pipeline.id"
                      class="animate-spin h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <svg
                      v-else
                      class="h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    {{ deletingPipeline === pipeline.id ? 'Deleting...' : 'Delete' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Upload Modal -->
    <UploadPipelineModal
      v-if="showUploadModal"
      @close="showUploadModal = false"
      @uploaded="handlePipelineUploaded"
    />

    <!-- Pipeline Detail Modal -->
    <PipelineDetailModal
      v-if="selectedPipeline"
      :pipeline="selectedPipeline"
      @close="selectedPipeline = null"
      @view-logs="handleViewLogs"
      @run-pipeline="handleRunPipeline"
    />

    <!-- Delete Confirmation Modal -->
    <div
      v-if="deleteConfirmPipeline"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="cancelDelete"
    >
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div class="p-6">
          <div class="flex items-center mb-4">
            <div class="flex-shrink-0 mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
              <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
          </div>
          <h3 class="text-lg font-medium text-gray-900 text-center mb-2">Delete Pipeline</h3>
          <p class="text-sm text-gray-500 text-center mb-4">
            Are you sure you want to delete <strong>{{ deleteConfirmPipeline.name }}</strong>?
          </p>
          <p class="text-xs text-red-600 text-center mb-6 bg-red-50 p-3 rounded">
            ⚠️ This will permanently delete the pipeline and all associated jobs and job steps. This action cannot be undone.
          </p>
          <div class="flex gap-3">
            <button
              @click="cancelDelete"
              class="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="deletePipeline(deleteConfirmPipeline.id)"
              :disabled="deletingPipeline === deleteConfirmPipeline.id"
              class="flex-1 px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors"
            >
              {{ deletingPipeline === deleteConfirmPipeline.id ? 'Deleting...' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPipelines, runPipeline as runPipelineAPI, deletePipeline as deletePipelineAPI } from '../api/client'
import UploadPipelineModal from './UploadPipelineModal.vue'
import PipelineDetailModal from './PipelineDetailModal.vue'

const emit = defineEmits(['view-logs'])

const loading = ref(false)
const pipelines = ref([])
const showUploadModal = ref(false)
const runningPipeline = ref(null)
const selectedPipeline = ref(null)
const deletingPipeline = ref(null)
const deleteConfirmPipeline = ref(null)

function formatDate(dateString) {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function loadPipelines() {
  loading.value = true
  try {
    const pipelinesData = await getPipelines()
    pipelines.value = Array.isArray(pipelinesData) ? pipelinesData : []
  } catch (error) {
    console.error('Error loading pipelines:', error)
    pipelines.value = []
  } finally {
    loading.value = false
  }
}

async function runPipeline(id, name) {
  runningPipeline.value = id
  try {
    await runPipelineAPI(id, {
      name: `${name} - Run`,
      triggered_by: 'UI'
    })
    // Show success notification (you could use a toast library here)
    await loadPipelines() // Refresh to show new job
  } catch (error) {
    alert('Error running pipeline: ' + error.message)
  } finally {
    runningPipeline.value = null
  }
}

function viewPipeline(id) {
  const pipeline = pipelines.value.find(p => p.id === id)
  if (pipeline) {
    selectedPipeline.value = pipeline
  }
}

function handleViewLogs(jobId) {
  emit('view-logs', jobId)
  selectedPipeline.value = null
}

function handleRunPipeline(pipelineId) {
  const pipeline = pipelines.value.find(p => p.id === pipelineId)
  if (pipeline) {
    selectedPipeline.value = null
    runPipeline(pipelineId, pipeline.name)
  }
}

function handlePipelineUploaded() {
  showUploadModal.value = false
  loadPipelines()
}

function confirmDeletePipeline(id, name) {
  deleteConfirmPipeline.value = { id, name }
}

function cancelDelete() {
  deleteConfirmPipeline.value = null
}

async function deletePipeline(id) {
  deletingPipeline.value = id
  try {
    await deletePipelineAPI(id)
    deleteConfirmPipeline.value = null
    await loadPipelines()
    // Show success message (you could use a toast library here)
  } catch (error) {
    alert('Error deleting pipeline: ' + error.message)
  } finally {
    deletingPipeline.value = null
  }
}

onMounted(() => {
  loadPipelines()
})
</script>
