<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <button
              @click="$emit('close')"
              class="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
            </button>
              <div>
                <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ pipelineId ? 'Edit Pipeline' : 'Create Pipeline' }}
                </h1>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  {{ pipelineId ? `Editing pipeline #${pipelineId}` : 'Create a new pipeline from scratch' }}
                </p>
              </div>
          </div>
          <div class="flex items-center gap-3">
            <button
              @click="$emit('close')"
              class="btn btn-secondary px-4 py-2"
              :disabled="saving"
            >
              Cancel
            </button>
            <button
              @click="handleSave"
              :disabled="saving || !formData.name || !formData.definition"
              class="btn btn-primary px-6 py-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <svg
                v-if="saving"
                class="animate-spin h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ saving ? 'Saving...' : (pipelineId ? 'Save Changes' : 'Create Pipeline') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left Column: Pipeline Info -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Basic Info -->
          <div class="card">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Pipeline Information</h2>
            <div class="space-y-4">
              <FormField label="Pipeline Name" required>
                <TextInput
                  v-model="formData.name"
                  type="text"
                  placeholder="Enter pipeline name"
                  :disabled="loading || saving"
                  required
                />
              </FormField>

              <FormField label="Description" description="Optional: Brief description of the pipeline">
                <TextInput
                  v-model="formData.description"
                  type="text"
                  placeholder="Enter description"
                  :disabled="loading || saving"
                />
              </FormField>
            </div>
          </div>

          <!-- Variables & Secrets -->
          <div class="card">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Variables & Secrets</h2>
              <button
                type="button"
                @click="addVariable"
                class="btn btn-secondary text-sm px-3 py-1.5 flex items-center gap-1.5"
                :disabled="loading || saving"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                Add Variable
              </button>
            </div>

            <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
              Define variables that can be used in your pipeline using <code class="bg-gray-100 dark:bg-gray-700 px-1 rounded text-gray-900 dark:text-gray-100">${VAR_NAME}</code> or <code class="bg-gray-100 dark:bg-gray-700 px-1 rounded text-gray-900 dark:text-gray-100">{{VAR_NAME}}</code> syntax.
            </p>

            <div v-if="formData.variables.length === 0" class="text-sm text-gray-500 dark:text-gray-400 text-center py-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
              No variables defined. Click "Add Variable" to add one.
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="(variable, index) in formData.variables"
                :key="index"
                class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600 space-y-3"
              >
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Variable Name</label>
                  <TextInput
                    v-model="variable.name"
                    type="text"
                    placeholder="VAR_NAME"
                    :disabled="loading || saving"
                    class="font-mono"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Value</label>
                  <TextInput
                    v-model="variable.value"
                    :type="variable.isSecret && variable.value !== '***MASKED***' ? 'password' : 'text'"
                    :placeholder="variable.isSecret && variable.value === '***MASKED***' ? 'Secret value (enter new value to update)' : variable.isSecret ? 'Enter secret value' : 'Enter value'"
                    :disabled="loading || saving"
                  />
                  <p v-if="variable.isSecret && variable.value === '***MASKED***'" class="text-xs text-gray-500 mt-1">
                    Secret value is masked. Enter a new value to update it.
                  </p>
                </div>
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <input
                      v-model="variable.isSecret"
                      type="checkbox"
                      :id="`secret-${index}`"
                      :disabled="loading || saving"
                      class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <label :for="`secret-${index}`" class="text-sm text-gray-700 cursor-pointer">
                      Mark as secret
                    </label>
                  </div>
                  <button
                    type="button"
                    @click="removeVariable(index)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    :disabled="loading || saving"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column: Code Editor -->
        <div class="lg:col-span-2">
          <div class="card p-0 overflow-hidden">
            <div class="bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600 px-4 py-3 flex items-center justify-between">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Pipeline Definition (YAML)</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Edit your pipeline YAML definition with syntax highlighting</p>
              </div>
            </div>
            <div ref="editorContainer" class="codemirror-container"></div>
          </div>
        </div>
      </div>

      <!-- Error Alert -->
      <Alert
        v-if="error"
        type="error"
        :message="error"
        dismissible
        @dismiss="error = ''"
        class="mt-6"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { EditorView, lineNumbers, keymap } from '@codemirror/view'
import { EditorState } from '@codemirror/state'
import { yaml } from '@codemirror/lang-yaml'
import { oneDark } from '@codemirror/theme-one-dark'
import { defaultKeymap } from '@codemirror/commands'
import { useTheme } from '../composables/useTheme'
import { createPipeline, updatePipeline, getPipeline } from '../api/client'
import FormField from './FormField.vue'
import TextInput from './TextInput.vue'
import Alert from './Alert.vue'

