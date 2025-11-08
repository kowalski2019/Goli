<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
      <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">Upload Pipeline YAML</h3>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Pipeline Name (optional)
          </label>
          <input
            v-model="formData.name"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            placeholder="Will use YAML name if not provided"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Description (optional)
          </label>
          <input
            v-model="formData.description"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            YAML File
          </label>
          <input
            @change="handleFileSelect"
            type="file"
            accept=".yaml,.yml"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        </div>
        <div class="flex items-center">
          <input
            v-model="formData.runImmediately"
            type="checkbox"
            id="runImmediately"
            class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
          />
          <label for="runImmediately" class="ml-2 block text-sm text-gray-700">
            Run pipeline immediately after upload
          </label>
        </div>
        <div class="flex justify-end space-x-3 pt-4">
          <button type="button" @click="$emit('close')" class="btn btn-secondary">
            Cancel
          </button>
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? 'Uploading...' : 'Upload & Create' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { uploadPipeline } from '../api/client'

const emit = defineEmits(['close', 'uploaded'])

const loading = ref(false)
const formData = ref({
  name: '',
  description: '',
  file: null,
  runImmediately: false
})

function handleFileSelect(event) {
  formData.value.file = event.target.files[0]
}

async function handleSubmit() {
  if (!formData.value.file) {
    alert('Please select a YAML file')
    return
  }

  loading.value = true
  try {
    const uploadFormData = new FormData()
    uploadFormData.append('yaml_file', formData.value.file)
    if (formData.value.name) {
      uploadFormData.append('name', formData.value.name)
    }
    if (formData.value.description) {
      uploadFormData.append('description', formData.value.description)
    }
    if (formData.value.runImmediately) {
      uploadFormData.append('run', 'true')
    }

    const result = await uploadPipeline(uploadFormData)
    
    if (result.job_started) {
      alert('Pipeline uploaded and started successfully!')
    } else {
      alert('Pipeline uploaded successfully!')
    }
    
    emit('uploaded')
  } catch (error) {
    alert('Error uploading pipeline: ' + error.message)
  } finally {
    loading.value = false
  }
}
</script>

