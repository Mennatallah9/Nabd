import React, { useState, useEffect } from 'react';
import { autoHealAPI } from '../services/api';
import LoadingSpinner from '../components/LoadingSpinner';

const AutoHeal = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [triggering, setTriggering] = useState(false);

  const fetchEvents = async () => {
    try {
      const response = await autoHealAPI.getHistory();
      setEvents(response.data.data);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch auto-heal events');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
    const interval = setInterval(() => {
      fetchEvents();
    }, 15000); // Refresh every 15 seconds
    return () => clearInterval(interval);
  }, []);

  const handleTriggerAutoHeal = async () => {
    setTriggering(true);
    try {
      await autoHealAPI.trigger();
      // Refresh events after a short delay
      setTimeout(fetchEvents, 2000);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to trigger auto-heal');
    } finally {
      setTriggering(false);
    }
  };

  // const handleToggleContainer = async (containerId, currentEnabled) => {
  //   try {
  //     await autoHealAPI.updateContainerConfig(containerId, !currentEnabled);
  //     await fetchContainers();
  //   } catch (err) {
  //     setError(err.response?.data?.error || 'Failed to update container configuration');
  //   }
  // };

  const getStatusIcon = (success) => {
    return success ? (
      <div className="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center">
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
    ) : (
      <div className="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center">
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </div>
    );
  };

  const getStatusColor = (success) => {
    return success 
      ? 'text-green-400' 
      : 'text-red-400';
  };

  const truncateContainerId = (containerId) => {
    if (!containerId) return '';
    // Show first 12 characters of container ID (standard Docker short ID length)
    return containerId.substring(0, 12);
  };

  const cleanReason = (reason) => {
    if (!reason) return '';
    
    //remove common timestamp patterns from the reason
    return reason
      .replace(/,\s*Status:\s*.*?(\d+\s*(hours?|minutes?|seconds?|days?)\s*ago)/gi, '')
      .replace(/,\s*Status:\s*.*?ago/gi, '')
      .replace(/\(\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2}.*?\)/g, '')
      .replace(/\d{4}-\d{2}-\d{2}[T\s]\d{2}:\d{2}:\d{2}[.\d]*[Z]?/g, '')
      .replace(/\s+/g, ' ')
      .trim();
  };

  const formatTimestamp = (timestamp) => {
    if (!timestamp) return 'Unknown';
    
    try {
      const date = new Date(timestamp);
      if (isNaN(date.getTime())) {
        return 'Invalid date';
      }
      
      const now = new Date();
      const diffMs = now - date;
      const diffMinutes = Math.floor(diffMs / (1000 * 60));
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
      
      let relativeTime = '';
      if (diffMinutes < 1) {
        relativeTime = 'Just now';
      } else if (diffMinutes < 60) {
        relativeTime = `${diffMinutes} minute${diffMinutes === 1 ? '' : 's'} ago`;
      } else if (diffHours < 24) {
        relativeTime = `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;
      } else if (diffDays < 7) {
        relativeTime = `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
      } else {
        relativeTime = date.toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric'
        });
      }
      
      const fullDateTime = date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        timeZoneName: 'short'
      });
      
      return { relativeTime, fullDateTime };
    } catch (error) {
      console.error('Error formatting timestamp:', timestamp, error);
      return { relativeTime: 'Invalid date', fullDateTime: 'Invalid date' };
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 relative">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="py-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-4xl text-text-primary">
            Auto-Healing
          </h1>
          <p className="text-xl text-text-secondary mt-2">
            Monitor and manage automatic container healing actions
          </p>
        </div>
        <button
          onClick={handleTriggerAutoHeal}
          disabled={triggering}
          className="px-6 py-3 bg-primary-500 text-white text-lg font-medium rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg"
        >
          {triggering ? 'Triggering...' : 'Trigger Check'}
        </button>
      </div>

      {error && (
        <div className="mb-8 bg-red-900 bg-opacity-30 border border-red-400 text-red-200 px-6 py-4 rounded-lg backdrop-blur-sm">
          {error}
        </div>
      )}

      {/* Info Card */}
      <div className="bg-primary-900 bg-opacity-40 border border-primary-400 border-opacity-30 rounded-lg p-8 mb-8 backdrop-blur-sm">
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <div className="w-8 h-8 bg-primary-400 rounded-full flex items-center justify-center">
              <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div className="ml-4">
            <h3 className="text-xl text-text-primary">
              Auto-Healing Information
            </h3>
            <div className="mt-3 text-lg text-text-secondary">
              <p>
                Auto-healing automatically monitors containers and performs recovery actions:
              </p>
              <ul className="list-disc list-inside mt-3 space-y-2">
                <li>Detects stopped or unhealthy containers</li>
                <li>Automatically restarts failed containers</li>
                <li>Logs all recovery actions</li>
                <li>Runs health checks every 15 seconds</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Events List */}
      <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl backdrop-blur-sm border border-primary-400 border-opacity-30">
        <div className="px-8 py-6 border-b border-primary-400 border-opacity-30">
          <h2 className="text-2xl text-text-primary">
            Recent Auto-Heal Events
          </h2>
        </div>
        
        {events.length > 0 ? (
          <div className="divide-y divide-primary-400 divide-opacity-20">
            {events.map((event) => (
              <div key={event.id} className="p-8">
                <div className="flex items-start justify-between">
                  <div className="flex items-start space-x-4">
                    {getStatusIcon(event.success)}
                    <div>
                      <div className="flex items-center space-x-3">
                        <h3 className="text-lg text-text-primary">
                          {event.name}
                        </h3>
                        <span className="px-3 py-1 text-sm font-medium rounded-full bg-primary-400 bg-opacity-30 text-text-primary border border-primary-400 border-opacity-50">
                          {event.action}
                        </span>
                      </div>
                      <p className="text-base text-text-secondary mt-2" title={event.container_id}>
                        Container ID: {truncateContainerId(event.container_id)}
                      </p>
                      <p className="text-lg text-text-secondary mt-3">
                        {cleanReason(event.reason)}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className={`text-lg ${getStatusColor(event.success)}`}>
                      {event.success ? 'Success' : 'Failed'}
                    </p>
                    <p 
                      className="text-sm text-text-secondary mt-1" 
                      title={(() => {
                        const formatted = formatTimestamp(event.timestamp);
                        return typeof formatted === 'object' ? formatted.fullDateTime : formatted;
                      })()}
                    >
                      {(() => {
                        const formatted = formatTimestamp(event.timestamp);
                        return typeof formatted === 'object' ? formatted.relativeTime : formatted;
                      })()}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="p-16 text-center">
            <h3 className="mt-2 text-xl text-text-primary">
              No auto-heal events
            </h3>
            <p className="mt-3 text-lg text-text-secondary">
              All containers are healthy! Auto-heal events will appear here when actions are taken.
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default AutoHeal;