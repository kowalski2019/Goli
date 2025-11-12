<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-600 via-blue-600 to-indigo-700 p-4">
    <div class="w-full max-w-2xl bg-white rounded-2xl shadow-2xl p-8">
      <div class="flex flex-col items-center mb-6">
        <img src="/goli-logo.jpg" alt="Goli Logo" class="h-16 w-16 rounded-full shadow-md mb-3" />
        <h2 class="text-2xl font-bold text-gray-900">Goli Setup Wizard</h2>
        <p class="text-gray-500">Complete the setup to get started</p>
      </div>

      <!-- Step 0: Setup Password -->
      <div v-if="currentStep === 0" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Setup Password</label>
          <p class="text-xs text-gray-500 mb-2">Enter the setup password shown during installation</p>
          <input 
            v-model="setupPassword" 
            type="password" 
            required 
            class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500"
            placeholder="Enter setup password"
          />
          <div v-if="passwordError" class="mt-2 text-sm text-red-600">{{ passwordError }}</div>
        </div>
        <button 
          @click="verifySetupPasswordHandler" 
          :disabled="isVerifying"
          class="w-full py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition disabled:opacity-50"
        >
          {{ isVerifying ? 'Verifying...' : 'Continue' }}
        </button>
      </div>

      <!-- Step 1: Admin User -->
      <div v-if="currentStep === 1" class="space-y-4">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Create Admin User</h3>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input v-model="adminUser.username" type="text" required class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="adminUser.email" type="email" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input v-model="adminUser.password" type="password" required class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
          <input v-model="adminUser.confirmPassword" type="password" required class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div v-if="step1Error" class="text-sm text-red-600">{{ step1Error }}</div>
        <div class="flex space-x-3">
          <button @click="prevStep" class="flex-1 py-2.5 rounded-lg bg-gray-200 text-gray-700 hover:bg-gray-300 transition">
            Back
          </button>
          <button @click="nextStep" :disabled="!canProceedStep1" class="flex-1 py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition disabled:opacity-50">
            Next
          </button>
        </div>
      </div>

      <!-- Step 2: System Settings -->
      <div v-if="currentStep === 2" class="space-y-4">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">System Settings</h3>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Port</label>
          <input v-model.number="settings.port" type="number" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Auth Key</label>
          <div class="flex space-x-2">
            <input v-model="settings.authKey" type="text" readonly class="flex-1 rounded-lg border-gray-300 bg-gray-100" />
            <button @click="generateNewAuthKey" class="px-4 py-2 rounded-lg bg-gray-200 text-gray-700 hover:bg-gray-300 transition">
              Generate
            </button>
          </div>
        </div>
        <div class="flex space-x-3">
          <button @click="prevStep" class="flex-1 py-2.5 rounded-lg bg-gray-200 text-gray-700 hover:bg-gray-300 transition">
            Back
          </button>
          <button @click="nextStep" class="flex-1 py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
            Next
          </button>
        </div>
      </div>

      <!-- Step 3: Review & Complete -->
      <div v-if="currentStep === 3" class="space-y-4">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Review & Complete</h3>
        <div class="bg-gray-50 rounded-lg p-4 space-y-2">
          <div><strong>Username:</strong> {{ adminUser.username }}</div>
          <div><strong>Email:</strong> {{ adminUser.email || 'Not provided' }}</div>
          <div><strong>Port:</strong> {{ settings.port }}</div>
        </div>
        <div v-if="setupError" class="text-sm text-red-600">{{ setupError }}</div>
        <div class="flex space-x-3">
          <button @click="prevStep" class="flex-1 py-2.5 rounded-lg bg-gray-200 text-gray-700 hover:bg-gray-300 transition">
            Back
          </button>
          <button 
            @click="completeSetup" 
            :disabled="isCompleting"
            class="flex-1 py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition disabled:opacity-50"
          >
            {{ isCompleting ? 'Completing...' : 'Complete Setup' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { createUser, getConfig, updateConfig, verifySetupPassword, clearAuthorizationKey } from '../api/client'

const emit = defineEmits(['setup-complete', 'setup-already-complete'])

const currentStep = ref(0)
const isCompleting = ref(false)
const isVerifying = ref(false)
const step1Error = ref('')
const setupError = ref('')
const passwordError = ref('')
const setupPassword = ref('')

const adminUser = ref({
  username: 'goli',
  password: '',
  confirmPassword: '',
  email: ''
})

const settings = ref({
  port: 8125,
  authKey: ''
})

const canProceedStep1 = computed(() => {
  return adminUser.value.username &&
         adminUser.value.password &&
         adminUser.value.password === adminUser.value.confirmPassword &&
         adminUser.value.password.length >= 6
})

function generateNewAuthKey() {
  // Generate a random auth key (similar to backend)
  const randomBytes = new Uint8Array(64)
  crypto.getRandomValues(randomBytes)
  const base64 = btoa(String.fromCharCode(...randomBytes))
  // Simple hash simulation (in real app, use proper hashing)
  settings.value.authKey = base64.substring(0, 64)
}

async function verifySetupPasswordHandler() {
  if (!setupPassword.value) {
    passwordError.value = 'Please enter the setup password'
    return
  }

  isVerifying.value = true
  passwordError.value = ''

  try {
    const result = await verifySetupPassword(setupPassword.value)
    // Password verified - auth key is automatically stored by verifySetupPassword
    // Load config now that we have the auth key
    await loadConfig()
    // Proceed to next step
    currentStep.value = 1
  } catch (error) {
    const errorMessage = error.message || 'Invalid setup password. Please check the installation output.'
    
    // If setup is already completed, jump to login page
    if (errorMessage.includes('Setup has already been completed')) {
      emit('setup-already-complete')
      return
    }
    
    passwordError.value = errorMessage
  } finally {
    isVerifying.value = false
  }
}

function nextStep() {
  if (currentStep.value === 1) {
    if (!canProceedStep1.value) {
      step1Error.value = 'Please fill all required fields and ensure passwords match (min 6 characters)'
      return
    }
    step1Error.value = ''
  }
  
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

async function completeSetup() {
  isCompleting.value = true
  setupError.value = ''

  try {
    // Create admin user (using Goli-Auth-Key during setup)
    await createUser({
      username: adminUser.value.username,
      password: adminUser.value.password,
      email: adminUser.value.email || '',
      role: 'admin'
    }, true) // useAuthKey = true

    // Update config (using Goli-Auth-Key during setup)
    const configData = {
      port: settings.value.port.toString(),
      auth_key: settings.value.authKey,
      setup_complete: true
    }
    await updateConfig(configData, true) // useAuthKey = true

    // Clear the authorization key from storage (no longer needed after setup)
    // User will use Bearer token after login
    clearAuthorizationKey()
    
    // Emit completion event
    emit('setup-complete')
  } catch (error) {
    setupError.value = error.message || 'Failed to complete setup. Please try again.'
    console.error('Setup error:', error)
  } finally {
    isCompleting.value = false
  }
}

// Load current config (only after setup password is verified and auth key is available)
async function loadConfig() {
  try {
    // Try to load config using auth key (if available from setup password verification)
    const config = await getConfig(true) // useAuthKey = true
    if (config.port) settings.value.port = parseInt(config.port) || 8125
    if (config.auth_key) settings.value.authKey = config.auth_key
    else generateNewAuthKey()
  } catch (error) {
    console.error('Error loading config:', error)
    // If we can't load config, generate a new auth key
    generateNewAuthKey()
  }
}
</script>

