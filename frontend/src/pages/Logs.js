import React, { useState, useEffect } from 'react';
import { containerAPI } from '../services/api';
import LoadingSpinner from '../components/LoadingSpinner';

const Logs = ({ selectedContainer, onClose }) => {
  const [containers, setContainers] = useState([]);
  const [currentContainer, setCurrentContainer] = useState(selectedContainer || '');
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [lines, setLines] = useState(100);
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [summary, setSummary] = useState(null);
  const [summarizing, setSummarizing] = useState(false);
  const [summaryError, setSummaryError] = useState('');
  const [showSummary, setShowSummary] = useState(false);

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
    setSummary(null);
    setSummaryError('');
    setShowSummary(false);
    fetchLogs(containerName);
  };

  const handleSummarizeLogs = async () => {
    if (!currentContainer) return;
    
    setSummarizing(true);
    setSummaryError('');
    
    try {
      if (!('Summarizer' in window)) {
        throw new Error('AI Summarizer is not available. Please ensure you are using Chrome 138+ and enable experimental AI features in chrome://flags/');
      }

      const response = await containerAPI.getLogs(currentContainer, 500);
      let logText = response.data.data.join('\n');
      
      if (!logText || logText.trim().length === 0) {
        throw new Error('No logs available to summarize');
      }

      const maxChars = 10000;
      if (logText.length > maxChars) {
        //most recent logs
        logText = '...[earlier logs truncated]...\n' + logText.slice(-maxChars);
      }

      // check availability
      let availability;
      try {
        if (typeof window.Summarizer.availability === 'function') {
          availability = await window.Summarizer.availability();
          console.log('Summarizer availability:', availability);
        } else if (typeof window.Summarizer.capabilities === 'function') {
          const caps = await window.Summarizer.capabilities();
          availability = caps.available;
        } else {
          availability = 'ready';
        }
      } catch (e) {
        console.warn('Could not check availability:', e);
        availability = 'ready';
      }
      
      if (availability === 'no' || availability === 'unavailable') {
        setSummaryError('⚠️ System requirements not fully met');
      }

      if (availability === 'after-download' || availability === 'downloadable') {
        setSummaryError('Downloading AI model for the first time. This may take several minutes...');
      }

      let summarizer;
      try {
        const optimizedOptions = {
          type: 'key-points',     
          format: 'plain-text',
          length: 'short',
          sharedContext: `Docker container "${currentContainer}" logs. Focus on critical errors and warnings only.`,
          monitor(m) {
            m.addEventListener('downloadprogress', (e) => {
              const percentComplete = Math.round(e.loaded * 100);
              console.log(`Downloaded ${percentComplete}%`);
              setSummaryError(`Downloading AI model: ${percentComplete}% complete...`);
            });
          }
        };
        
        summarizer = await window.Summarizer.create(optimizedOptions);
      } catch (createError) {
        console.error('First create attempt failed:', createError);
        
        try {
          const minimalOptions = {
            type: 'tldr',
            format: 'plain-text',
            length: 'short',
            monitor(m) {
              m.addEventListener('downloadprogress', (e) => {
                const percentComplete = Math.round(e.loaded * 100);
                console.log(`Downloaded ${percentComplete}% (minimal config)`);
                setSummaryError(`Downloading AI model (minimal): ${percentComplete}% complete...`);
              });
            }
          };
          
          summarizer = await window.Summarizer.create(minimalOptions);
        } catch (minimalError) {
          try {
            const basicOptions = {
              monitor(m) {
                m.addEventListener('downloadprogress', (e) => {
                  const percentComplete = Math.round(e.loaded * 100);
                  console.log(`Downloaded ${percentComplete}% (basic config)`);
                  setSummaryError(`Downloading AI model (basic): ${percentComplete}% complete...`);
                });
              }
            };
            
            summarizer = await window.Summarizer.create(basicOptions);
          } catch (finalError) {
            throw new Error(`Failed to create summarizer with all attempted configurations. Original error: ${createError.message}.`);
          }
        }
      }

      // Clear any download messages
      setSummaryError('');

      // Generate the summary - try different method names
      let summaryText;
      try {
        if (typeof summarizer.summarize === 'function') {
          summaryText = await summarizer.summarize(logText, {
            context: 'Analyze these container logs for errors, warnings, performance issues, and significant events. Provide actionable insights and recommendations.'
          });
        } else if (typeof summarizer.generateSummary === 'function') {
          summaryText = await summarizer.generateSummary(logText);
        } else if (typeof summarizer.run === 'function') {
          summaryText = await summarizer.run(logText);
        } else {
          throw new Error('No summarization method found on the summarizer object');
        }
      } catch (summaryError) {
        throw new Error(`Failed to generate summary: ${summaryError.message}`);
      }

      const summaryData = {
        container_name: currentContainer,
        summary: summaryText,
        original_lines: response.data.data.length,
        timestamp: new Date().toISOString()
      };

      setSummary(summaryData);
      setShowSummary(true);

      // Clean up the summarizer
      if (summarizer && typeof summarizer.destroy === 'function') {
        summarizer.destroy();
      }
      
    } catch (err) {
      console.error('Summarization error:', err);
      
      let errorMessage = err.message;

      if (err.name === 'NotSupportedError') {
        errorMessage = 'AI features are not enabled. Please enable experimental AI flags in Chrome settings.';
      } else if (err.name === 'InvalidStateError') {
        errorMessage = 'AI model is not ready. Please try again in a few minutes.';
      } else if (err.name === 'AbortError') {
        errorMessage = 'AI operation was cancelled. Please try again.';
      } else if (err.message.includes('memory') || err.message.includes('RAM')) {
        errorMessage = `Memory limitation detected. Try these steps:\n\n1. Close other Chrome tabs/applications\n2. Go to chrome://flags/ and enable:\n   - "#enable-experimental-web-platform-features"\n   - "#optimization-guide-on-device-model" \n   - "#prompt-api-for-gemini-nano"\n3. Restart Chrome\n4. If still failing, the model may need to download first. Wait 10-15 minutes and try again.\n\nOriginal error: ${err.message}`;
      }
      
      const debugInfo = [];
      debugInfo.push(`Chrome version: ${navigator.userAgent}`);
      debugInfo.push(`Has window.ai: ${!!window.ai}`);
      debugInfo.push(`Has Summarizer: ${!!window.Summarizer}`);
      
      if (window.Summarizer) {
        debugInfo.push(`Summarizer type: ${typeof window.Summarizer}`);
        debugInfo.push(`Summarizer methods: ${Object.getOwnPropertyNames(window.Summarizer).join(', ')}`);
      }
      
      setSummaryError(errorMessage + '\n\nTroubleshooting:\n' + debugInfo.join('\n'));
      setSummary(null);
    } finally {
      setSummarizing(false);
    }
  };

  const checkAISupport = async () => {
    const info = [];
    info.push(`Browser: ${navigator.userAgent}`);
    info.push(`window.ai: ${!!window.ai}`);
    info.push(`window.Summarizer: ${!!window.Summarizer}`);
    info.push(`chrome.aiOriginTrial: ${!!(window.chrome && window.chrome.aiOriginTrial)}`);
    
    if (window.ai) {
      info.push(`window.ai.summarizer: ${!!window.ai.summarizer}`);
      info.push(`window.ai keys: ${Object.keys(window.ai).join(', ')}`);
    }
    
    if (window.Summarizer) {
      info.push(`Summarizer type: ${typeof window.Summarizer}`);
      info.push(`Summarizer prototype: ${Object.getOwnPropertyNames(window.Summarizer).join(', ')}`);
      
      if (window.Summarizer.prototype) {
        info.push(`Summarizer prototype methods: ${Object.getOwnPropertyNames(window.Summarizer.prototype).join(', ')}`);
      }
      
      try {
        const descriptor = Object.getOwnPropertyDescriptor(window, 'Summarizer');
        info.push(`Summarizer descriptor: ${JSON.stringify(descriptor, null, 2)}`);
      } catch (e) {
        info.push(`Error getting descriptor: ${e.message}`);
      }
    }
    
    alert('AI Support Debug Info:\n\n' + info.join('\n'));
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
            className="px-4 py-2 bg-primary-500 bg-opacity-20 text-primary-200 border border-primary-400 border-opacity-50 rounded-lg hover:bg-opacity-30 transition-all duration-200 mb-4"
          >
            ← Back to Dashboard
          </button>
        )}
      </div>

      {/* Controls */}
      <div className="bg-primary-900 bg-opacity-40 rounded-lg shadow-xl p-8 mb-8 border border-primary-400 border-opacity-30 backdrop-blur-sm">
        <div className="grid grid-cols-1 md:grid-cols-6 gap-6">
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
              {loading ? <LoadingSpinner size="small" /> : 'Refresh'}
            </button>
          </div>

          <div className="flex items-end">
            <button
              onClick={handleSummarizeLogs}
              disabled={!currentContainer || summarizing}
              className="w-full px-6 py-3 bg-blue-500 text-white text-lg font-medium rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg"
            >
              {summarizing ? 'Summarizing...' : 'Summarise Logs'}
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

      {summaryError && (
        <div className="mb-8 bg-red-900 bg-opacity-30 border border-red-400 text-red-200 px-6 py-4 rounded-lg backdrop-blur-sm">
          <strong>Summary Error:</strong> {summaryError}
        </div>
      )}

      {/* AI Summary Display */}
      {summary && showSummary && (
        <div className="bg-blue-900 bg-opacity-40 rounded-lg shadow-xl p-8 mb-8 border border-blue-400 border-opacity-30 backdrop-blur-sm">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-2xl font-semibold text-blue-300">
              AI Log Summary - {summary.container_name}
            </h3>
            <button
              onClick={() => setShowSummary(false)}
              className="text-blue-300 hover:text-blue-200 transition-colors duration-200"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div className="bg-blue-900 bg-opacity-30 rounded-lg p-6 border border-blue-400 border-opacity-20">
            <p className="text-blue-100 leading-relaxed text-lg whitespace-pre-wrap">
              {summary.summary}
            </p>
            <div className="mt-4 pt-4 border-t border-blue-400 border-opacity-20">
              <div className="flex justify-between text-blue-300 text-sm">
                <span>Analyzed {summary.original_lines} log lines</span>
                <span>Generated at {new Date(summary.timestamp).toLocaleString()}</span>
              </div>
            </div>
          </div>
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
              <div className="text-text-secondary text-center py-12 text-lg flex flex-col items-center">
                {loading ? (
                  <>
                    <LoadingSpinner size="small" />
                  </>
                ) : (
                  'No logs available'
                )}
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