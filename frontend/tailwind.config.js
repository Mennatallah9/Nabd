/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#d9d9d9',
          400: '#6895fd',
          500: '#6895fd',
          600: '#4e7dfc',
          700: '#3464fb',
          800: '#1a4bfa',
          900: '#001f65',
        },
        dark: {
          800: '#001f65',
          900: '#001330',
        },
        text: {
          primary: '#d9d9d9',
          secondary: '#b8b8b8',
        }
      },
      backgroundImage: {
        'gradient-primary': 'linear-gradient(135deg, #6895fd 0%, #001f65 100%)',
        'gradient-primary-dark': 'linear-gradient(135deg, #4e7dfc 0%, #001330 100%)',
      },
      fontSize: {
        'xs': ['0.6rem', { lineHeight: '0.8rem' }],
        'sm': ['0.7rem', { lineHeight: '1rem' }],
        'base': ['0.8rem', { lineHeight: '1.25rem' }],
        'lg': ['0.9rem', { lineHeight: '1.5rem' }],
        'xl': ['1rem', { lineHeight: '1.625rem' }],
        '2xl': ['1.1rem', { lineHeight: '1.75rem' }],
        '3xl': ['1.2rem', { lineHeight: '2rem' }],
        '4xl': ['1.3rem', { lineHeight: '2.25rem' }],
      }
    },
  },
  plugins: [],
}