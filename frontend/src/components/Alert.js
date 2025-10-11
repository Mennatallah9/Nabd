import React from 'react';

const Alert = ({ alert, onClose }) => {
  const getSeverityColor = (severity) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-900 bg-opacity-30 border-red-400 text-red-200 backdrop-blur-sm';
      case 'warning':
        return 'bg-yellow-900 bg-opacity-30 border-yellow-400 text-yellow-200 backdrop-blur-sm';
      default:
        return 'bg-blue-900 bg-opacity-30 border-blue-400 text-blue-200 backdrop-blur-sm';
    }
  };

  const getSeverityIcon = (severity) => {
    switch (severity) {
      case 'critical':
        return (
          <div className="w-6 h-6 bg-red-500 rounded-full flex items-center justify-center">
            <span className="text-white text-sm">!</span>
          </div>
        );
      case 'warning':
        return (
          <div className="w-6 h-6 bg-yellow-500 rounded-full flex items-center justify-center">
            <span className="text-white text-sm">!</span>
          </div>
        );
      default:
        return (
          <div className="w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center">
            <span className="text-white text-sm">i</span>
          </div>
        );
    }
  };

  return (
    <div className={`border-l-4 p-6 rounded-lg ${getSeverityColor(alert.severity)}`}>
      <div className="flex justify-between items-start">
        <div className="flex">
          <div className="flex-shrink-0">
            {getSeverityIcon(alert.severity)}
          </div>
          <div className="ml-4">
            <h3 className="text-lg text-text-primary">
              {alert.name} - {alert.type.replace('_', ' ').toUpperCase()}
            </h3>
            <p className="text-base mt-2 text-text-secondary">{alert.message}</p>
            <p className="text-sm mt-2 text-text-secondary opacity-75">
              {new Date(alert.timestamp).toLocaleString()}
            </p>
          </div>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="text-text-secondary hover:text-text-primary transition-colors duration-200 p-1"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>
    </div>
  );
};

export default Alert;