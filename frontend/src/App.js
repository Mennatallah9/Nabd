import React, { useState, useEffect } from 'react';
import Layout from './components/Layout';
import Header from './components/Header';
import Dashboard from './pages/Dashboard';
import Logs from './pages/Logs';
import AutoHeal from './pages/AutoHeal';
import Login from './pages/Login';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [activeTab, setActiveTab] = useState('dashboard');
  const [selectedContainer, setSelectedContainer] = useState(null);

  useEffect(() => {
    // Check if user is already authenticated
    const token = localStorage.getItem('nabd_token');
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const handleLogin = () => {
    setIsAuthenticated(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('nabd_token');
    setIsAuthenticated(false);
  };

  const handleViewLogs = (containerName) => {
    setSelectedContainer(containerName);
    setActiveTab('logs');
  };

  const handleBackToDashboard = () => {
    setSelectedContainer(null);
    setActiveTab('dashboard');
  };

  if (!isAuthenticated) {
    return <Login onLogin={handleLogin} />;
  }

  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return <Dashboard onViewLogs={handleViewLogs} />;
      case 'logs':
        return (
          <Logs 
            selectedContainer={selectedContainer} 
            onClose={selectedContainer ? handleBackToDashboard : null}
          />
        );
      case 'autoheal':
        return <AutoHeal />;
      default:
        return <Dashboard onViewLogs={handleViewLogs} />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <Header activeTab={activeTab} setActiveTab={setActiveTab} />
      <Layout>
        {renderContent()}
      </Layout>
    </div>
  );
}

export default App;