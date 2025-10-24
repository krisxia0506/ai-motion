import React from 'react';
import './LoadingSpinner.css';

interface LoadingSpinnerProps {
  size?: 'small' | 'medium' | 'large';
  text?: string;
  fullScreen?: boolean;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'medium',
  text,
  fullScreen = false,
}) => {
  const sizeMap = {
    small: 24,
    medium: 48,
    large: 72,
  };

  const spinnerSize = sizeMap[size];

  const spinner = (
    <div className="loading-spinner-wrapper">
      <svg
        className={`loading-spinner loading-spinner-${size}`}
        width={spinnerSize}
        height={spinnerSize}
        viewBox="0 0 50 50"
      >
        <circle
          className="loading-spinner-circle"
          cx="25"
          cy="25"
          r="20"
          fill="none"
          strokeWidth="4"
        />
      </svg>
      {text && <p className="loading-spinner-text">{text}</p>}
    </div>
  );

  if (fullScreen) {
    return <div className="loading-spinner-fullscreen">{spinner}</div>;
  }

  return spinner;
};
