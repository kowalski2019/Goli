import { ref, computed, watch, onMounted } from 'vue'

const THEME_STORAGE_KEY = 'goli_theme'
const THEMES = {
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system'
}

// Global state
const currentTheme = ref(THEMES.SYSTEM)
const systemTheme = ref('light')

// Get system theme preference
function getSystemTheme() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? THEMES.DARK : THEMES.LIGHT
}

// Computed dark mode state
const isDark = computed(() => {
  if (currentTheme.value === THEMES.SYSTEM) {
    return systemTheme.value === THEMES.DARK
  }
  return currentTheme.value === THEMES.DARK
})

// Apply theme to document
function applyTheme() {
  const root = document.documentElement
  root.classList.toggle('dark', isDark.value)
}

// Load theme from localStorage or use system default
function loadTheme() {
  const saved = localStorage.getItem(THEME_STORAGE_KEY)
  if (saved && Object.values(THEMES).includes(saved)) {
    currentTheme.value = saved
  } else {
    currentTheme.value = THEMES.SYSTEM
  }
  systemTheme.value = getSystemTheme()
  applyTheme()
}

// Save theme preference
function saveTheme(theme) {
  localStorage.setItem(THEME_STORAGE_KEY, theme)
  currentTheme.value = theme
  applyTheme()
}

// Watch for system theme changes
let systemThemeListener = null

function setupSystemThemeListener() {
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleChange = () => {
    systemTheme.value = getSystemTheme()
    if (currentTheme.value === THEMES.SYSTEM) {
      applyTheme()
    }
  }
  mediaQuery.addEventListener('change', handleChange)
  systemThemeListener = () => mediaQuery.removeEventListener('change', handleChange)
}

export function useTheme() {
  // Initialize theme on first use
  if (typeof window !== 'undefined' && !document.documentElement.classList.contains('dark') && !document.documentElement.classList.contains('light-initialized')) {
    loadTheme()
    setupSystemThemeListener()
    document.documentElement.classList.add('light-initialized')
  }

  // Watch for theme changes
  watch([currentTheme, systemTheme, isDark], () => {
    applyTheme()
  }, { immediate: true })

  return {
    currentTheme,
    isDark,
    themes: THEMES,
    setTheme: saveTheme,
    toggleTheme: () => {
      if (currentTheme.value === THEMES.LIGHT) {
        saveTheme(THEMES.DARK)
      } else if (currentTheme.value === THEMES.DARK) {
        saveTheme(THEMES.SYSTEM)
      } else {
        saveTheme(THEMES.LIGHT)
      }
    }
  }
}

