import { useId, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { register } from '../services/api';
import './Register.css';

const Register = () => {
  const [formData, setFormData] = useState({
    nombre: '',
    apellidoPaterno: '',
    apellidoMaterno: '',
    correo: '',
    contraseña: '',
  });
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const nombreId = useId();
  const apellidoPaternoId = useId();
  const apellidoMaternoId = useId();
  const correoId = useId();
  const contraseñaId = useId();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prevState) => ({ ...prevState, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      await register(formData);
      navigate('/login'); // Redirect to login page on successful registration
    } catch (err) {
      setError('Failed to register. Please try again.');
      console.error(err);
    }
  };

  return (
    <div className="register-container">
      <div className="register-form">
        <h1>Crear Cuenta</h1>
        <p>Únete a AcademiaSys</p>
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label htmlFor={nombreId}>Nombre</label>
            <input
              type="text"
              id={nombreId}
              name="nombre"
              value={formData.nombre}
              onChange={handleChange}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor={apellidoPaternoId}>Apellido Paterno</label>
            <input
              type="text"
              id={apellidoPaternoId}
              name="apellidoPaterno"
              value={formData.apellidoPaterno}
              onChange={handleChange}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor={apellidoMaternoId}>Apellido Materno</label>
            <input
              type="text"
              id={apellidoMaternoId}
              name="apellidoMaterno"
              value={formData.apellidoMaterno}
              onChange={handleChange}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor={correoId}>Correo Electrónico</label>
            <input
              type="email"
              id={correoId}
              name="correo"
              value={formData.correo}
              onChange={handleChange}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor={contraseñaId}>Contraseña</label>
            <input
              type="password"
              id={contraseñaId}
              name="contraseña"
              value={formData.contraseña}
              onChange={handleChange}
              required
            />
          </div>
          {error && <p className="error-message">{error}</p>}
          <button type="submit" className="register-button">
            Registrarse
          </button>
        </form>
        <p className="login-link">
          ¿Ya tienes una cuenta? <a href="/login">Inicia Sesión</a>
        </p>
      </div>
      <div className="register-image">
        {/* You can place an image here as in the mockup */}
      </div>
    </div>
  );
};

export default Register;
