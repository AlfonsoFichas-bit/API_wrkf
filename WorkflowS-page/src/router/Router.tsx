import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Login from '../pages/Login';
import Projects from '../pages/Projects';
import Register from '../pages/Register';
import ProtectedRoute from './ProtectedRoute';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/register',
    element: <Register />,
  },
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        path: 'projects',
        element: <Projects />,
      },
      // Other protected routes go here
    ],
  },
]);

const AppRouter = () => {
  return <RouterProvider router={router} />;
};

export default AppRouter;
