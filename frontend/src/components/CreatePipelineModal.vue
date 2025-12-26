<template>
  <Modal :show="true" title="Create Pipeline" size="lg" @close="$emit('close')">
    <form @submit.prevent="handleSubmit" class="space-y-5">
      <FormField label="Pipeline Name" required>
        <TextInput
          v-model="formData.name"
          type="text"
          placeholder="Enter pipeline name"
          :disabled="loading"
          required
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

      <FormField label="Pipeline Definition (YAML)" required>
        <textarea
          v-model="formData.definition"
          rows="15"
          placeholder="Enter your pipeline YAML definition here..."
          :disabled="loading"
          required
          class="w-full px-4 py-2.5 text-sm font-mono border border-gray-300 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white disabled:bg-gray-50 disabled:cursor-not-allowed resize-y"
        ></textarea>
      </FormField>

      <!-- Variables Section -->
      <div class="border-t pt-5">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Variables & Secrets</h3>
            <p class="text-sm text-gray-500 mt-1">Define variables that can be used in your pipeline (e.g., ${VAR_NAME} or {{VAR_NAME}})</p>
          </div>
          <button
            type="button"
            @click="addVariable"
            class="btn btn-secondary text-sm px-3 py-1.5 flex items-center gap-1.5"
            :disabled="loading"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Variable
          </button>
        </div>

        <div v-if="formData.variables.length === 0" class="text-sm text-gray-500 text-center py-4 bg-gray-50 rounded-lg">
          No variables defined. Click "Add Variable" to add one.
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="(variable, index) in formData.variables"
            :key="index"
            class="flex items-start gap-3 p-4 bg-gray-50 rounded-lg border border-gray-200"
          >
            <div class="flex-1 space-y-3">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Variable Name</label>
                <TextInput
                  v-model="variable.name"
                  type="text"
                  placeholder="VAR_NAME"
                  :disabled="loading"
                  class="font-mono"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Value</label>
                <TextInput
                  v-model="variable.value"
                  :type="variable.isSecret ? 'password' : 'text'"
                  :placeholder="variable.isSecret ? 'Enter secret value' : 'Enter value'"
                  :disabled="loading"
                />
              </div>
              <div class="flex items-center gap-2">
                <input
                  v-model="variable.isSecret"
                  type="checkbox"
                  :id="`secret-${index}`"
                  :disabled="loading"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                />
                <label :for="`secret-${index}`" class="text-sm text-gray-700 cursor-pointer">
                  Mark as secret (value will be masked)
                </label>
              </div>
            </div>
            <button
              type="button"
              @click="removeVariable(index)"
              class="mt-1 p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              :disabled="loading"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
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
          :disabled="loading || !formData.name || !formData.definition"
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
          {{ loading ? 'Creating...' : 'Create Pipeline' }}
        </button>
      </div>
    </template>
  </Modal>
</template>

<script setup>
import { ref } from 'vue'
import { createPipeline } from '../api/client'
import Modal from './Modal.vue'
import FormField from './FormField.vue'
import TextInput from './TextInput.vue'
import Alert from './Alert.vue'

const emit = defineEmits(['close', 'created'])

const loading = ref(false)
const error = ref('')
const formData = ref({
  name: '',
  description: '',
  definition: '',
  variables: []
})

function addVariable() {
  formData.value.variables.push({
    name: '',
    value: '',
    isSecret: false
  })
}

function removeVariable(index) {
  formData.value.variables.splice(index, 1)
}

async function handleSubmit() {
  if (!formData.value.name || !formData.value.definition) {
    error.value = 'Pipeline name and definition are required'
    return
  }

  loading.value = true
  error.value = ''

  try {
    // Prepare variables object
    const variables = {}
    for (const variable of formData.value.variables) {
      if (variable.name && variable.value) {
        variables[variable.name] = {
          value: variable.value,
          is_secret: variable.isSecret
        }
      }
    }

    const pipelineData = {
      name: formData.value.name,
      description: formData.value.description || '',
      definition: formData.value.definition,
      variables: Object.keys(variables).length > 0 ? variables : undefined
    }

    await createPipeline(pipelineData)
    emit('created')
  } catch (e) {
    error.value = e.message || 'Error creating pipeline'
  } finally {
    loading.value = false
  }
}
</script>

