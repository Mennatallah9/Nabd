import React from 'react';

const LoadingSpinner = ({ className = "", size = "default" }) => {
  const isSmall = size === "small";
  const spinnerClass = isSmall ? "loading-small" : "loading";

  return (
    <div className={`${spinnerClass} ${className}`}>
      <svg width="16px" height="12px">
        <polyline 
          id="back"
          points="1 6 4 6 6 11 10 1 12 6 15 6"
          fill="none"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          stroke="rgba(216, 217, 217, 0.18)"
        />
        <polyline 
          id="front"
          className="loading-front" 
          points="1 6 4 6 6 11 10 1 12 6 15 6"
          fill="none"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          stroke="#d9d9d9"
          strokeDasharray="12, 36"
          strokeDashoffset="48"
        />
      </svg>
    </div>
  );
};

export default LoadingSpinner;