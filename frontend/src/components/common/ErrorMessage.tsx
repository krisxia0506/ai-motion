import React from 'react';
import './ErrorMessage.css';

interface ErrorMessageProps {
  title?: string;
  message: string;
  onRetry?: () => void;
  className?: string;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({
  title = 'Error',
  message,
  onRetry,
  className = '',
}) => {
  return (
    <div className={`error-message ${className}`}>
      <div className="error-message-icon">⚠️</div>
      <div className="error-message-content">
        <h3 className="error-message-title">{title}</h3>
        <p className="error-message-text">{message}</p>
        {onRetry && (
          <button onClick={onRetry} className="error-message-retry-btn">
            Retry
          </button>
        )}
      </div>
    </div>
  );
};
