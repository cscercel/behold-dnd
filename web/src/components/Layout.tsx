import { type ReactNode } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { IconUsers, IconLogOut, IconSkull, IconSword } from './Icon';

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
        <div className="flex h-screen overflow-hidden">
            <aside className="w-[220px] shrink-0 bg-stone border-r border-stone-border flex flex-col">
                <div className="flex items-center gap-2.5 px-5 pt-5 pb-4 border-b border-stone-border">
                    <IconSkull size={28} className="text-crimson-light" />
                    <span className="font-display text-xl font-bold text-parchment tracking-wide">Behold</span>
                </div>
                <nav className="flex-1 p-3 flex flex-col gap-1">
                    {navItems.map(item => (
                        <Link key={item.to} to={item.to}
                            className={`flex items-center gap-2.5 px-3 py-2.5 rounded-md text-ash-light text-[13px] font-medium transition-all duration-180 hover:bg-stone-mid hover:text-parchment ${location.pathname.startsWith(item.to) ? 'bg-stone-mid text-gold border-l-[3px] border-gold' : ''}`}>
                            {item.icon}
                            <span>{item.label}</span>
                        </Link>
                    ))}
                </nav>
                <div className="p-4 border-t border-stone-border flex items-center gap-2.5">
                    <div className="flex items-center gap-2.5 flex-1 min-w-0">
                        <div className="w-8 h-8 rounded-full bg-crimson flex items-center justify-center font-display text-[13px] font-bold text-parchment shrink-0">
                            {user?.username?.[0]?.toUpperCase()}
                        </div>
                        <div>
                            <div className="text-[13px] font-medium text-parchment whitespace-nowrap overflow-hidden text-ellipsis">{user?.username}</div>
                            <div className="text-[11px] text-ash">{user?.role === 'dm' ? 'Dungeon Master' : 'Player'}</div>
                        </div>
                    </div>
                    <button className="bg-transparent border-none text-ash p-1.5 rounded-sm flex items-center transition-colors duration-180 hover:text-crimson-light" onClick={handleLogout} title="Log out">
                        <IconLogOut size={16} />
                    </button>
                </div>
            </aside>
            <main className="flex-1 overflow-y-auto bg-ink">{children}</main>
        </div>
    );
}
