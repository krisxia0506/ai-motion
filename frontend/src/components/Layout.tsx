import React, { useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { MdMenu, MdClose, MdHome, MdLibraryBooks, MdPerson, MdMovie, MdFileDownload } from 'react-icons/md';
import './Layout.css';

export const Layout: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const location = useLocation();

  const navItems = [
    { path: '/', label: '首页', icon: <MdHome size={20} /> },
    { path: '/novels', label: '小说列表', icon: <MdLibraryBooks size={20} /> },
    { path: '/characters', label: '角色管理', icon: <MdPerson size={20} /> },
    { path: '/generation', label: '生成动漫', icon: <MdMovie size={20} /> },
    { path: '/export', label: '导出视频', icon: <MdFileDownload size={20} /> },
  ];

  const isActivePath = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
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
        </div>
      </footer>
    </div>
  );
};
