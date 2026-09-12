/**
 * Single source of truth for platform branding and auth configuration.
 * Mirrors the pattern used in `hub` (lattiq/hub/src/config/app.config.ts).
 */
export const APP_CONFIG = {
  platform: {
    name: 'Accessdesk',
    description: 'Dataset access request and approval desk',
  },

  auth: {
    /** Access identifier sent to the server on login (webtk's AuthConfig.access) */
    access: 'portal:accessdesk',
    /** localStorage key the auth store persists to */
    storageKey: 'accessdesk-auth',
  },
} as const;

export type AppConfig = typeof APP_CONFIG;
