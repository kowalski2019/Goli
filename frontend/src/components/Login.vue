<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-700 via-indigo-700 to-blue-700 p-4">
    <div class="w-full max-w-md bg-white rounded-2xl shadow-2xl p-8">
      <div class="flex flex-col items-center mb-6">
        <img src="/goli-logo.jpg" alt="Goli Logo" class="h-16 w-16 rounded-full shadow-md mb-3" />
        <h2 class="text-2xl font-bold text-gray-900">Welcome to Goli</h2>
        <p class="text-gray-500">Please sign in to continue</p>
      </div>
      <form @submit.prevent="handleSubmit" v-if="step==='login'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input v-model="username" type="text" required class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input v-model="password" type="password" required class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" />
        </div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
        <button type="submit" class="w-full py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
          Sign in
        </button>
      </form>

      <form @submit.prevent="handleVerify" v-else class="space-y-4">
        <div class="text-gray-700">
          We sent a verification code. Choose a channel and enter the code.
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Channel</label>
          <select v-model="channel" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500">
            <option v-for="c in channels" :key="c" :value="c">{{ c.toUpperCase() }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Code</label>
          <input v-model="code" type="tel" inputmode="numeric" pattern="[0-9]{6}" maxlength="6" class="w-full rounded-lg border-gray-300 focus:ring-primary-500 focus:border-primary-500" placeholder="000000" />
        </div>
        <div v-if="error" class="text-sm text-red-600">{{ error }}</div>
        <button type="submit" class="w-full py-2.5 rounded-lg bg-primary-600 text-white hover:bg-primary-700 transition">
          Verify
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { login, verify2FA, setToken } from '../api/client'

const emit = defineEmits(['logged-in'])

const step = ref('login')
const username = ref('')
const password = ref('')
const channel = ref('email')
const channels = ref(['email', 'sms'])
const code = ref('')
const error = ref('')

async function handleSubmit() {
  error.value = ''
  console.log('handleSubmit', username.value, password.value)
  try {
    const res = await login({ username: username.value, password: password.value })
    console.log('res', res)
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
  }
}

async function handleVerify() {
  error.value = ''
  // Trim whitespace from code
  const trimmedCode = code.value.trim()
  
  // Validate code format (6 digits)
  if (!/^\d{6}$/.test(trimmedCode)) {
    error.value = 'Please enter a 6-digit code'
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
  }
}
</script>


