import { NavLink } from 'react-router-dom';
import './Sidebar.css';

const Sidebar = () => {
  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <h2>AcademiaSys</h2>
      </div>
      <nav className="sidebar-nav">
        <ul>
          <li>
            <NavLink
              to="/projects"
              className={({ isActive }) => (isActive ? 'active' : '')}
            >
              Projects
            </NavLink>
          </li>
          <li>
            <NavLink
              to="/tasks"
              className={({ isActive }) => (isActive ? 'active' : '')}
            >
              Tasks
            </NavLink>
          </li>
          <li>
            <NavLink
              to="/evaluations"
              className={({ isActive }) => (isActive ? 'active' : '')}
            >
              Evaluations
            </NavLink>
          </li>
        </ul>
      </nav>
      <div className="sidebar-footer">
        {/* User info and logout button will go here */}
        <p>User Info</p>
        <button type="button">Logout</button>
      </div>
    </aside>
  );
};

export default Sidebar;
