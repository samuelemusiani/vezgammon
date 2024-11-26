import { ref, onMounted } from 'vue'

export function useTheme() {
  // Custom theme options based on Tailwind config
  const themeOptions = ['retro', 'retroPastel', 'retroDark', 'light', 'dark']

  // Current theme state
  const currentTheme = ref('retro')

  // Function to change theme
  function changeTheme(theme: string) {
    if (themeOptions.includes(theme)) {
      currentTheme.value = theme

      // Use DaisyUI's theme change method
      document.documentElement.setAttribute('data-theme', theme)

      // Optional: Persist theme in localStorage
      localStorage.setItem('app-theme', theme)
    }
  }

  // On mount, check for saved theme or set default
  onMounted(() => {
    const savedTheme = localStorage.getItem('app-theme')
    if (savedTheme && themeOptions.includes(savedTheme)) {
      changeTheme(savedTheme)
    } else {
      changeTheme('retro')  // Default theme
    }
  })

  return {
    currentTheme,
    themeOptions,
    changeTheme
  }
}
