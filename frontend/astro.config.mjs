import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'astro/config';

import react from '@astrojs/react';

export default defineConfig({
  output: 'static',

  vite: {
    plugins: [tailwindcss()],
    server: {
      proxy: {
        '/api': 'http://192.168.3.16:8000'
      }
    }
  },

  integrations: [react()],
});