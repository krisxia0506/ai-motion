import React from 'react';
import './ProgressBar.css';

interface ProgressBarProps {
  value: number;
  max?: number;
  label?: string;
  showPercentage?: boolean;
  size?: 'small' | 'medium' | 'large';
  variant?: 'primary' | 'success' | 'warning' | 'danger';
  className?: string;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  value,
  max = 100,
  label,
  showPercentage = true,
  size = 'medium',
  variant = 'primary',
  className = '',
}) => {
  const percentage = Math.min(100, Math.max(0, (value / max) * 100));
  const classes = ['progress-bar', `progress-bar-${size}`, className]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes}>
      {(label || showPercentage) && (
        <div className="progress-bar-header">
          {label && <span className="progress-bar-label">{label}</span>}
          {showPercentage && (
            <span className="progress-bar-percentage">{Math.round(percentage)}%</span>
          )}
        </div>
      )}
      <div className="progress-bar-track">
        <div
          className={`progress-bar-fill progress-bar-fill-${variant}`}
          style={{ width: `${percentage}%` }}
          role="progressbar"
          aria-valuenow={value}
          aria-valuemin={0}
          aria-valuemax={max}
        />
      </div>
    </div>
  );
};
