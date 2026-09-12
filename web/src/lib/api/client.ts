import { createApiClient } from '@lattiq/webtk';

import { API_BASE_URL } from '@/lib/api/baseUrl';
import { useAuthStore } from '@/stores/useAuthStore';

// Authenticated client: attaches the bearer token from useAuthStore to every
// request and transparently refreshes + retries once on a 401.
export const api = createApiClient({
  baseUrl: API_BASE_URL,
  authStore: useAuthStore,
});

export const apiUrl = api.buildUrl;
