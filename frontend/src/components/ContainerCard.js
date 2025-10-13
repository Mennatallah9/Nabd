import React from 'react';

const ContainerCard = ({ container, metric, onRestart, onViewLogs, onViewCharts }) => {
  const getStatusColor = (status) => {
    if (status.includes('Up')) return 'text-green-400';
    if (status.includes('Exited')) return 'text-red-400';
    return 'text-yellow-400';
  };

  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatPercentage = (usage, limit) => {
    if (!limit) return '0%';
    return ((usage / limit) * 100).toFixed(1) + '%';
  };

  const truncateImageName = (imageName) => {
    if (!imageName) return '';
    // If the image name is longer than 40 characters, truncate it
    if (imageName.length > 40) {
      return imageName.substring(0, 37) + '...';
    }
    return imageName;
  };

  const truncateContainerId = (containerId) => {
    if (!containerId) return '';
    // Show first 12 characters of container ID (standard Docker short ID length)
    return containerId.substring(0, 12);
  };

  return (
    <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
      <div className="flex items-center justify-between mb-6">
        <div className="flex-1 min-w-0 mr-4">
          <h3 className="text-xl text-text-primary truncate">
            {container.name}
          </h3>
          <p className="text-base text-text-secondary mt-1 truncate" title={container.image}>
            {truncateImageName(container.image)}
          </p>
        </div>
        <div className="flex space-x-2 flex-shrink-0">
          <button
            onClick={() => onViewCharts(container.name)}
            className="px-3 py-2 bg-purple-500 bg-opacity-20 text-purple-200 border border-purple-400 border-opacity-50 rounded-lg text-sm hover:bg-opacity-30 transition-all duration-200 backdrop-blur-sm"
          >
            Charts
          </button>
          <button
            onClick={() => onViewLogs(container.name)}
            className="px-3 py-2 bg-blue-500 bg-opacity-20 text-blue-200 border border-blue-400 border-opacity-50 rounded-lg text-sm hover:bg-opacity-30 transition-all duration-200 backdrop-blur-sm"
          >
            Logs
          </button>
          <button
            onClick={() => onRestart(container.name)}
            className="px-3 py-2 bg-green-500 bg-opacity-20 text-green-200 border border-green-400 border-opacity-50 rounded-lg text-sm hover:bg-opacity-30 transition-all duration-200 backdrop-blur-sm"
          >
            Restart
          </button>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6 mb-6">
        <div>
          <p className="text-base text-text-secondary">Status</p>
          <p className={`text-lg ${getStatusColor(container.status)}`}>
            {container.state}
          </p>
        </div>
        <div>
          <p className="text-base text-text-secondary">Container ID</p>
          <p className="text-base font-mono text-text-primary truncate" title={container.id}>
            {truncateContainerId(container.id)}
          </p>
        </div>
      </div>

      {metric && (
        <div className="grid grid-cols-3 gap-6 pt-6 border-t border-primary-400 border-opacity-30">
          <div>
            <p className="text-base text-text-secondary">CPU</p>
            <p className="text-2xl text-text-primary">
              {metric.cpu_percent.toFixed(1)}%
            </p>
          </div>
          <div>
            <p className="text-base text-text-secondary">Memory</p>
            <p className="text-2xl text-text-primary">
              {formatPercentage(metric.memory_usage, metric.memory_limit)}
            </p>
            <p className="text-sm text-text-secondary">
              {formatBytes(metric.memory_usage)} / {formatBytes(metric.memory_limit)}
            </p>
          </div>
          <div>
            <p className="text-base text-text-secondary">Network</p>
            <p className="text-base text-text-primary">
              ↓ {formatBytes(metric.network_rx)}
            </p>
            <p className="text-base text-text-primary">
              ↑ {formatBytes(metric.network_tx)}
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default ContainerCard;