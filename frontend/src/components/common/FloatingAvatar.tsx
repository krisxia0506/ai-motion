import React, { useState } from 'react';
import './FloatingAvatar.css';

interface FloatingAvatarProps {
  name: string;
  avatarUrl: string;
  profileUrl: string;
}

export const FloatingAvatar: React.FC<FloatingAvatarProps> = ({ 
  name, 
  avatarUrl, 
  profileUrl 
}) => {
  const [isHovered, setIsHovered] = useState(false);

  const handleClick = () => {
    window.open(profileUrl, '_blank', 'noopener,noreferrer');
  };

  return (
    <div 
      className="floating-avatar"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      onClick={handleClick}
      role="button"
      tabIndex={0}
      aria-label={`Visit ${name}'s profile`}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          handleClick();
        }
      }}
    >
      <div className={`avatar-container ${isHovered ? 'expanded' : ''}`}>
        <div className="avatar-name" aria-hidden={!isHovered}>
          {name}
        </div>
        <div className="avatar-circle">
          <img 
            src={avatarUrl} 
            alt={`${name}'s avatar`}
            className="avatar-image"
          />
        </div>
      </div>
    </div>
  );
};
