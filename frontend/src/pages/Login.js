import React, { useState } from 'react';
import { authAPI } from '../services/api';

const Login = ({ onLogin }) => {
  const [token, setToken] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await authAPI.login(token);
      const jwtToken = response.data.token;
      
      localStorage.setItem('nabd_token', jwtToken);
      onLogin();
    } catch (err) {
      setError(err.response?.data?.error || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-primary py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <div className="flex justify-center mb-6">
            <img 
              src="/nabd-logo.png" 
              alt="Nabd Logo" 
              className="h-24 w-auto"
            />
          </div>
          <h2 className="mt-6 text-center text-4xl text-text-primary">
            Nabd Dashboard
          </h2>
          <p className="mt-4 text-center text-xl text-text-secondary">
            Container Observability & Auto-Healing Tool
          </p>
          <p className="mt-6 text-center text-lg text-text-secondary">
            Enter your admin token to access the dashboard
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div>
            <label htmlFor="token" className="sr-only">
              Admin Token
            </label>
            <input
              id="token"
              name="token"
              type="password"
              required
              className="appearance-none rounded-lg relative block w-full px-4 py-3 border border-primary-400 border-opacity-50 placeholder-text-secondary text-text-primary bg-primary-900 bg-opacity-30 backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-primary-400 focus:border-primary-400 focus:z-10 text-lg"
              placeholder="Admin Token"
              value={token}
              onChange={(e) => setToken(e.target.value)}
            />
          </div>

          {error && (
            <div className="text-red-400 text-lg text-center bg-red-900 bg-opacity-30 border border-red-400 border-opacity-50 rounded-lg py-3 backdrop-blur-sm">
              {error}
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={loading}
              className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-lg rounded-lg text-white bg-primary-500 hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg"
            >
              {loading ? 'Authenticating...' : 'Sign in'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;