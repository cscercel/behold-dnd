import { type ReactNode } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { IconUsers, IconLogOut, IconSkull, IconSword } from './Icon';
import styles from './Layout.module.css';

export function Layout({ children }: { children: ReactNode }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => { logout(); navigate('/login'); };

  const navItems = [
    { to: '/characters', label: 'Characters', icon: <IconUsers size={18} /> },
    ...(user?.role === 'dm' ? [{ to: '/combat', label: 'Combat', icon: <IconSword size={18} /> }] : []),
  ];

  return (
    <div className={styles.shell}>
      <aside className={styles.sidebar}>
        <div className={styles.logo}>
          <IconSkull size={28} className={styles.logoIcon} />
          <span className={styles.logoText}>Behold</span>
        </div>
        <nav className={styles.nav}>
          {navItems.map(item => (
            <Link key={item.to} to={item.to}
              className={`${styles.navItem} ${location.pathname.startsWith(item.to) ? styles.navActive : ''}`}>
              {item.icon}
              <span>{item.label}</span>
            </Link>
          ))}
        </nav>
        <div className={styles.sidebarFooter}>
          <div className={styles.userInfo}>
            <div className={styles.userAvatar}>{user?.username?.[0]?.toUpperCase()}</div>
            <div>
              <div className={styles.userName}>{user?.username}</div>
              <div className={styles.userRole}>{user?.role === 'dm' ? 'Dungeon Master' : 'Player'}</div>
            </div>
          </div>
          <button className={styles.logoutBtn} onClick={handleLogout} title="Log out">
            <IconLogOut size={16} />
          </button>
        </div>
      </aside>
      <main className={styles.main}>{children}</main>
    </div>
  );
}
