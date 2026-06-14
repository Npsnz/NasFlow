/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#f5f3ff',
          100: '#ede9fe',
          200: '#ddd6fe',
          300: '#c084fc',
          400: '#b666ff',
          500: '#aa3bff',
          600: '#9333ea',
          700: '#7e22ce',
        },
        surface: {
          light: '#f8f9fa',
          DEFAULT: '#ffffff',
          dark: '#1e293b',
        },
        bg: {
          light: '#ffffff',
          dark: '#0f172a',
        }
      },
      backgroundColor: {
        'surface-light': '#f8f9fa',
        'surface-dark': '#1e293b',
      },
      fontFamily: {
        sans: ['Inter', 'Noto Sans Thai', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'sans-serif'],
      },
      keyframes: {
        'fade-up': {
          from: { transform: 'translateY(6px)', opacity: '0' },
          to: { transform: 'translateY(0)', opacity: '1' },
        },
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
      },
      animation: {
        'fade-up': 'fade-up 0.18s ease-out',
        'fade-in': 'fade-in 0.15s ease-out',
      },
    },
  },
  plugins: [],
}
