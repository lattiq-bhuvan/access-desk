import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';

import { type AuthState, useAuthStore } from '@/stores/useAuthStore';

interface PublicLayoutProps {
  children: ReactNode;
}

// Wraps unauthenticated pages (login). Bounces an already-authenticated
// user straight to the app instead of showing the login form again.
export function PublicLayout({ children }: PublicLayoutProps) {
  const isAuthenticated = useAuthStore((state: AuthState) => state.isAuthenticated);

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
