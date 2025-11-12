<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-700 via-indigo-700 to-blue-700 p-4">
    <div class="w-full max-w-md bg-white rounded-2xl shadow-2xl overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-primary-600 to-primary-700 px-8 py-6">
        <div class="flex flex-col items-center">
          <div class="mb-4">
            <img 
              src="/goli-logo.jpg" 
              alt="Goli Logo" 
              class="h-16 w-16 rounded-full shadow-lg border-4 border-white/20"
              onerror="this.style.display='none'"
            />
          </div>
          <h2 class="text-2xl font-bold text-white mb-1">Welcome to Goli</h2>
          <p class="text-sm text-primary-100">Please sign in to continue</p>
        </div>
      </div>

      <!-- Form Content -->
      <div class="px-8 py-6">
        <!-- Error Alert -->
        <Alert
          v-if="error"
          type="error"
          :message="error"
          dismissible
          @dismiss="error = ''"
          class="mb-4"
        />

        <!-- Login Form -->
        <form @submit.prevent="handleSubmit" v-if="step === 'login'" class="space-y-5">
          <FormField label="Username" required>
            <TextInput
              v-model="username"
              type="text"
              placeholder="Enter your username"
              required
              :disabled="loading"
            />
          </FormField>

          <FormField label="Password" required>
            <TextInput
              v-model="password"
              type="password"
              placeholder="Enter your password"
              required
              :disabled="loading"
            />
          </FormField>

          <button
            type="submit"
            :disabled="loading"
            class="w-full btn btn-primary py-3 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg
              v-if="loading"
              class="animate-spin h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ loading ? 'Signing in...' : 'Sign in' }}</span>
          </button>
        </form>

        <!-- 2FA Verification Form -->
        <form @submit.prevent="handleVerify" v-else class="space-y-5">
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
              </svg>
              <div class="text-sm text-blue-800">
                <p class="font-medium mb-1">Verification Code Required</p>
                <p class="text-blue-700">We sent a verification code. Choose a channel and enter the code.</p>
              </div>
            </div>
          </div>

          <FormField label="Channel" required>
            <select
              v-model="channel"
              :disabled="loading"
              class="w-full px-4 py-2.5 text-sm border border-gray-300 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white disabled:bg-gray-50 disabled:cursor-not-allowed"
            >
              <option v-for="c in channels" :key="c" :value="c">
                {{ c === 'email' ? '📧 Email' : c === 'sms' ? '📱 SMS' : c.toUpperCase() }}
              </option>
            </select>
          </FormField>

          <FormField label="Verification Code" required description="Enter the 6-digit code">
            <input
              v-model="code"
              type="tel"
              inputmode="numeric"
              pattern="[0-9]{6}"
              maxlength="6"
              placeholder="000000"
              required
              :disabled="loading"
              class="w-full px-4 py-2.5 text-center text-lg tracking-widest font-mono border border-gray-300 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white disabled:bg-gray-50 disabled:cursor-not-allowed"
            />
          </FormField>

          <button
            type="submit"
            :disabled="loading || code.length !== 6"
            class="w-full btn btn-primary py-3 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg
              v-if="loading"
              class="animate-spin h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ loading ? 'Verifying...' : 'Verify' }}</span>
          </button>

          <button
            type="button"
            @click="step = 'login'; error = ''"
            class="w-full text-sm text-gray-600 hover:text-gray-800 py-2"
          >
            ← Back to login
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { login, verify2FA, setToken } from '../api/client'
import TextInput from './TextInput.vue'
import FormField from './FormField.vue'
import Alert from './Alert.vue'

const emit = defineEmits(['logged-in'])

const step = ref('login')
const username = ref('')
const password = ref('')
const channel = ref('email')
const channels = ref(['email', 'sms'])
const code = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  
  try {
    const res = await login({ username: username.value, password: password.value })
    
    if (res.two_fa_required) {
      channels.value = res.channels || ['email', 'sms']
      step.value = '2fa'
      return
    }
    
    if (res.token) {
      setToken(res.token)
      emit('logged-in', res.user)
    } else {
      error.value = 'Unexpected server response'
    }
  } catch (e) {
    error.value = e.message || 'Login failed'
  } finally {
    loading.value = false
  }
}

async function handleVerify() {
  error.value = ''
  loading.value = true
  
  // Trim whitespace from code
  const trimmedCode = code.value.trim()
  
  // Validate code format (6 digits)
  if (!/^\d{6}$/.test(trimmedCode)) {
    error.value = 'Please enter a 6-digit code'
    loading.value = false
    return
  }
  
  try {
    const res = await verify2FA({ username: username.value, channel: channel.value, code: trimmedCode })
    
    if (res.token) {
      setToken(res.token)
      emit('logged-in', res.user)
    } else {
      error.value = 'Invalid code'
    }
  } catch (e) {
    error.value = e.message || 'Verification failed'
  } finally {
    loading.value = false
  }
}
</script>
