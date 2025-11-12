<template>
  <Modal :show="true" title="Create New Job" size="md" @close="$emit('close')">
    <form @submit.prevent="handleSubmit" class="space-y-5">
      <FormField label="Job Name" required>
        <TextInput
          v-model="formData.name"
          type="text"
          placeholder="Enter job name"
          required
          :disabled="loading"
        />
      </FormField>

      <FormField label="Triggered By" description="Optional: e.g., Manual, Webhook, Schedule">
        <TextInput
          v-model="formData.triggered_by"
          type="text"
          placeholder="e.g., Manual, Webhook, Schedule"
          :disabled="loading"
        />
      </FormField>

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
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleSubmit"
          :disabled="loading || !formData.name"
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
          {{ loading ? 'Creating...' : 'Create Job' }}
        </button>
      </div>
    </template>
  </Modal>
</template>

<script setup>
import { ref } from 'vue'
import { createJob } from '../api/client'
import Modal from './Modal.vue'
import FormField from './FormField.vue'
import TextInput from './TextInput.vue'
import Alert from './Alert.vue'

const emit = defineEmits(['close', 'created'])

const loading = ref(false)
const error = ref('')
const formData = ref({
  name: '',
  triggered_by: ''
})

async function handleSubmit() {
  if (!formData.value.name) {
    error.value = 'Job name is required'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    await createJob(formData.value)
    emit('created')
  } catch (e) {
    error.value = e.message || 'Error creating job'
  } finally {
    loading.value = false
  }
}
</script>
