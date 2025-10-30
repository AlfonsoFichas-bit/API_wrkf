import { useId, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/api';
import './Login.css';

const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const emailId = useId();
  const passwordId = useId();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const data = await login(email, password);
      if (data.token) {
        localStorage.setItem('token', data.token);
        navigate('/projects'); // Redirect to projects page on successful login
      }
    } catch (err) {
      setError('Invalid email or password');
      console.error(err);
    }
  };

  return (
    <div className="login-container">
      <div className="login-form">
        <h1>Iniciar Sesión</h1>
        <p>Bienvenido de vuelta a AcademiaSys</p>
        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <label htmlFor={emailId}>Correo Electrónico</label>
            <input
              type="email"
              id={emailId}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor={passwordId}>Contraseña</label>
            <input
              type="password"
              id={passwordId}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          {error && <p className="error-message">{error}</p>}
          <button type="submit" className="login-button">
            Iniciar Sesión
          </button>
        </form>
        <p className="signup-link">
          ¿No tienes una cuenta? <a href="/register">Regístrate</a>
        </p>
      </div>
      <div className="login-image">
        {/* You can place an image here as in the mockup */}
      </div>
    </div>
  );
};

export default Login;
