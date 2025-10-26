import React, { useState } from 'react';
import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { MdMenu, MdClose, MdHome, MdList, MdLogout, MdAccountCircle, MdMovie } from 'react-icons/md';
import { FloatingAvatar } from './common/FloatingAvatar';
import { useAuth } from '../contexts/AuthContext';
import './Layout.css';

export const Layout: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [showUserMenu, setShowUserMenu] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { user, signOut } = useAuth();

  const navItems = [
    { path: '/', label: '首页', icon: <MdHome size={20} /> },
    { path: '/tasks', label: '任务列表', icon: <MdList size={20} /> },
  ];

  const isActivePath = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    // 任务详情页也应该高亮任务列表
    if (path === '/tasks' && location.pathname.startsWith('/task')) {
      return true;
    }
    return location.pathname.startsWith(path);
  };

  return (
    <div className="layout">
      <header className="layout-header">
        <div className="container">
          <div className="header-content">
            <Link to="/" className="header-logo">
              <MdMovie size={32} />
              <span>AI-Motion</span>
            </Link>

            <nav className="header-nav">
              {navItems.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`nav-link ${isActivePath(item.path) ? 'nav-link-active' : ''}`}
                >
                  {item.icon}
                  <span>{item.label}</span>
                </Link>
              ))}
            </nav>

            <div className="header-user">
              {user ? (
                <div className="user-menu-container">
                  <button
                    className="user-menu-button"
                    onClick={() => setShowUserMenu(!showUserMenu)}
                  >
                    <MdAccountCircle size={24} />
                    <span className="user-email">{user.email}</span>
                  </button>
                  {showUserMenu && (
                    <div className="user-menu-dropdown">
                      <button
                        className="user-menu-item"
                        onClick={async () => {
                          await signOut();
                          navigate('/login');
                        }}
                      >
                        <MdLogout size={20} />
                        <span>退出登录</span>
                      </button>
                    </div>
                  )}
                </div>
              ) : (
                <Link to="/login" className="login-button">
                  登录
                </Link>
              )}
            </div>

            <button
              className="mobile-menu-button"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              aria-label="Toggle menu"
            >
              {isMobileMenuOpen ? <MdClose size={24} /> : <MdMenu size={24} />}
            </button>
          </div>
        </div>
      </header>

      {isMobileMenuOpen && (
        <div className="mobile-menu" onClick={() => setIsMobileMenuOpen(false)}>
          <nav className="mobile-nav">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className={`mobile-nav-link ${isActivePath(item.path) ? 'mobile-nav-link-active' : ''}`}
              >
                {item.icon}
                <span>{item.label}</span>
              </Link>
            ))}
          </nav>
        </div>
      )}

      <main className="layout-main">
        <Outlet />
      </main>

      <footer className="layout-footer">
        <div className="container">
          <p>&copy; {new Date().getFullYear()} AI-Motion. 智能动漫生成系统</p>
          <p className="footer-credit">
            页面由{' '}
            <a 
              href="https://github.com/xgopilot" 
              target="_blank" 
              rel="noopener noreferrer"
              className="footer-link"
            >
              xgopilot
            </a>
            {' '}设计和开发
          </p>
        </div>
      </footer>

      <FloatingAvatar
        name="xgopilot"
        avatarUrl="https://github.com/xgopilot.png"
        profileUrl="https://github.com/xgopilot"
      />
    </div>
  );
};
