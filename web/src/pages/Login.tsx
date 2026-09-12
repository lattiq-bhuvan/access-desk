import { LoginPage } from '@lattiq/auth';
import { useNavigate } from 'react-router-dom';

import { APP_CONFIG } from '@/config/app.config';
import { useAuthStore } from '@/stores/useAuthStore';

export function Login() {
  const navigate = useNavigate();

  return (
    <LoginPage
      useAuthStore={useAuthStore}
      onLoginSuccess={() => navigate('/')}
      config={{
        platform: APP_CONFIG.platform.name,
        // PDF §2.3 M6: "Set primaryMethod: 'password' on the login page and
        // you can skip the OTP three" — our backend never returns a 202, so
        // OTP is neither implemented nor reachable.
        primaryMethod: 'password',
        allowToggle: false,
      }}
    />
  );
}
