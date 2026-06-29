// Zen Theme Configuration
// Warm earth tones + Cool minimal + Dark mode

export const zenTheme = {
  // Light Mode - Warm Earth
  light: {
    warmEarth: {
      primary: '#A89968', // Sepia
      secondary: '#8B7355', // Brown
      accent: '#D4A574', // Warm tan
      bg: '#FAF8F3', // Off-white
      surface: '#FFFFFF', // White
      border: '#E8DDD2', // Warm beige
      text: '#3D3D3D', // Dark gray
      textSecondary: '#8B8B8B', // Medium gray
      muted: '#D9CEC0', // Muted brown
    },
    coolMinimal: {
      primary: '#4A5568', // Slate
      secondary: '#718096', // Cool gray
      accent: '#63B3ED', // Sky blue
      bg: '#F7FAFC', // Cool white
      surface: '#FFFFFF',
      border: '#E2E8F0', // Cool border
      text: '#1A202C', // Dark slate
      textSecondary: '#718096',
      muted: '#CBD5E0', // Cool gray
    },
  },

  // Dark Mode
  dark: {
    warmEarth: {
      primary: '#A89968',
      secondary: '#D4A574',
      accent: '#F5DEB3', // Wheat
      bg: '#1A1410', // Dark brown
      surface: '#2D2420', // Darker brown
      border: '#3D3530', // Brown border
      text: '#F5F1E8', // Light cream
      textSecondary: '#C9C1B8', // Warm gray
      muted: '#5C5450', // Muted dark
    },
    coolMinimal: {
      primary: '#E2E8F0', // Light gray
      secondary: '#CBD5E0',
      accent: '#63B3ED',
      bg: '#0F1419', // Very dark blue
      surface: '#1A202C', // Dark slate
      border: '#2D3748', // Dark border
      text: '#F7FAFC', // Almost white
      textSecondary: '#A0AEC0', // Medium gray
      muted: '#4A5568', // Muted slate
    },
  },

  // Typography
  fonts: {
    serif: 'Georgia, "Times New Roman", serif',
    sans: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    modern: 'Poppins, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },

  // Spacing (8px base)
  spacing: {
    xs: '0.5rem', // 8px
    sm: '1rem', // 16px
    md: '1.5rem', // 24px
    lg: '2rem', // 32px
    xl: '3rem', // 48px
    '2xl': '4rem', // 64px
  },

  // Border radius
  radius: {
    none: '0',
    sm: '0.25rem', // 4px
    base: '0.5rem', // 8px
    md: '0.75rem', // 12px
    lg: '1rem', // 16px
    full: '9999px',
  },

  // Shadows (minimal)
  shadows: {
    none: 'none',
    sm: '0 1px 2px rgba(0, 0, 0, 0.04)',
    base: '0 2px 4px rgba(0, 0, 0, 0.06)',
    md: '0 4px 8px rgba(0, 0, 0, 0.08)',
  },
}

export type ThemeMode = 'warmEarth' | 'coolMinimal'
export type ColorMode = 'light' | 'dark'
