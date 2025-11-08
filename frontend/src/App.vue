<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700">
    <!-- Header -->
    <header class="bg-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <h1 class="text-2xl font-bold text-gray-900">🚀 Goli CI/CD</h1>
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
          </nav>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <Dashboard v-if="activeTab === 'dashboard'" />
      <Pipelines v-if="activeTab === 'pipelines'" />
      <Jobs v-if="activeTab === 'jobs'" />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createWebSocket } from './api/client'
import Dashboard from './components/Dashboard.vue'
import Pipelines from './components/Pipelines.vue'
import Jobs from './components/Jobs.vue'

const activeTab = ref('dashboard')
const wsConnected = ref(false)

onMounted(() => {
  // Setup WebSocket connection
  createWebSocket((message) => {
    wsConnected.value = true
    // Handle WebSocket messages
    console.log('WebSocket message:', message)
  })
  
  // Set connected after a short delay
  setTimeout(() => {
    wsConnected.value = true
  }, 500)
})
</script>

