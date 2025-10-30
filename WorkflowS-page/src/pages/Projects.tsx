import { useEffect, useState } from 'react';
import Layout from '../components/Layout';
import { getProjects } from '../services/api';

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
      <div>
        <header className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Mis Proyectos</h1>
          <a
            href="/projects/new"
            className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            Crear Nuevo Proyecto
          </a>
        </header>

        {error && <p className="mt-4 text-sm text-red-600">{error}</p>}

        <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {projects.length > 0 ? (
            projects.map((project) => (
              <div
                key={project.ID}
                className="overflow-hidden rounded-lg bg-white shadow"
              >
                <div className="p-5">
                  <h3 className="text-lg font-medium text-gray-900">
                    {project.name}
                  </h3>
                  <p className="mt-2 text-sm text-gray-500">
                    {project.description}
                  </p>
                </div>
              </div>
            ))
          ) : (
            <p className="text-sm text-gray-500">
              No projects found.{' '}
              <a
                href="/projects/new"
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Create one!
              </a>
            </p>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Projects;
