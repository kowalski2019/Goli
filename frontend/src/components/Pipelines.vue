<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-3xl font-bold text-white">Pipelines</h2>
      <button @click="showUploadModal = true" class="btn btn-primary">
        + Upload Pipeline
      </button>
    </div>

    <!-- Pipelines List -->
    <div class="card">
      <div v-if="loading" class="text-center py-8 text-gray-500">
        Loading pipelines...
      </div>
      <div v-else-if="pipelines.length === 0" class="text-center py-8 text-gray-500">
        No pipelines found. Upload a YAML file to get started!
      </div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="pipeline in pipelines" :key="pipeline.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">#{{ pipeline.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ pipeline.name }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ pipeline.description || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(pipeline.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                <button 
                  @click="runPipeline(pipeline.id, pipeline.name)"
                  :disabled="runningPipeline === pipeline.id"
                  class="btn btn-primary text-sm"
                >
                  {{ runningPipeline === pipeline.id ? 'Running...' : 'Run' }}
                </button>
                <button 
                  @click="viewPipeline(pipeline.id)"
                  class="btn btn-secondary text-sm"
                >
                  View
                </button>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPipelines, runPipeline } from '../api/client'
import UploadPipelineModal from './UploadPipelineModal.vue'

const loading = ref(false)
const pipelines = ref([])
const showUploadModal = ref(false)
const runningPipeline = ref(null)

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
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

async function runPipelineHandler(id, name) {
  runningPipeline.value = id
  try {
    await runPipeline(id, {
      name: `${name} - Run`,
      triggered_by: 'UI'
    })
    alert('Pipeline started successfully!')
  } catch (error) {
    alert('Error running pipeline: ' + error.message)
  } finally {
    runningPipeline.value = null
  }
}

function viewPipeline(id) {
  // Navigate to pipeline details
  console.log('View pipeline:', id)
}

function handlePipelineUploaded() {
  showUploadModal.value = false
  loadPipelines()
}

onMounted(() => {
  loadPipelines()
})
</script>

