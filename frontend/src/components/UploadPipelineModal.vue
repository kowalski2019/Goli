<template>
  <Modal :show="true" title="Upload Pipeline YAML" size="md" @close="$emit('close')">
    <form @submit.prevent="handleSubmit" class="space-y-5">
      <FormField label="Pipeline Name" description="Optional: Will use YAML name if not provided">
        <TextInput
          v-model="formData.name"
          type="text"
          placeholder="Enter pipeline name"
          :disabled="loading"
        />
      </FormField>

      <FormField label="Description" description="Optional: Brief description of the pipeline">
        <TextInput
          v-model="formData.description"
          type="text"
          placeholder="Enter description"
          :disabled="loading"
        />
      </FormField>

      <FormField label="YAML File" required description="Select a .yaml or .yml file">
        <input
          @change="handleFileSelect"
          type="file"
          accept=".yaml,.yml"
          required
          :disabled="loading"
          class="w-full px-4 py-2.5 text-sm border border-gray-300 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white disabled:bg-gray-50 disabled:cursor-not-allowed file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
        />
        <div v-if="formData.file" class="mt-2 text-sm text-gray-600 flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          {{ formData.file.name }}
        </div>
      </FormField>

      <div class="flex items-start gap-3 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <input
          v-model="formData.runImmediately"
          type="checkbox"
          id="runImmediately"
          :disabled="loading"
          class="mt-1 h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
        />
        <label for="runImmediately" class="text-sm text-gray-700 cursor-pointer">
          <span class="font-medium">Run pipeline immediately after upload</span>
          <p class="text-xs text-gray-500 mt-0.5">The pipeline will start executing as soon as it's uploaded</p>
        </label>
      </div>

      <Alert
        v-if="error"
        type="error"
        :message="error"
        dismissible
        @dismiss="error = ''"
      />
    </form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button
          type="button"
          @click="$emit('close')"
          class="btn btn-secondary"
          :disabled="loading"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleSubmit"
          :disabled="loading || !formData.file"
          class="btn btn-primary disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <svg
            v-if="loading"
            class="animate-spin h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ loading ? 'Uploading...' : 'Upload & Create' }}
        </button>
      </div>
    </template>
  </Modal>
</template>

<script setup>
import { ref } from 'vue'
import { uploadPipeline } from '../api/client'
import Modal from './Modal.vue'
import FormField from './FormField.vue'
import TextInput from './TextInput.vue'
import Alert from './Alert.vue'

const emit = defineEmits(['close', 'uploaded'])

const loading = ref(false)
const error = ref('')
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
    error.value = 'Please select a YAML file'
    return
  }

  loading.value = true
  error.value = ''

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
      // Success - pipeline uploaded and started
      emit('uploaded')
    } else {
      // Success - pipeline uploaded
      emit('uploaded')
    }
  } catch (e) {
    error.value = e.message || 'Error uploading pipeline'
  } finally {
    loading.value = false
  }
}
</script>
