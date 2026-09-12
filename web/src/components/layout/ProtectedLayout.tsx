import { useLogout } from '@lattiq/auth';
import { AppLayout } from '@lattiq/design-system';
import type { ReactNode } from 'react';
import { Link, Navigate, useLocation, useNavigate } from 'react-router-dom';

import { APP_CONFIG } from '@/config/app.config';
import { type AuthState, useAuthStore } from '@/stores/useAuthStore';

interface ProtectedLayoutProps {
  children: ReactNode;
}

// AppLayout wants a LinkComponent so its nav items render through the host
// app's router instead of a hard <a>.
function RouterLink({ href, children }: { href: string; children: ReactNode }) {
  return <Link to={href}>{children}</Link>;
}

export function isApprover(auth: AuthState['auth']): boolean {
  return !!auth?.user?.roles?.includes('approver');
}

// Shell for authenticated pages: header, nav, user menu + logout (all from
// AppLayout), dark-mode toggle (built into AppLayout's header), and a
// redirect-to-login guard for anyone who isn't authenticated.
export function ProtectedLayout({ children }: ProtectedLayoutProps) {
  const isAuthenticated = useAuthStore((state: AuthState) => state.isAuthenticated);
  const auth = useAuthStore((state: AuthState) => state.auth);
  const location = useLocation();
  const navigate = useNavigate();
  const handleLogout = useLogout(useAuthStore, {
    onComplete: () => navigate('/login', { replace: true }),
  });

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const navigationItems = [
    { label: 'Catalog', href: '/', active: location.pathname === '/' },
    { label: 'Requests', href: '/requests', active: location.pathname === '/requests' },
  ];

  return (
    <AppLayout
      platformName={APP_CONFIG.platform.name}
      user={{
        name: auth?.user?.name || 'User',
        email: auth?.user?.email || '',
      }}
      LinkComponent={RouterLink}
      navigationItems={navigationItems}
      onLogout={handleLogout}
    >
      {children}
    </AppLayout>
  );
}
