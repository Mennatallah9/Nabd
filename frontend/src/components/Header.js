import React from 'react';
import { useNavigate } from 'react-router-dom';

const Header = ({ activeTab, setActiveTab, onLogout }) => {
  const navigate = useNavigate();
  
  const tabs = [
    { id: 'dashboard', name: 'Dashboard', path: '/dashboard'},
    { id: 'logs', name: 'Logs', path: '/logs'},
    { id: 'autoheal', name: 'Auto-Heal', path: '/autoheal'},
  ];

  const handleTabClick = (tab) => {
    setActiveTab(tab.id);
    navigate(tab.path);
  };

  const handleLogoClick = () => {
    setActiveTab('dashboard');
    navigate('/dashboard');
  };

  const handleLogout = () => {
    if (onLogout) {
      onLogout();
    }
  };

  return (
    <header className="bg-primary-900 bg-opacity-90 backdrop-blur-sm shadow-lg border-b border-primary-400 border-opacity-30">
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              {/* Nabd Logo - Clickable */}
              <button
                onClick={handleLogoClick}
                className="hover:opacity-80 transition-opacity duration-200 focus:outline-none"
              >
                <img 
                  src="/nabd-logo.png" 
                  alt="Nabd Logo" 
                  className="h-20 w-auto mr-4"
                />
              </button>
            </div>
            <nav className="hidden sm:ml-8 sm:flex sm:space-x-8">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => handleTabClick(tab)}
                  className={`${
                    activeTab === tab.id
                      ? 'border-primary-400 text-text-primary bg-primary-400 bg-opacity-20'
                      : 'border-transparent text-text-secondary hover:text-text-primary hover:border-primary-400 hover:border-opacity-50 hover:bg-primary-400 hover:bg-opacity-10'
                  } whitespace-nowrap py-3 px-4 border-b-3 font-medium text-lg flex items-center space-x-2 transition-all duration-200 rounded-t-lg`}
                >
                  <span>{tab.name}</span>
                </button>
              ))}
            </nav>
          </div>
          
          {/* Logout Button */}
          <div className="flex items-center">
            <button
              onClick={handleLogout}
              className="border-transparent text-text-secondary hover:text-red-400 hover:border-red-400 hover:border-opacity-50 hover:bg-red-400 hover:bg-opacity-10 whitespace-nowrap py-3 px-4 border-b-3 font-medium text-lg flex items-center space-x-2 transition-all duration-200 rounded-t-lg"
            >
              <span>Logout</span>
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;