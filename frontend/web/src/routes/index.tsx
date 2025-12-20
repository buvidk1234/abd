import { createBrowserRouter, Navigate } from 'react-router'

import LoginPage from '@/views/auth/Login'
import RegisterPage from '@/views/auth/Register'
import { ChatPage } from '@/views/Home/pages/ChatPage'
import { ContactPage } from '@/views/Home/pages/ContactPage'
import App from '@/App'
import HomeLayout from '@/views/Home'

export const router = createBrowserRouter([
  {
    element: <App />,
    children: [
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
      {
        path: '/',
        element: <HomeLayout />,
        children: [
          {
            index: true,
            element: <Navigate to="/chat" replace />,
          },
          {
            path: 'chat',
            element: <ChatPage />,
          },
          {
            path: 'chat/:id',
            element: <ChatPage />,
          },
          {
            path: 'contact',
            element: <ContactPage />,
          },
          {
            path: 'contact/:id',
            element: <ContactPage />,
          },
        ],
      },
      { path: '*', element: <Navigate to="/" replace /> },
    ],
  },
])

export default router
