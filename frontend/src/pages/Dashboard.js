import React, { useState, useEffect } from 'react';
import { containerAPI, alertAPI } from '../services/api';
import ContainerCard from '../components/ContainerCard';
import Alert from '../components/Alert';
import LoadingSpinner from '../components/LoadingSpinner';

const Dashboard = ({ onViewLogs, onViewCharts }) => {
  const [containers, setContainers] = useState([]);
  const [metrics, setMetrics] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchData = async () => {
    try {
      const [containersRes, metricsRes, alertsRes] = await Promise.all([
        containerAPI.getContainers(),
        containerAPI.getMetrics(),
        alertAPI.getAlerts(),
      ]);

      setContainers(containersRes.data.data || []);
      setMetrics(metricsRes.data.data || []);
      setAlerts(alertsRes.data.data || []);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch data');
      // Keep existing data if fetch fails, don't reset to null
      // If we don't have data yet, initialize with empty arrays
      setContainers(prev => prev || []);
      setMetrics(prev => prev || []);
      setAlerts(prev => prev || []);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const handleRestart = async (containerName) => {
    try {
      await containerAPI.restartContainer(containerName);
      // Refresh data after restart
      setTimeout(fetchData, 2000);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to restart container');
    }
  };

  const getMetricForContainer = (containerName) => {
    return metrics ? metrics.find(m => m.name === containerName) : null;
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
      <div className="mb-8">
        <h1 className="text-4xl text-text-primary">
          Container Dashboard
        </h1>
        <p className="text-xl text-text-secondary mt-2">
          Monitor your Docker containers in real-time
        </p>
      </div>

      {error && (
        <div className="mb-8 bg-red-900 bg-opacity-30 border border-red-400 text-red-200 px-6 py-4 rounded-lg backdrop-blur-sm">
          {error}
        </div>
      )}

      {/* Alerts Section */}
      {alerts && alerts.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl text-text-primary mb-4">
            Active Alerts
          </h2>
          <div className="space-y-4">
            {alerts.map((alert) => (
              <Alert key={alert.id} alert={alert} />
            ))}
          </div>
        </div>
      )}

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-12 h-12 bg-primary-400 rounded-lg flex items-center justify-center">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
            </div>
            <div className="ml-4">
              <p className="text-lg font-medium text-text-secondary">
                Total Containers
              </p>
              <p className="text-3xl text-text-primary">
                {containers ? containers.length : 0}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-12 h-12 bg-green-500 rounded-lg flex items-center justify-center">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <div className="ml-4">
              <p className="text-lg font-medium text-text-secondary">
                Running
              </p>
              <p className="text-3xl text-green-400">
                {containers ? containers.filter(c => c.state === 'running').length : 0}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-12 h-12 bg-red-500 rounded-lg flex items-center justify-center">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <div className="ml-4">
              <p className="text-lg font-medium text-text-secondary">
                Stopped
              </p>
              <p className="text-3xl text-red-400">
                {containers ? containers.filter(c => c.state === 'exited').length : 0}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-12 h-12 bg-yellow-500 rounded-lg flex items-center justify-center">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.664-.833-2.464 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
            </div>
            <div className="ml-4">
              <p className="text-lg font-medium text-text-secondary">
                Active Alerts
              </p>
              <p className="text-3xl text-yellow-400">
                {alerts ? alerts.length : 0}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Containers Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8">
        {containers && containers.map((container) => (
          <ContainerCard
            key={container.id}
            container={container}
            metric={getMetricForContainer(container.name)}
            onRestart={handleRestart}
            onViewLogs={onViewLogs}
            onViewCharts={onViewCharts}
          />
        ))}
      </div>

      {(!containers || containers.length === 0) && !loading && (
        <div className="text-center py-16">
          <h3 className="mt-2 text-xl text-text-primary">
            No containers found
          </h3>
          <p className="mt-3 text-lg text-text-secondary">
            Make sure Docker is running and you have containers to monitor.
          </p>
        </div>
      )}
    </div>
  );
};

export default Dashboard;