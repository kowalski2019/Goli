<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
      <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">Create New Job</h3>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Job Name</label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            placeholder="Enter job name"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Triggered By (optional)</label>
          <input
            v-model="formData.triggered_by"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            placeholder="e.g., Manual, Webhook, Schedule"
          />
        </div>
        <div class="flex justify-end space-x-3 pt-4">
          <button type="button" @click="$emit('close')" class="btn btn-secondary">
            Cancel
          </button>
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? 'Creating...' : 'Create Job' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { createJob } from '../api/client'

const emit = defineEmits(['close', 'created'])

const loading = ref(false)
const formData = ref({
  name: '',
  triggered_by: ''
})

async function handleSubmit() {
  loading.value = true
  try {
    await createJob(formData.value)
    emit('created')
  } catch (error) {
    alert('Error creating job: ' + error.message)
  } finally {
    loading.value = false
  }
}
</script>

