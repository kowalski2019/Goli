<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700">
    <!-- Setup Wizard (shown if setup not complete) -->
    <SetupWizard v-if="showSetupWizard" @setup-complete="handleSetupComplete" @setup-already-complete="() => handleSetupComplete(true)" />

    <!-- Main App (shown if setup is complete and authenticated) -->
    <template v-else>
      <template v-if="isAuthenticated">
        <!-- Header -->
        <header class="bg-white shadow-lg">
          <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <img src="/goli-logo.jpg" alt="Goli Logo" class="h-10 w-auto" />
                <h1 class="text-2xl font-bold text-gray-900">Goli CI/CD</h1>
                <div class="flex items-center space-x-2">
                  <div 
                    :class="[
                      'w-3 h-3 rounded-full',
                      wsConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'
                    ]"
                  ></div>
                  <span class="text-sm text-gray-600">
                    {{ wsConnected ? 'Online' : 'Offline' }}
                  </span>
                </div>
              </div>
              <nav class="flex space-x-4">
                <button
                  @click="activeTab = 'dashboard'"
                  :class="[
                    'px-4 py-2 rounded-lg transition-colors',
                    activeTab === 'dashboard' 
                      ? 'bg-primary-600 text-white' 
                      : 'text-gray-700 hover:bg-gray-100'
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
                      : 'text-gray-700 hover:bg-gray-100'
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
                      : 'text-gray-700 hover:bg-gray-100'
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
                      : 'text-gray-700 hover:bg-gray-100'
                  ]"
                >
                  Settings
                </button>
                <button
                  @click="handleLogout"
                  class="px-4 py-2 rounded-lg text-gray-700 hover:bg-gray-100"
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
          <Pipelines v-if="activeTab === 'pipelines'" @view-logs="showLogsView" />
          <Jobs v-if="activeTab === 'jobs'" @view-logs="showLogsView" />
          <Settings v-if="activeTab === 'settings'" />
          <LogsView v-if="activeTab === 'logs'" :job-id="logsJobId" @close="closeLogsView" />
        </main>
      </template>
      <Login v-else @logged-in="handleLoggedIn" />
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createWebSocket, getToken, getSetupStatus, logout } from './api/client'
import Dashboard from './components/Dashboard.vue'
import Pipelines from './components/Pipelines.vue'
import Jobs from './components/Jobs.vue'
import Settings from './components/Settings.vue'
import SetupWizard from './components/SetupWizard.vue'
import LogsView from './components/LogsView.vue'
import Login from './components/Login.vue'

const activeTab = ref('dashboard')
const previousTab = ref('dashboard')
const wsConnected = ref(false)
const showSetupWizard = ref(false)
const logsJobId = ref(null)
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

onMounted(() => {
  // Check setup status first
  checkSetupStatus()

  // Auth check
  isAuthenticated.value = !!getToken()

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
</script>
