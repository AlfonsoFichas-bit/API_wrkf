import { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import { getProjects } from '../services/api';
import './Projects.css';

const Projects = () => {
  const [projects, setProjects] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const data = await getProjects();
        setProjects(data || []);
      } catch (err) {
        setError('Failed to fetch projects.');
        console.error(err);
      }
    };

    fetchProjects();
  }, []);

  return (
    <Layout>
      <div className="projects-container">
        <header className="projects-header">
          <h1>Mis Proyectos</h1>
          <button type="button" className="create-project-btn">
            Crear Nuevo Proyecto
          </button>
        </header>
        {error && <p className="error-message">{error}</p>}
        <div className="projects-list">
          {projects.length > 0 ? (
            projects.map((project) => (
              <div key={project.ID} className="project-card">
                <h2>{project.name}</h2>
                <p>{project.description}</p>
              </div>
            ))
          ) : (
            <p>
              No projects found. <a href="/projects/new">Create one!</a>
            </p>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Projects;
