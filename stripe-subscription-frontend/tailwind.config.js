/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,jsx,ts}'],
  theme: {
    extend: {}
  },
  plugins: [require('tailwind-scrollbar')({ nocompatible: true })]
};
