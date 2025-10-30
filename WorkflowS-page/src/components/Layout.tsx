import { NavLink } from 'react-router-dom';

const Sidebar = () => {
  // TODO: Add real user info and logout functionality
  const handleLogout = () => {
    // localStorage.removeItem('token');
    // window.location.href = '/login';
    console.log('Logout clicked');
  };

  return (
    <aside className="flex w-64 flex-col border-r bg-white">
      <div className="flex h-16 shrink-0 items-center border-b px-6">
        <h2 className="text-xl font-bold">AcademiaSys</h2>
      </div>
      <nav className="flex-1 space-y-1 p-4">
        <NavLink
          to="/projects"
          className={({ isActive }) =>
            `flex items-center rounded-lg px-4 py-2 text-gray-700 hover:bg-gray-100 ${isActive ? 'bg-gray-200' : ''}`
          }
        >
          Projects
        </NavLink>
        <NavLink
          to="/tasks"
          className={({ isActive }) =>
            `flex items-center rounded-lg px-4 py-2 text-gray-700 hover:bg-gray-100 ${isActive ? 'bg-gray-200' : ''}`
          }
        >
          Tasks
        </NavLink>
        <NavLink
          to="/evaluations"
          className={({ isActive }) =>
            `flex items-center rounded-lg px-4 py-2 text-gray-700 hover:bg-gray-100 ${isActive ? 'bg-gray-200' : ''}`
          }
        >
          Evaluations
        </NavLink>
      </nav>
      <div className="mt-auto border-t p-4">
        <div className="flex items-center">
          {/* Replace with actual user info */}
          <div className="ml-3">
            <p className="text-sm font-medium text-gray-700">User Name</p>
            <p className="text-xs text-gray-500">user.email@example.com</p>
          </div>
        </div>
        <button
          type="button"
          onClick={handleLogout}
          className="mt-4 w-full rounded-lg bg-red-500 px-4 py-2 text-sm font-medium text-white hover:bg-red-600"
        >
          Logout
        </button>
      </div>
    </aside>
  );
};

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="flex h-screen bg-gray-50">
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-8">{children}</main>
    </div>
  );
};

export default Layout;
