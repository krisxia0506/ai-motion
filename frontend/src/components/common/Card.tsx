import React from 'react';
import './Card.css';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  hoverable?: boolean;
  onClick?: () => void;
}

export const Card: React.FC<CardProps> = ({
  children,
  className = '',
  hoverable = false,
  onClick,
}) => {
  const classes = [
    'card',
    hoverable && 'card-hoverable',
    onClick && 'card-clickable',
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <div className={classes} onClick={onClick}>
      {children}
    </div>
  );
};

export const CardHeader: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => <div className="card-header">{children}</div>;

export const CardBody: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => <div className="card-body">{children}</div>;

export const CardFooter: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => <div className="card-footer">{children}</div>;