const props = defineProps({
  pipelineId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['close', 'saved'])

const theme = useTheme()
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const editorContainer = ref(null)
let editorView = null

const formData = ref({
  name: '',
  description: '',
  definition: '',
  variables: []
})

// Initialize CodeMirror editor
function initEditor() {
  if (!editorContainer.value) return

  const updateListener = EditorView.updateListener.of((update) => {
    if (update.docChanged) {
      formData.value.definition = update.state.doc.toString()
    }
  })

  const extensions = [
    lineNumbers(),
    yaml(),
    keymap.of(defaultKeymap),
    updateListener,
    EditorView.theme({
      '&': {
        fontSize: '14px',
        minHeight: '600px'
      },
      '.cm-content': {
        padding: '16px',
        minHeight: '600px',
        fontFamily: 'Monaco, Menlo, "Ubuntu Mono", Consolas, "source-code-pro", monospace'
      },
      '.cm-scroller': {
        overflow: 'auto'
      },
      '.cm-editor': {
        height: '100%'
      },
      '.cm-editor.cm-focused': {
        outline: 'none'
      }
    })
  ]

  // Apply dark theme if global theme is dark
  if (theme.isDark.value) {
    extensions.push(oneDark)
  }

  const state = EditorState.create({
    doc: formData.value.definition,
    extensions: extensions
  })

  editorView = new EditorView({
    state: state,
    parent: editorContainer.value
  })
}

// Watch for theme changes and update editor
watch(() => theme.isDark.value, (isDark) => {
  if (editorView) {
    // Recreate editor with new theme
    const currentContent = editorView.state.doc.toString()
    editorView.destroy()
    editorView = null
    // Small delay to ensure cleanup
    setTimeout(() => {
      initEditor()
      // Restore content after initialization
      if (editorView) {
        const state = editorView.state
        editorView.dispatch({
          changes: {
            from: 0,
            to: state.doc.length,
            insert: currentContent
          }
        })
      }
    }, 10)
  }
})

// Load pipeline if editing
async function loadPipeline() {
  if (!props.pipelineId) {
    // Set default YAML template for new pipelines
    formData.value.definition = `name: "My Pipeline"
description: "Pipeline description"
steps:
  - name: "Example Step"
    type: "docker"
    action: "pull"
    config:
      image: "nginx:latest"
`
    initEditor()
    return
  }

  loading.value = true
  try {
    const pipeline = await getPipeline(props.pipelineId)
    formData.value.name = pipeline.name || ''
    formData.value.description = pipeline.description || ''
    formData.value.definition = pipeline.definition || ''
    
    // Load variables
    if (pipeline.variables && typeof pipeline.variables === 'object') {
      formData.value.variables = Object.entries(pipeline.variables).map(([name, value]) => {
        const isSecret = value === '***MASKED***'
        return {
          name,
          value: isSecret ? '***MASKED***' : String(value),
          isSecret
        }
      })
    } else {
      formData.value.variables = []
    }

    // Initialize or update editor
    if (!editorView) {
      initEditor()
    } else {
      // Update editor content
      const currentContent = editorView.state.doc.toString()
      if (currentContent !== formData.value.definition) {
        editorView.dispatch({
          changes: {
            from: 0,
            to: currentContent.length,
            insert: formData.value.definition
          }
        })
      }
    }
  } catch (e) {
    error.value = e.message || 'Error loading pipeline'
  } finally {
    loading.value = false
  }
}

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

async function handleSave() {
  if (!formData.value.name || !formData.value.definition) {
    error.value = 'Pipeline name and definition are required'
    return
  }

  saving.value = true
  error.value = ''

  try {
    // Prepare variables object
    const variables = {}
    for (const variable of formData.value.variables) {
      if (variable.name) {
        variables[variable.name] = {
          value: variable.value || '',
          is_secret: variable.isSecret
        }
      }
    }

    const pipelineData = {
      name: formData.value.name,
      description: formData.value.description || '',
      definition: formData.value.definition,
      variables: variables
    }

    if (props.pipelineId) {
      await updatePipeline(props.pipelineId, pipelineData)
    } else {
      await createPipeline(pipelineData)
    }

    emit('saved')
  } catch (e) {
    error.value = e.message || (props.pipelineId ? 'Error updating pipeline' : 'Error creating pipeline')
  } finally {
    saving.value = false
  }
}

// Watch for definition changes from external sources (like loading pipeline)
watch(() => formData.value.definition, (newVal) => {
  if (editorView) {
    const currentContent = editorView.state.doc.toString()
    // Only update if different to avoid infinite loop
    if (currentContent !== newVal) {
      editorView.dispatch({
        changes: {
          from: 0,
          to: currentContent.length,
          insert: newVal
        }
      })
    }
  }
})

onMounted(() => {
  loadPipeline()
})

onBeforeUnmount(() => {
  if (editorView) {
    editorView.destroy()
  }
})
</script>

<style scoped>
.codemirror-container {
  width: 100%;
  height: 600px;
  overflow: auto;
}

.codemirror-container :deep(.cm-editor) {
  height: 100%;
}

.codemirror-container :deep(.cm-scroller) {
  overflow: auto;
}
</style>

