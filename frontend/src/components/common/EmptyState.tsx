import React from 'react';
import './EmptyState.css';

interface EmptyStateProps {
  icon?: React.ReactNode | string;
  title: string;
  description?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
  className?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  icon = '📭',
  title,
  description,
  action,
  className = '',
}) => {
  return (
    <div className={`empty-state ${className}`}>
      <div className="empty-state-icon">
        {typeof icon === 'string' ? <span>{icon}</span> : icon}
      </div>
      <h3 className="empty-state-title">{title}</h3>
      {description && <p className="empty-state-description">{description}</p>}
      {action && (
        <button onClick={action.onClick} className="empty-state-action-btn">
          {action.label}
        </button>
      )}
    </div>
  );
};
