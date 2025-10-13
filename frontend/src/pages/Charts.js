import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { containerAPI } from '../services/api';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  AreaChart,
  Area
} from 'recharts';
import LoadingSpinner from '../components/LoadingSpinner';

const Charts = () => {
  const { containerName } = useParams();
  const navigate = useNavigate();
  const [metricsData, setMetricsData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [timeRange, setTimeRange] = useState(24); // hours
  const [container, setContainer] = useState(null);

  // Format bytes to human readable format
  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  // Format timestamp for chart
  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit',
      hour12: false 
    });
  };

  // Custom tooltip for charts
  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-primary-900 bg-opacity-90 p-4 rounded-lg border border-primary-400 border-opacity-50 backdrop-blur-sm">
          <p className="text-text-primary font-semibold">{`Time: ${new Date(label).toLocaleString()}`}</p>
          {payload.map((entry, index) => (
            <p key={index} className="text-text-secondary">
              <span style={{ color: entry.color }}>{entry.name}: </span>
              {entry.dataKey.includes('memory') ? formatBytes(entry.value) : 
               entry.dataKey.includes('network') ? formatBytes(entry.value) :
               entry.dataKey.includes('cpu') ? `${entry.value.toFixed(2)}%` : entry.value}
            </p>
          ))}
        </div>
      );
    }
    return null;
  };

  const fetchMetricsData = async () => {
    try {
      setLoading(true);
      //get container info
      const containersRes = await containerAPI.getContainers();
      const containers = containersRes.data.data || [];
      const currentContainer = containers.find(c => c.name === containerName);
      
      if (!currentContainer) {
        setError('Container not found');
        return;
      }
      
      setContainer(currentContainer);

      //get metrics history
      const metricsRes = await containerAPI.getMetricsHistory(currentContainer.id, timeRange);
      const metrics = metricsRes.data.data || [];
      
      // Process and format the data for charts
      const processedData = metrics.map(metric => ({
        timestamp: metric.timestamp,
        time: formatTime(metric.timestamp),
        cpu_percent: metric.cpu_percent,
        memory_usage: metric.memory_usage,
        memory_limit: metric.memory_limit,
        memory_percentage: metric.memory_limit ? (metric.memory_usage / metric.memory_limit) * 100 : 0,
        network_rx: metric.network_rx,
        network_tx: metric.network_tx
      })).sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

      setMetricsData(processedData);
      setError('');
    } catch (err) {
      console.error('Error fetching metrics:', err);
      setError(err.response?.data?.error || 'Failed to fetch metrics data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (containerName) {
      fetchMetricsData();
    }
  }, [containerName, timeRange]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 relative">
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <div className="py-8">
        <div className="mb-4">
          <button
            onClick={() => navigate('/')}
            className="px-4 py-2 bg-primary-500 bg-opacity-20 text-primary-200 border border-primary-400 border-opacity-50 rounded-lg hover:bg-opacity-30 transition-all duration-200"
          >
            ← Back to Dashboard
          </button>
        </div>
        <div className="bg-red-900 bg-opacity-30 border border-red-400 text-red-200 px-6 py-4 rounded-lg backdrop-blur-sm">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div>
            <button
              onClick={() => navigate('/')}
              className="px-4 py-2 bg-primary-500 bg-opacity-20 text-primary-200 border border-primary-400 border-opacity-50 rounded-lg hover:bg-opacity-30 transition-all duration-200 mb-4"
            >
              ← Back to Dashboard
            </button>
            <h1 className="text-4xl text-text-primary">
              Container Metrics
            </h1>
            <p className="text-xl text-text-secondary mt-2">
              {container ? `${container.name} (${container.image})` : containerName}
            </p>
          </div>
          
          {/* Time Range Selector */}
          <div className="flex items-center space-x-4">
            <label className="text-text-secondary">Time Range:</label>
            <select
              value={timeRange}
              onChange={(e) => setTimeRange(parseInt(e.target.value))}
              className="px-3 py-2 bg-primary-900 bg-opacity-40 border border-primary-400 border-opacity-30 rounded-lg text-text-primary focus:outline-none focus:border-primary-400"
            >
              <option value={1}>Last Hour</option>
              <option value={6}>Last 6 Hours</option>
              <option value={24}>Last 24 Hours</option>
              <option value={72}>Last 3 Days</option>
              <option value={168}>Last Week</option>
            </select>
            <button
              onClick={fetchMetricsData}
              className="px-4 py-2 bg-green-500 bg-opacity-20 text-green-200 border border-green-400 border-opacity-50 rounded-lg hover:bg-opacity-30 transition-all duration-200"
            >
              Refresh
            </button>
          </div>
        </div>
      </div>

      {metricsData.length === 0 ? (
        <div className="text-center py-16">
          <h3 className="mt-2 text-xl text-text-primary">
            No metrics data available
          </h3>
          <p className="mt-3 text-lg text-text-secondary">
            Metrics data will appear here once the container generates some activity.
          </p>
        </div>
      ) : (
        <div className="space-y-8">
          {/* CPU Usage Chart */}
          <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-6 border border-primary-400 border-opacity-30 backdrop-blur-sm">
            <h3 className="text-2xl text-text-primary mb-4">CPU Usage</h3>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={metricsData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="timestamp" 
                  tickFormatter={formatTime}
                  stroke="#9CA3AF"
                />
                <YAxis 
                  domain={[0, 'dataMax']}
                  stroke="#9CA3AF"
                  tickFormatter={(value) => `${value.toFixed(1)}%`}
                />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="cpu_percent"
                  stroke="#3B82F6"
                  fill="#3B82F6"
                  fillOpacity={0.3}
                  name="CPU Usage (%)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>

          {/* Memory Usage Chart */}
          <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-6 border border-primary-400 border-opacity-30 backdrop-blur-sm">
            <h3 className="text-2xl text-text-primary mb-4">Memory Usage</h3>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={metricsData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="timestamp" 
                  tickFormatter={formatTime}
                  stroke="#9CA3AF"
                />
                <YAxis 
                  stroke="#9CA3AF"
                  tickFormatter={(value) => formatBytes(value)}
                />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="memory_usage"
                  stroke="#10B981"
                  strokeWidth={2}
                  dot={false}
                  name="Memory Usage"
                />
                <Line
                  type="monotone"
                  dataKey="memory_limit"
                  stroke="#F59E0B"
                  strokeWidth={2}
                  strokeDasharray="5 5"
                  dot={false}
                  name="Memory Limit"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>

          {/* Memory Percentage Chart */}
          <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-6 border border-primary-400 border-opacity-30 backdrop-blur-sm">
            <h3 className="text-2xl text-text-primary mb-4">Memory Usage Percentage</h3>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={metricsData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="timestamp" 
                  tickFormatter={formatTime}
                  stroke="#9CA3AF"
                />
                <YAxis 
                  domain={[0, 100]}
                  stroke="#9CA3AF"
                  tickFormatter={(value) => `${value.toFixed(1)}%`}
                />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Area
                  type="monotone"
                  dataKey="memory_percentage"
                  stroke="#10B981"
                  fill="#10B981"
                  fillOpacity={0.3}
                  name="Memory Usage (%)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>

          {/* Network Usage Chart */}
          <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-6 border border-primary-400 border-opacity-30 backdrop-blur-sm">
            <h3 className="text-2xl text-text-primary mb-4">Network Usage</h3>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={metricsData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis 
                  dataKey="timestamp" 
                  tickFormatter={formatTime}
                  stroke="#9CA3AF"
                />
                <YAxis 
                  stroke="#9CA3AF"
                  tickFormatter={(value) => formatBytes(value)}
                />
                <Tooltip content={<CustomTooltip />} />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="network_rx"
                  stroke="#8B5CF6"
                  strokeWidth={2}
                  dot={false}
                  name="Network RX (Download)"
                />
                <Line
                  type="monotone"
                  dataKey="network_tx"
                  stroke="#F97316"
                  strokeWidth={2}
                  dot={false}
                  name="Network TX (Upload)"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}
    </div>
  );
};

export default Charts;