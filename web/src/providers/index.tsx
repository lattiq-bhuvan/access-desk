import { TokenRefreshProvider } from '@lattiq/auth';
import { PageErrorBoundary, ThemeProvider, Toaster } from '@lattiq/design-system';
import type { ReactNode } from 'react';

import { useAuthStore } from '@/stores/useAuthStore';

interface ProvidersProps {
  children: ReactNode;
}

// Wraps the whole app: dark-mode (theme-provider + persisted preference),
// background token refresh so an expired access token is renewed
// transparently instead of surfacing a 401, an error boundary, and the
// toast host used for the "request submitted" success toast (§2.4).
export function Providers({ children }: ProvidersProps) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <TokenRefreshProvider useAuthStore={useAuthStore} />
      <PageErrorBoundary>{children}</PageErrorBoundary>
      <Toaster richColors position="bottom-right" />
    </ThemeProvider>
  );
}
