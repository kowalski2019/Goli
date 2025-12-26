const API_BASE = '/api/v1'
const WS_BASE = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
const WS_URL = `${WS_BASE}//${window.location.host}/ws`

// Token management
export function setToken(token) {
  localStorage.setItem('goli_token', token)
}

export function getToken() {
  return localStorage.getItem('goli_token')
}

export function clearToken() {
  localStorage.removeItem('goli_token')
}

// Event name for automatic logout
export const AUTO_LOGOUT_EVENT = 'goli:auto-logout'

// Dispatch auto-logout event to notify App.vue
function triggerAutoLogout() {
  clearToken()
  window.dispatchEvent(new CustomEvent(AUTO_LOGOUT_EVENT))
}

// Wrapper for fetch that handles authentication errors
// Automatically logs out user if token is expired (401/403)
export async function fetchWithAuth(url, options = {}) {
  const response = await fetch(url, options)
  
  // Check for authentication errors (401 Unauthorized or 403 Forbidden)
  if (response.status === 401 || response.status === 403) {
    // Only trigger logout if we have a token (to avoid infinite loops)
    if (getToken()) {
      console.warn('Session expired. Logging out automatically.')
      triggerAutoLogout()
    }
    // Throw an error so the calling code can handle it
    const error = await response.json().catch(() => ({ 
      error: 'Authentication failed',
      description: 'Your session has expired. Please log in again.' 
    }))
    throw new Error(error.description || error.error || 'Authentication failed')
  }
  
  return response
}

// Authorization key management (for setup)
export function setAuthorizationKey(key) {
  localStorage.setItem('goli_authorization_key', key)
}

export function getAuthorizationKey() {
  return localStorage.getItem('goli_authorization_key')
}

export function clearAuthorizationKey() {
  localStorage.removeItem('goli_authorization_key')
}

