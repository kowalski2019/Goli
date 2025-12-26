<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
    <!-- Setup Wizard (shown if setup not complete) -->
    <SetupWizard v-if="showSetupWizard" @setup-complete="handleSetupComplete" @setup-already-complete="() => handleSetupComplete(true)" />

    <!-- Main App (shown if setup is complete and authenticated) -->
    <template v-else>
      <template v-if="isAuthenticated">
        <!-- Header -->
        <header class="bg-white dark:bg-gray-800 shadow-lg">
          <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <img src="/goli-logo.jpg" alt="Goli Logo" class="h-10 w-auto" />
                <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Goli CI/CD</h1>
                <div class="flex items-center space-x-2">
                  <div 
                    :class="[
                      'w-3 h-3 rounded-full',
                      wsConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
                    ]"
                  ></div>
                  <span class="text-sm text-gray-600 dark:text-gray-300">
                    {{ wsConnected ? 'Online' : 'Offline' }}
                  </span>
                </div>
              </div>
              <nav class="flex items-center space-x-4">
                <button
                  @click="activeTab = 'dashboard'"
                  :class="[
                    'px-4 py-2 rounded-lg transition-colors',
                    activeTab === 'dashboard' 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  ]"
                >
                  Dashboard
                </button>
                <button
                  @click="activeTab = 'pipelines'"
                  :class="[
                    'px-4 py-2 rounded-lg transition-colors',
                    activeTab === 'pipelines' 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  ]"
                >
                  Pipelines
                </button>
                <button
                  @click="activeTab = 'jobs'"
                  :class="[
                    'px-4 py-2 rounded-lg transition-colors',
                    activeTab === 'jobs' 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  ]"
                >
                  Jobs
                </button>
                <button
                  @click="activeTab = 'settings'"
                  :class="[
                    'px-4 py-2 rounded-lg transition-colors',
                    activeTab === 'settings' 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  ]"
                >
                  Settings
                </button>
                <button
                  @click="theme.toggleTheme()"
                  class="p-2 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                  :title="theme.currentTheme === 'light' ? 'Switch to dark mode' : theme.currentTheme === 'dark' ? 'Switch to system theme' : 'Switch to light mode'"
                >
                  <svg v-if="theme.currentTheme === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                  </svg>
                  <svg v-else-if="theme.currentTheme === 'dark'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                </button>
                <button
                  @click="handleLogout"
                  class="px-4 py-2 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                >
                  Logout
                </button>
              </nav>
            </div>
          </div>
        </header>

        <!-- Main Content -->
        <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Dashboard v-if="activeTab === 'dashboard'" @view-logs="showLogsView" />
          <Pipelines 
            v-if="activeTab === 'pipelines'" 
            @view-logs="showLogsView"
            @create-pipeline="showPipelineEditor"
            @edit-pipeline="showPipelineEditor"
          />
          <Jobs v-if="activeTab === 'jobs'" @view-logs="showLogsView" />
          <Settings v-if="activeTab === 'settings'" />
          <LogsView v-if="activeTab === 'logs'" :job-id="logsJobId" @close="closeLogsView" />
          <PipelineEditor 
            v-if="activeTab === 'pipeline-editor'" 
            :pipeline-id="editingPipelineId"
            @close="closePipelineEditor"
            @saved="handlePipelineSaved"
          />
        </main>
      </template>
      <Login v-else @logged-in="handleLoggedIn" />
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { createWebSocket, getToken, getSetupStatus, logout, AUTO_LOGOUT_EVENT } from './api/client'
import { useTheme } from './composables/useTheme'
import Dashboard from './components/Dashboard.vue'
import Pipelines from './components/Pipelines.vue'
import Jobs from './components/Jobs.vue'
import Settings from './components/Settings.vue'
import SetupWizard from './components/SetupWizard.vue'
import LogsView from './components/LogsView.vue'
import Login from './components/Login.vue'
import PipelineEditor from './components/PipelineEditor.vue'

const theme = useTheme()

const activeTab = ref('dashboard')
const previousTab = ref('dashboard')
const wsConnected = ref(false)
const showSetupWizard = ref(false)
const logsJobId = ref(null)
const editingPipelineId = ref(null)
const isAuthenticated = ref(false)

function showLogsView(jobId) {
  previousTab.value = activeTab.value
  logsJobId.value = jobId
  activeTab.value = 'logs'
}

function closeLogsView() {
  activeTab.value = previousTab.value
  logsJobId.value = null
}

function showPipelineEditor(pipelineId = null) {
  previousTab.value = activeTab.value
  editingPipelineId.value = pipelineId
  activeTab.value = 'pipeline-editor'
}

function closePipelineEditor() {
  activeTab.value = previousTab.value
  editingPipelineId.value = null
}

function handlePipelineSaved() {
  // Return to pipelines tab and refresh
  activeTab.value = 'pipelines'
  editingPipelineId.value = null
  // Trigger a refresh by emitting an event that Pipelines component can listen to
  // For now, we'll just navigate back - the component will refresh on mount
}

async function checkSetupStatus() {
  try {
    const status = await getSetupStatus()
    // Check if setup_complete is false or undefined
    showSetupWizard.value = !status.setup_complete || status.setup_complete === false
  } catch (error) {
    console.error('Error checking setup status:', error)
    // If we can't check, assume setup is needed
    showSetupWizard.value = true
  }
}

function handleSetupComplete(skipReload = false) {
  showSetupWizard.value = false
  // Only reload if setup was just completed, not if it was already complete
  if (!skipReload) {
    // Optionally reload the page to ensure everything is fresh
    window.location.reload()
  }
}

function handleLoggedIn() {
  isAuthenticated.value = true
  // after login, setup websocket
  createWebSocket((message) => {
    wsConnected.value = true
    console.log('WebSocket message:', message)
  })
}

async function handleLogout() {
  try {
    await logout()
    isAuthenticated.value = false
    activeTab.value = 'dashboard'
  } catch (error) {
    console.error('Logout error:', error)
    // Even if logout fails, clear local state
    isAuthenticated.value = false
  }
}

// Handle automatic logout when session expires
function handleAutoLogout() {
  console.log('Session expired. Logging out automatically.')
  isAuthenticated.value = false
  activeTab.value = 'dashboard'
  wsConnected.value = false
}

onMounted(() => {
  // Check setup status first
  checkSetupStatus()

  // Auth check
  isAuthenticated.value = !!getToken()

  // Listen for automatic logout events (when token expires)
  window.addEventListener(AUTO_LOGOUT_EVENT, handleAutoLogout)

  // Setup WebSocket connection (only if setup is complete and authenticated)
  if (!showSetupWizard.value && isAuthenticated.value) {
    createWebSocket((message) => {
      wsConnected.value = true
      // Handle WebSocket messages
      console.log('WebSocket message:', message)
    })
    
    // Set connected after a short delay
    setTimeout(() => {
      wsConnected.value = true
    }, 500)
  }
})

onUnmounted(() => {
  // Clean up event listener
  window.removeEventListener(AUTO_LOGOUT_EVENT, handleAutoLogout)
})
</script>
