import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login as apiLogin, register as apiRegister } from '../lib/api';
import { useAuth } from '../context/AuthContext';
import { IconSkull } from '../components/Icon';

export function Login() {
    const [mode, setMode] = useState<'login' | 'register'>('login');
    const [form, setForm] = useState({
        username: '',
        email: '',
        password: '',
        registration_code: '',
        role: '',
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const set = (k: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
        setForm(f => ({ ...f, [k]: e.target.value }));

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(''); setLoading(true);
        try {
            const res = mode === 'login'
                ? await apiLogin(form.username, form.password)
                : await apiRegister(form.username, form.email, form.password, form.registration_code, form.role);
            login(res.token);
            navigate('/characters');
        } catch (err: any) {
            setError(err?.data?.error || 'Something went wrong');
        } finally { setLoading(false); }
    };

    return (
        <div
            className="min-h-screen flex items-center justify-center bg-ink p-5"
            style={{
                backgroundImage:
                    'radial-gradient(ellipse at 20% 50%, rgba(139,26,26,0.08) 0%, transparent 60%), radial-gradient(ellipse at 80% 20%, rgba(201,168,76,0.05) 0%, transparent 50%)',
            }}
        >
            <div className="w-full max-w-[400px] bg-stone border border-stone-border rounded-lg p-10 shadow-[0_24px_64px_rgba(0,0,0,0.6)]">
                <div className="text-center mb-7">
                    <IconSkull size={40} className="text-crimson-light mx-auto mb-3" />
                    <h1 className="font-display text-[32px] font-black text-parchment tracking-[0.1em]">Behold</h1>
                    <p className="font-body text-[15px] text-ash mt-1">Your digital grimoire</p>
                </div>
                <div className="flex gap-1 bg-stone-mid rounded-md p-1 mb-6">
                    <button className={`flex-1 bg-transparent border-none py-2 rounded-sm text-[13px] font-medium transition-all duration-180 ${mode === 'login' ? 'bg-stone-light text-parchment' : 'text-ash'}`} onClick={() => setMode('login')}>Sign In</button>
                    <button className={`flex-1 bg-transparent border-none py-2 rounded-sm text-[13px] font-medium transition-all duration-180 ${mode === 'register' ? 'bg-stone-light text-parchment' : 'text-ash'}`} onClick={() => setMode('register')}>Register</button>
                </div>
                <form onSubmit={submit} className="flex flex-col gap-4">
                    <div className="flex flex-col gap-1.5">
                        <label className="text-xs font-medium text-ash-light uppercase tracking-wide">Username</label>
                        <input value={form.username} onChange={set('username')} required placeholder="Your adventurer name" />
                    </div>
                    {mode === 'register' && (
                        <div className="flex flex-col gap-1.5">
                            <label className="text-xs font-medium text-ash-light uppercase tracking-wide">Email</label>
                            <input type="email" value={form.email} onChange={set('email')} required placeholder="your@email.com" />
                        </div>
                    )}
                    <div className="flex flex-col gap-1.5">
                        <label className="text-xs font-medium text-ash-light uppercase tracking-wide">Password</label>
                        <input type="password" value={form.password} onChange={set('password')} required placeholder="••••••••" />
                    </div>
                    {mode === 'register' && (
                        <div className="flex flex-col gap-1.5">
                            <label className="text-xs font-medium text-ash-light uppercase tracking-wide">Registration Code</label>
                            <input value={form.registration_code} onChange={set('registration_code')} required placeholder="Provided by your DM" />
                        </div>
                    )}
                    {mode === 'register' && (
                        <div className="flex flex-col gap-1.5">
                            <label className="text-xs font-medium text-ash-light uppercase tracking-wide">Role</label>
                            <input value={form.role} onChange={set('role')} required placeholder="player/dm" />
                        </div>
                    )}
                    {error && <p className="text-[13px] text-crimson-light bg-crimson/15 border border-crimson/30 rounded-sm px-3 py-2">{error}</p>}
                    <button
                        className="bg-crimson border-none text-parchment py-3 rounded-md font-display text-sm font-semibold tracking-wide transition-colors duration-180 mt-1 hover:enabled:bg-crimson-light disabled:opacity-50 disabled:cursor-not-allowed"
                        type="submit" disabled={loading}>
                        {loading ? 'Loading…' : mode === 'login' ? 'Enter the Realm' : 'Join the Party'}
                    </button>
                </form>
            </div>
        </div>
    );
}
