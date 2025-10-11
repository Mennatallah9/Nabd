import React, { useState, useEffect } from 'react';
import { containerAPI } from '../services/api';

const Logs = ({ selectedContainer, onClose }) => {
  const [containers, setContainers] = useState([]);
  const [currentContainer, setCurrentContainer] = useState(selectedContainer || '');
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [lines, setLines] = useState(100);
  const [autoRefresh, setAutoRefresh] = useState(false);

  useEffect(() => {
    fetchContainers();
  }, []);

  useEffect(() => {
    if (selectedContainer) {
      setCurrentContainer(selectedContainer);
      fetchLogs(selectedContainer);
    }
  }, [selectedContainer]);

  useEffect(() => {
    let interval;
    if (autoRefresh && currentContainer) {
      interval = setInterval(() => {
        fetchLogs(currentContainer);
      }, 5000);
    }
    return () => clearInterval(interval);
  }, [autoRefresh, currentContainer]);

  const fetchContainers = async () => {
    try {
      const response = await containerAPI.getContainers();
      setContainers(response.data.data.filter(c => c.state === 'running'));
    } catch (err) {
      setError('Failed to fetch containers');
    }
  };

  const fetchLogs = async (containerName) => {
    if (!containerName) return;
    
    setLoading(true);
    try {
      const response = await containerAPI.getLogs(containerName, lines);
      setLogs(response.data.data);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch logs');
      setLogs([]);
    } finally {
      setLoading(false);
    }
  };

  const handleContainerChange = (containerName) => {
    setCurrentContainer(containerName);
    fetchLogs(containerName);
  };

  return (
    <div className="py-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-4xl text-text-primary">
            Container Logs
          </h1>
          <p className="text-xl text-text-secondary mt-2">
            View real-time logs from your containers
          </p>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="px-6 py-3 bg-primary-500 text-white text-lg font-medium rounded-lg hover:bg-primary-600 transition-all duration-200 shadow-lg"
          >
            Back to Dashboard
          </button>
        )}
      </div>

      {/* Controls */}
      <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 mb-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div>
            <label className="block text-lg font-medium text-text-primary mb-3">
              Container
            </label>
            <select
              value={currentContainer}
              onChange={(e) => handleContainerChange(e.target.value)}
              className="w-full px-4 py-3 border border-primary-400 border-opacity-50 rounded-lg bg-primary-900 bg-opacity-30 text-text-primary backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-primary-400 text-base"
            >
              <option value="">Select a container</option>
              {containers.map((container) => (
                <option key={container.id} value={container.name}>
                  {container.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-lg font-medium text-text-primary mb-3">
              Lines
            </label>
            <select
              value={lines}
              onChange={(e) => setLines(parseInt(e.target.value))}
              className="w-full px-4 py-3 border border-primary-400 border-opacity-50 rounded-lg bg-primary-900 bg-opacity-30 text-text-primary backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-primary-400 text-base"
            >
              <option value={50}>50 lines</option>
              <option value={100}>100 lines</option>
              <option value={200}>200 lines</option>
              <option value={500}>500 lines</option>
            </select>
          </div>

          <div className="flex items-end">
            <button
              onClick={() => fetchLogs(currentContainer)}
              disabled={!currentContainer || loading}
              className="w-full px-6 py-3 bg-primary-500 text-white text-lg font-medium rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg"
            >
              {loading ? 'Loading...' : 'Refresh'}
            </button>
          </div>

          <div className="flex items-end">
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="mr-3 h-5 w-5 text-primary-500 focus:ring-primary-400 border-primary-400 border-opacity-50 rounded bg-primary-900 bg-opacity-30"
              />
              <span className="text-lg text-text-primary">
                Auto-refresh (5s)
              </span>
            </label>
          </div>
        </div>
      </div>

      {error && (
        <div className="mb-8 bg-red-900 bg-opacity-30 border border-red-400 text-red-200 px-6 py-4 rounded-lg backdrop-blur-sm">
          {error}
        </div>
      )}

      {/* Logs Display */}
      <div className="bg-primary-900 bg-opacity-60 rounded-lg shadow-xl overflow-hidden border border-primary-400 border-opacity-30 backdrop-blur-sm">
        <div className="bg-primary-800 bg-opacity-60 px-6 py-4 flex justify-between items-center border-b border-primary-400 border-opacity-30">
          <span className="text-green-400 font-mono text-lg">
            {currentContainer ? `${currentContainer} logs` : 'No container selected'}
          </span>
          {logs.length > 0 && (
            <span className="text-text-secondary text-base">
              {logs.length} lines
            </span>
          )}
        </div>
        <div className="p-6 h-96 overflow-y-auto font-mono text-base">
          {currentContainer ? (
            logs.length > 0 ? (
              <div className="space-y-1">
                {logs.map((log, index) => (
                  <div key={index} className="text-green-400 whitespace-pre-wrap">
                    {log}
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-text-secondary text-center py-12 text-lg">
                {loading ? 'Loading logs...' : 'No logs available'}
              </div>
            )
          ) : (
            <div className="text-text-secondary text-center py-12 text-lg">
              Select a container to view logs
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Logs;