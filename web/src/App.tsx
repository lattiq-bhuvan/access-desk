import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { ProtectedLayout, PublicLayout } from '@/components/layout';
import { Catalog, Login, Requests } from '@/pages';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <PublicLayout>
              <Login />
            </PublicLayout>
          }
        />

        <Route
          path="/*"
          element={
            <ProtectedLayout>
              <Routes>
                <Route path="/" element={<Catalog />} />
                <Route path="/requests" element={<Requests />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </ProtectedLayout>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
