import React from 'react';

const Header = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'dashboard', name: 'Dashboard'},
    { id: 'logs', name: 'Logs'},
    { id: 'autoheal', name: 'Auto-Heal'},
  ];

  return (
    <header className="bg-primary-900 bg-opacity-90 backdrop-blur-sm shadow-lg border-b border-primary-400 border-opacity-30">
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex">
            <div className="flex-shrink-0 flex items-center">
              {/* Nabd Logo */}
              <img 
                src="/nabd-logo.png" 
                alt="Nabd Logo" 
                className="h-20 w-auto mr-4"
              />
            </div>
            <nav className="hidden sm:ml-8 sm:flex sm:space-x-8">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`${
                    activeTab === tab.id
                      ? 'border-primary-400 text-text-primary bg-primary-400 bg-opacity-20'
                      : 'border-transparent text-text-secondary hover:text-text-primary hover:border-primary-400 hover:border-opacity-50'
                  } whitespace-nowrap py-3 px-4 border-b-3 font-medium text-lg flex items-center space-x-2 transition-all duration-200 rounded-t-lg`}
                >
                  <span>{tab.name}</span>
                </button>
              ))}
            </nav>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;