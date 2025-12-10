import { createBrowserRouter, Navigate } from 'react-router'

import LoginPage from '@/views/auth/Login'
import RegisterPage from '@/views/auth/Register'
import HomePage from '@/views/Home'
import App from '@/App'

export const router = createBrowserRouter([
  {
    element: <App />,
    children: [
      {
        path: '/',
        element: <HomePage />,
      },
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
])

export default router
