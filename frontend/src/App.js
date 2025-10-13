import React, { useState, useEffect } from 'react';
import { Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom';
import Layout from './components/Layout';
import Header from './components/Header';
import Dashboard from './pages/Dashboard';
import Logs from './pages/Logs';
import Charts from './pages/Charts';
import AutoHeal from './pages/AutoHeal';
import Login from './pages/Login';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [activeTab, setActiveTab] = useState('dashboard');
  const [selectedContainer, setSelectedContainer] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Check if user is already authenticated
    const token = localStorage.getItem('nabd_token');
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const handleLogin = () => {
    setIsAuthenticated(true);
    // Redirect to dashboard after login
    navigate('/dashboard');
  };

  const handleLogout = () => {
    localStorage.removeItem('nabd_token');
    setIsAuthenticated(false);
    navigate('/login');
  };

  const handleViewLogs = (containerName) => {
    setSelectedContainer(containerName);
    setActiveTab('logs');
    navigate('/logs');
  };

  const handleViewCharts = (containerName) => {
    setSelectedContainer(containerName);
    navigate(`/charts/${containerName}`);
  };

  const handleBackToDashboard = () => {
    setSelectedContainer(null);
    setActiveTab('dashboard');
    navigate('/dashboard');
  };

  // Protected Route component
  const ProtectedRoute = ({ children }) => {
    if (!isAuthenticated) {
      return <Navigate to="/login" replace />;
    }
    return children;
  };

  // Login Route component that redirects if already authenticated
  const LoginRoute = () => {
    if (isAuthenticated) {
      return <Navigate to="/dashboard" replace />;
    }
    return <Login onLogin={handleLogin} />;
  };

  return (
    <Routes>
      <Route path="/login" element={<LoginRoute />} />
      <Route 
        path="/dashboard" 
        element={
          <ProtectedRoute>
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
              <Header activeTab="dashboard" setActiveTab={setActiveTab} onLogout={handleLogout} />
              <Layout>
                <Dashboard onViewLogs={handleViewLogs} onViewCharts={handleViewCharts} />
              </Layout>
            </div>
          </ProtectedRoute>
        } 
      />
      <Route 
        path="/logs" 
        element={
          <ProtectedRoute>
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
              <Header activeTab="logs" setActiveTab={setActiveTab} onLogout={handleLogout} />
              <Layout>
                <Logs 
                  selectedContainer={selectedContainer} 
                  onClose={selectedContainer ? handleBackToDashboard : null}
                />
              </Layout>
            </div>
          </ProtectedRoute>
        } 
      />
      <Route 
        path="/charts/:containerName" 
        element={
          <ProtectedRoute>
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
              <Header activeTab="charts" setActiveTab={setActiveTab} onLogout={handleLogout} />
              <Layout>
                <Charts />
              </Layout>
            </div>
          </ProtectedRoute>
        } 
      />
      <Route 
        path="/autoheal" 
        element={
          <ProtectedRoute>
            <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
              <Header activeTab="autoheal" setActiveTab={setActiveTab} onLogout={handleLogout} />
              <Layout>
                <AutoHeal />
              </Layout>
            </div>
          </ProtectedRoute>
        } 
      />
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default App;