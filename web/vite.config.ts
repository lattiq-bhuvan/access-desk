import path from 'path';

import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, './src'),
    },
    // Dedupe so there's exactly one React/Zustand instance even though
    // @lattiq/auth and @lattiq/design-system both declare them as peers.
    dedupe: ['react', 'react-dom', 'zustand'],
  },
  server: {
    port: 5173,
  },
});