// Get Bearer auth headers (for authenticated requests)
export function getBearerAuthHeaders() {
  const token = getToken()
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

// Get auth headers (legacy Goli-Auth-Key for setup)
export function getAuthKeyHeaders() {
  const authKey = getAuthorizationKey()
  return {
    'Content-Type': 'application/json',
    'Authorization': authKey ? `Goli-Auth-Key ${authKey}` : ''
  }
}

// Get form headers (for multipart/form-data)
export function getFormHeaders(useAuthKey = false) {
  if (useAuthKey) {
    const authKey = getAuthorizationKey()
    return {
      'Authorization': authKey ? `Goli-Auth-Key ${authKey}` : ''
    }
  }
  const token = getToken()
  return {
    'Authorization': token ? `Bearer ${token}` : ''
  }
}

// Setup API
export async function getSetupStatus() {
  const response = await fetch(`${API_BASE}/setup/status`, {
    method: 'GET'
  })
  if (!response.ok) throw new Error('Failed to fetch setup status')
  return response.json()
}

export async function verifySetupPassword(setupPassword) {
  const response = await fetch(`${API_BASE}/setup/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ setup_password: setupPassword })
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Invalid setup password')
  }
  const data = await response.json()
  // Store the auth key for setup completion
  if (data.auth_key) {
    setAuthorizationKey(data.auth_key)
  }
  return data
}

// Auth API
export async function login(credentials) {
  console.log('login', credentials)
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(credentials)
  })
  console.log('response', response)
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Login failed')
  }
  const data = await response.json()
  if (data.token) {
    setToken(data.token)
  }
  return data
}

export async function verify2FA(data) {
  const response = await fetch(`${API_BASE}/auth/2fa/verify`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || '2FA verification failed')
  }
  const result = await response.json()
  if (result.token) {
    setToken(result.token)
  }
  return result
}

export async function logout() {
  try {
    await fetchWithAuth(`${API_BASE}/auth/logout`, {
      method: 'POST',
      headers: getBearerAuthHeaders()
    })
  } catch (error) {
    console.error('Logout request failed:', error)
  } finally {
    clearToken()
  }
}

// Config API
export async function getConfig(useAuthKey = false) {
  const headers = useAuthKey ? getAuthKeyHeaders() : getBearerAuthHeaders()
  const response = useAuthKey 
    ? await fetch(`${API_BASE}/config`, { headers })
    : await fetchWithAuth(`${API_BASE}/config`, { headers })
  if (!response.ok) throw new Error('Failed to fetch config')
  return response.json()
}

export async function updateConfig(configData, useAuthKey = false) {
  const headers = useAuthKey ? getAuthKeyHeaders() : getBearerAuthHeaders()
  const response = useAuthKey
    ? await fetch(`${API_BASE}/config`, {
        method: 'POST',
        headers,
        body: JSON.stringify(configData)
      })
    : await fetchWithAuth(`${API_BASE}/config`, {
    method: 'POST',
    headers,
    body: JSON.stringify(configData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to update config')
  }
  return response.json()
}

// Users API
export async function getUsers() {
  const response = await fetchWithAuth(`${API_BASE}/users`, {
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch users')
  return response.json()
}

export async function createUser(userData, useAuthKey = false) {
  const headers = useAuthKey ? getAuthKeyHeaders() : getBearerAuthHeaders()
  const response = useAuthKey
    ? await fetch(`${API_BASE}/users`, {
        method: 'POST',
        headers,
        body: JSON.stringify(userData)
      })
    : await fetchWithAuth(`${API_BASE}/users`, {
    method: 'POST',
    headers,
    body: JSON.stringify(userData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to create user')
  }
  return response.json()
}

export async function updateUser(id, userData) {
  const response = await fetchWithAuth(`${API_BASE}/users/${id}`, {
    method: 'PUT',
    headers: getBearerAuthHeaders(),
    body: JSON.stringify(userData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to update user')
  }
  return response.json()
}

export async function deleteUser(id) {
  const response = await fetchWithAuth(`${API_BASE}/users/${id}`, {
    method: 'DELETE',
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to delete user')
  }
  return response.json()
}

// Jobs API
export async function getJobs(params = {}) {
  const query = new URLSearchParams(params).toString()
  const response = await fetchWithAuth(`${API_BASE}/jobs${query ? '?' + query : ''}`, {
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch jobs')
  const data = await response.json()
  // Ensure we always return an array, even if the API returns null
  return Array.isArray(data) ? data : []
}

export async function getJob(id) {
  const response = await fetchWithAuth(`${API_BASE}/jobs/${id}`, {
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch job')
  const data = await response.json()
  // Ensure steps is always an array
  if (data && !Array.isArray(data.steps)) {
    data.steps = []
  }
  return data
}

export async function createJob(jobData) {
  const response = await fetchWithAuth(`${API_BASE}/jobs`, {
    method: 'POST',
    headers: getBearerAuthHeaders(),
    body: JSON.stringify(jobData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to create job')
  }
  return response.json()
}

export async function cancelJob(id) {
  const response = await fetchWithAuth(`${API_BASE}/jobs/${id}/cancel`, {
    method: 'POST',
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to cancel job')
  }
  return response.json()
}

// Pipelines API
export async function getPipelines() {
  const response = await fetchWithAuth(`${API_BASE}/pipelines`, {
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch pipelines')
  const data = await response.json()
  // Ensure we always return an array, even if the API returns null
  return Array.isArray(data) ? data : []
}

export async function getPipeline(id) {
  const response = await fetchWithAuth(`${API_BASE}/pipelines/${id}`, {
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch pipeline')
  return response.json()
}

export async function createPipeline(pipelineData) {
  const response = await fetchWithAuth(`${API_BASE}/pipelines`, {
    method: 'POST',
    headers: getBearerAuthHeaders(),
    body: JSON.stringify(pipelineData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || error.error || 'Failed to create pipeline')
  }
  return response.json()
}

export async function updatePipeline(id, pipelineData) {
  const response = await fetchWithAuth(`${API_BASE}/pipelines/${id}`, {
    method: 'PUT',
    headers: getBearerAuthHeaders(),
    body: JSON.stringify(pipelineData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || error.error || 'Failed to update pipeline')
  }
  return response.json()
}

export async function deletePipeline(id) {
  const response = await fetchWithAuth(`${API_BASE}/pipelines/${id}`, {
    method: 'DELETE',
    headers: getBearerAuthHeaders()
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Failed to delete pipeline' }))
    throw new Error(error.error || error.message || 'Failed to delete pipeline')
  }
  return response.json()
}

export async function uploadPipeline(formData, useAuthKey = false) {
  const headers = getFormHeaders(useAuthKey)
  const response = useAuthKey
    ? await fetch(`${API_BASE}/pipelines/upload`, {
        method: 'POST',
        headers,
        body: formData
      })
    : await fetchWithAuth(`${API_BASE}/pipelines/upload`, {
    method: 'POST',
    headers,
    body: formData
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to upload pipeline')
  }
  return response.json()
}

export async function runPipeline(id, jobData = {}) {
  const response = await fetchWithAuth(`${API_BASE}/pipelines/${id}/run`, {
    method: 'POST',
    headers: getBearerAuthHeaders(),
    body: JSON.stringify(jobData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to run pipeline')
  }
  return response.json()
}

// WebSocket connection
export function createWebSocket(onMessage) {
  const ws = new WebSocket(WS_URL)
  
  ws.onopen = () => {
    console.log('WebSocket connected')
  }
  
  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data)
      onMessage(message)
    } catch (error) {
      console.error('Error parsing WebSocket message:', error)
    }
  }
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }
  
  ws.onclose = () => {
    console.log('WebSocket disconnected')
    // Auto-reconnect after 3 seconds
    setTimeout(() => {
      createWebSocket(onMessage)
    }, 3000)
  }
  
  return ws
}
