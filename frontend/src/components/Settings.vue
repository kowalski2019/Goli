<template>
  <div class="space-y-6 pb-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold text-white mb-1">Settings</h2>
        <p class="text-gray-300 text-sm">Manage your account and system configuration</p>
      </div>
    </div>

    <!-- Alerts -->
    <Alert
      v-if="success"
      type="success"
      :message="success"
      dismissible
      @dismiss="success = ''"
    />
    <Alert
      v-if="error"
      type="error"
      :message="error"
      dismissible
      @dismiss="error = ''"
    />

    <!-- User Profile Section -->
    <div class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex items-center gap-3 mb-6">
        <div class="p-2 bg-primary-100 rounded-lg">
          <svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-gray-900">User Profile</h3>
          <p class="text-sm text-gray-500">Update your personal information</p>
        </div>
      </div>

      <div v-if="currentUser" class="space-y-5">
        <FormField label="Username" description="Your username cannot be changed">
          <TextInput
            :model-value="currentUser.username"
            disabled
            class="bg-gray-50"
          />
        </FormField>

        <FormField label="Email Address" required>
          <TextInput
            v-model="userForm.email"
            type="email"
            placeholder="your.email@example.com"
          />
        </FormField>

        <FormField label="Phone Number">
          <TextInput
            v-model="userForm.phone"
            type="tel"
            placeholder="+1 (555) 123-4567"
          />
        </FormField>

        <div class="flex justify-end pt-2">
          <button
            @click="updateProfile"
            :disabled="saving"
            class="btn btn-primary px-6 py-2.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
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
            {{ saving ? 'Saving...' : 'Update Profile' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Two-Factor Authentication Section -->
    <div class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex items-center gap-3 mb-6">
        <div class="p-2 bg-green-100 rounded-lg">
          <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-gray-900">Two-Factor Authentication</h3>
          <p class="text-sm text-gray-500">Add an extra layer of security to your account</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="flex items-center justify-between p-4 bg-gradient-to-r from-gray-50 to-gray-50/50 rounded-lg border border-gray-200 hover:border-gray-300 transition-colors">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-1">
              <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              <h4 class="font-semibold text-gray-900">Email 2FA</h4>
            </div>
            <p class="text-sm text-gray-600 ml-8">Receive verification codes via email</p>
          </div>
          <ToggleSwitch
            v-model="twoFA.email"
            @update:model-value="update2FA"
          />
        </div>

        <div class="flex items-center justify-between p-4 bg-gradient-to-r from-gray-50 to-gray-50/50 rounded-lg border border-gray-200 hover:border-gray-300 transition-colors">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-1">
              <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
              <h4 class="font-semibold text-gray-900">SMS 2FA</h4>
            </div>
            <p class="text-sm text-gray-600 ml-8">Receive verification codes via SMS</p>
          </div>
          <ToggleSwitch
            v-model="twoFA.sms"
            @update:model-value="update2FA"
          />
        </div>
      </div>
    </div>

    <!-- System Configuration Section (Admin only) -->
    <div v-if="currentUser && currentUser.role === 'admin'" class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex items-center gap-3 mb-6">
        <div class="p-2 bg-blue-100 rounded-lg">
          <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-gray-900">System Configuration</h3>
          <p class="text-sm text-gray-500">Configure system-wide settings</p>
        </div>
      </div>

      <div class="space-y-5">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <FormField label="Port" required description="Server port number">
            <TextInput
              :model-value="config.port"
              type="number"
              placeholder="8125"
              @update:model-value="config.port = parseInt($event) || 8125"
            />
          </FormField>

          <FormField label="GitHub Username">
            <TextInput
              v-model="config.gh_username"
              type="text"
              placeholder="your-github-username"
            />
          </FormField>
        </div>

        <FormField label="GitHub Access Token" description="Token for GitHub API access">
          <TextInput
            v-model="config.gh_access_token"
            type="password"
            placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxx"
          />
        </FormField>

        <div class="flex justify-end pt-2">
          <button
            @click="updateConfig"
            :disabled="savingConfig"
            class="btn btn-primary px-6 py-2.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <svg
              v-if="savingConfig"
              class="animate-spin h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ savingConfig ? 'Saving...' : 'Update Configuration' }}
          </button>
        </div>
      </div>
    </div>

    <!-- SMTP Configuration Section (Admin only) -->
    <div v-if="currentUser && currentUser.role === 'admin'" class="card hover:shadow-lg transition-shadow duration-200">
      <div class="flex items-center gap-3 mb-6">
        <div class="p-2 bg-purple-100 rounded-lg">
          <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-gray-900">SMTP Configuration</h3>
          <p class="text-sm text-gray-500">Configure email server settings for notifications</p>
        </div>
      </div>

      <div class="space-y-5">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <FormField label="SMTP Host" required description="Your SMTP server address">
            <TextInput
              v-model="config.smtp_host"
              type="text"
              placeholder="smtp.example.com"
            />
          </FormField>

          <FormField label="SMTP Port" required description="Usually 587 (TLS) or 465 (SSL)">
            <TextInput
              :model-value="config.smtp_port"
              type="number"
              placeholder="587"
              @update:model-value="config.smtp_port = parseInt($event) || 587"
            />
          </FormField>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <FormField label="SMTP Username" required description="Your email address or username">
            <TextInput
              v-model="config.smtp_user"
              type="text"
              placeholder="your-email@example.com"
            />
          </FormField>

          <FormField label="SMTP Password" required description="Your SMTP password or app password">
            <TextInput
              v-model="config.smtp_pass"
              type="password"
              placeholder="Your SMTP password"
            />
          </FormField>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <FormField label="From Email" required description="Sender email address">
            <TextInput
              v-model="config.smtp_from"
              type="email"
              placeholder="noreply@example.com"
            />
          </FormField>

          <FormField label="From Name" description="Display name for sent emails">
            <TextInput
              v-model="config.smtp_from_name"
              type="text"
              placeholder="Goli CI/CD"
            />
          </FormField>
        </div>

        <div class="flex justify-end pt-2">
          <button
            @click="updateConfig"
            :disabled="savingConfig"
            class="btn btn-primary px-6 py-2.5 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <svg
              v-if="savingConfig"
              class="animate-spin h-4 w-4"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ savingConfig ? 'Saving...' : 'Update SMTP Configuration' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUsers, updateUser, getConfig, updateConfig as updateConfigAPI } from '../api/client'
import TextInput from './TextInput.vue'
import FormField from './FormField.vue'
import Alert from './Alert.vue'
import ToggleSwitch from './ToggleSwitch.vue'

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
  gh_access_token: '',
  smtp_host: '',
  smtp_port: 587,
  smtp_user: '',
  smtp_pass: '',
  smtp_from: '',
  smtp_from_name: ''
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
    config.value.smtp_host = configData.smtp_host || ''
    config.value.smtp_port = parseInt(configData.smtp_port) || 587
    config.value.smtp_user = configData.smtp_user || ''
    config.value.smtp_pass = configData.smtp_pass || ''
    config.value.smtp_from = configData.smtp_from || ''
    config.value.smtp_from_name = configData.smtp_from_name || ''
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
    const updateData = {
      port: config.value.port.toString(),
      gh_username: config.value.gh_username,
      gh_access_token: config.value.gh_access_token,
      smtp_host: config.value.smtp_host || '',
      smtp_port: config.value.smtp_port ? config.value.smtp_port.toString() : '',
      smtp_user: config.value.smtp_user || '',
      smtp_pass: config.value.smtp_pass || '',
      smtp_from: config.value.smtp_from || '',
      smtp_from_name: config.value.smtp_from_name || ''
    }
    
    await updateConfigAPI(updateData)
    success.value = 'Configuration updated successfully'
    await loadConfig()
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
