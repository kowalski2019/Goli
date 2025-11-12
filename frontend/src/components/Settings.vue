<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-white mb-6">Settings</h2>

    <!-- User Profile Section -->
    <div class="card">
      <h3 class="text-xl font-bold text-gray-900 mb-4">User Profile</h3>
      <div v-if="currentUser" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input v-model="currentUser.username" type="text" disabled class="w-full rounded-lg border-gray-300 bg-gray-100" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="userForm.email" type="email" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
          <input v-model="userForm.phone" type="tel" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <button @click="updateProfile" :disabled="saving" class="btn btn-primary">
          {{ saving ? 'Saving...' : 'Update Profile' }}
        </button>
      </div>
    </div>

    <!-- Two-Factor Authentication Section -->
    <div class="card">
      <h3 class="text-xl font-bold text-gray-900 mb-4">Two-Factor Authentication</h3>
      <div class="space-y-4">
        <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
          <div>
            <h4 class="font-semibold text-gray-900">Email 2FA</h4>
            <p class="text-sm text-gray-600">Receive verification codes via email</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input 
              type="checkbox" 
              v-model="twoFA.email" 
              @change="update2FA"
              class="sr-only peer"
            />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
          </label>
        </div>

        <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
          <div>
            <h4 class="font-semibold text-gray-900">SMS 2FA</h4>
            <p class="text-sm text-gray-600">Receive verification codes via SMS</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input 
              type="checkbox" 
              v-model="twoFA.sms" 
              @change="update2FA"
              class="sr-only peer"
            />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
          </label>
        </div>

        <div v-if="error" class="p-3 bg-red-50 border border-red-200 rounded text-sm text-red-800">
          {{ error }}
        </div>
        <div v-if="success" class="p-3 bg-green-50 border border-green-200 rounded text-sm text-green-800">
          {{ success }}
        </div>
      </div>
    </div>

    <!-- System Configuration Section (Admin only) -->
    <div v-if="currentUser && currentUser.role === 'admin'" class="card">
      <h3 class="text-xl font-bold text-gray-900 mb-4">System Configuration</h3>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Port</label>
          <input v-model.number="config.port" type="number" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">GitHub Username</label>
          <input v-model="config.gh_username" type="text" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">GitHub Access Token</label>
          <input v-model="config.gh_access_token" type="password" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <button @click="updateConfig" :disabled="savingConfig" class="btn btn-primary">
          {{ savingConfig ? 'Saving...' : 'Update Configuration' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUsers, updateUser, getConfig, updateConfig as updateConfigAPI } from '../api/client'

const currentUser = ref(null)
const userForm = ref({
  email: '',
  phone: ''
})
const twoFA = ref({
  email: false,
  sms: false
})
const config = ref({
  port: 8125,
  gh_username: '',
  gh_access_token: ''
})
const saving = ref(false)
const savingConfig = ref(false)
const error = ref('')
const success = ref('')

async function loadUserData() {
  try {
    const users = await getUsers()
    // Get current user (you might want to get this from login response or a /me endpoint)
    // For now, assume first user or get from token
    if (users && users.length > 0) {
      currentUser.value = users[0] // In real app, get current user from auth
      userForm.value.email = currentUser.value.email || ''
      userForm.value.phone = currentUser.value.phone || ''
      twoFA.value.email = currentUser.value.two_fa_email_enabled === 1
      twoFA.value.sms = currentUser.value.two_fa_sms_enabled === 1
    }
  } catch (error) {
    console.error('Error loading user data:', error)
  }
}

async function loadConfig() {
  try {
    const configData = await getConfig()
    config.value.port = parseInt(configData.port) || 8125
    config.value.gh_username = configData.gh_username || ''
    config.value.gh_access_token = configData.gh_access_token || ''
  } catch (error) {
    console.error('Error loading config:', error)
  }
}

async function updateProfile() {
  if (!currentUser.value) return
  
  saving.value = true
  error.value = ''
  success.value = ''
  
  try {
    await updateUser(currentUser.value.id, {
      email: userForm.value.email,
      phone: userForm.value.phone
    })
    success.value = 'Profile updated successfully'
    await loadUserData()
  } catch (err) {
    error.value = err.message || 'Failed to update profile'
  } finally {
    saving.value = false
  }
}

async function update2FA() {
  if (!currentUser.value) return
  
  error.value = ''
  success.value = ''
  
  try {
    await updateUser(currentUser.value.id, {
      two_fa_email_enabled: twoFA.value.email ? 1 : 0,
      two_fa_sms_enabled: twoFA.value.sms ? 1 : 0
    })
    success.value = '2FA settings updated successfully'
  } catch (err) {
    error.value = err.message || 'Failed to update 2FA settings'
    // Revert checkboxes on error
    await loadUserData()
  }
}

async function updateConfig() {
  savingConfig.value = true
  error.value = ''
  success.value = ''
  
  try {
    await updateConfigAPI({
      port: config.value.port.toString(),
      gh_username: config.value.gh_username,
      gh_access_token: config.value.gh_access_token
    })
    success.value = 'Configuration updated successfully'
  } catch (err) {
    error.value = err.message || 'Failed to update configuration'
  } finally {
    savingConfig.value = false
  }
}

onMounted(() => {
  loadUserData()
  loadConfig()
})
</script>

