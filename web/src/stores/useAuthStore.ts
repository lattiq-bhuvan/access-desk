import { createAuthStore } from '@lattiq/webtk';

import { APP_CONFIG } from '@/config/app.config';
import { buildApiUrl } from '@/lib/api/baseUrl';

// webtk hard-codes the auth paths it calls (/auth/v1/login, /auth/v1/refresh,
// /auth/v1/logout, /v1/users/me, /auth/v1/otp/*) — buildApiUrl only supplies
// the host. Our backend serves those routes directly (no /iam prefix like
// hub's separate IAM service), so no path rewriting is needed here.
export const useAuthStore = createAuthStore({
  buildApiUrl,
  access: APP_CONFIG.auth.access,
  storageKey: APP_CONFIG.auth.storageKey,
});

// Re-export types for convenience, same as hub does.
export type {
  AuthData,
  AuthError,
  AuthState,
  LoginResult,
  OTPResponse,
  User,
} from '@lattiq/webtk';
