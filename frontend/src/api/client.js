const API_BASE = '/api/v1'
const WS_BASE = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
const WS_URL = `${WS_BASE}//${window.location.host}/ws`

// Get auth key from localStorage or prompt
function getAuthKey() {
  let authKey = localStorage.getItem('goli_auth_key')
  if (!authKey) {
    authKey = prompt('Enter your Goli Auth Key:') || 'dummy_key'
    localStorage.setItem('goli_auth_key', authKey)
  }
  return authKey
}

export function getAuthHeaders() {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Goli-Auth-Key ${getAuthKey()}`
  }
}

export function getFormHeaders() {
  return {
    'Authorization': `Goli-Auth-Key ${getAuthKey()}`
  }
}

// Jobs API
export async function getJobs(params = {}) {
  const query = new URLSearchParams(params).toString()
  const response = await fetch(`${API_BASE}/jobs${query ? '?' + query : ''}`, {
    headers: getAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch jobs')
  const data = await response.json()
  // Ensure we always return an array, even if the API returns null
  return Array.isArray(data) ? data : []
}

export async function getJob(id) {
  const response = await fetch(`${API_BASE}/jobs/${id}`, {
    headers: getAuthHeaders()
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
  const response = await fetch(`${API_BASE}/jobs`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(jobData)
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to create job')
  }
  return response.json()
}

// Pipelines API
export async function getPipelines() {
  const response = await fetch(`${API_BASE}/pipelines`, {
    headers: getAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch pipelines')
  const data = await response.json()
  // Ensure we always return an array, even if the API returns null
  return Array.isArray(data) ? data : []
}

export async function getPipeline(id) {
  const response = await fetch(`${API_BASE}/pipelines/${id}`, {
    headers: getAuthHeaders()
  })
  if (!response.ok) throw new Error('Failed to fetch pipeline')
  return response.json()
}

export async function uploadPipeline(formData) {
  const response = await fetch(`${API_BASE}/pipelines/upload`, {
    method: 'POST',
    headers: getFormHeaders(),
    body: formData
  })
  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.description || 'Failed to upload pipeline')
  }
  return response.json()
}

export async function runPipeline(id, jobData = {}) {
  const response = await fetch(`${API_BASE}/pipelines/${id}/run`, {
    method: 'POST',
    headers: getAuthHeaders(),
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

