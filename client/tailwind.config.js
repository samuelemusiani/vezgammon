/** @type {import('tailwindcss').Config} */
import daisyui from 'daisyui'

export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'retro-pattern': `
          repeating-linear-gradient(45deg, currentColor 0px, currentColor 2px, transparent 2px, transparent 10px),
          repeating-linear-gradient(-45deg, currentColor 0px, currentColor 2px, transparent 2px, transparent 10px)
        `,
        },
      cursor: {
        'tortellino': 'url("/tortellino.png"), auto',
      },
    },
  },
  daisyui: {
    themes: [
      {
        retro: {
          primary: '#d2691e',
          'primary-focus': '#8b4513',
          'primary-content': '#ffffff',
          secondary: '#8b4513',
          accent: '#cd853f',
          neutral: '#2c1810',
          'base-100': '#ffe5c9',
          'base-200': '#f0e6d2',
          'base-300': '#e6d8c0',
          'base-content': '#3c1e10',
          background: '--var(retro-background-base)',
        },
        retroPastel: {
          primary: '#cd853f',
          'primary-focus': '#a0522d',
          'primary-content': '#ffffff',
          secondary: '#a0522d',
          accent: '#8b4513',
          neutral: '#f0e6d2',
          'base-100': '#fff5e6',
          'base-200': '#f5eee0',
          'base-300': '#ebe4d6',
          'base-content': '#5d4037',
        },
        retroDark: {
          primary: '#8b4513',
          'primary-focus': '#5d4037',
          'primary-content': '#ffffff',
          secondary: '#4a4a4a',
          accent: '#cd853f',
          neutral: '#1a1a1a',
          'base-100': '#2c2c2c',
          'base-200': '#262626',
          'base-300': '#1f1f1f',
          'base-content': '#d2d2d2',
        },
      },
      'light',
      'dark',
    ],
  },
  plugins: [daisyui],
}
