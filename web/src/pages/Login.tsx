import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login as apiLogin, register as apiRegister } from '../api';
import { useAuth } from '../context/AuthContext';
import { IconSkull } from '../components/Icon';
import styles from './Login.module.css';

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
    <div className={styles.page}>
      <div className={styles.card}>
        <div className={styles.header}>
          <IconSkull size={40} className={styles.icon} />
          <h1 className={styles.title}>Behold</h1>
          <p className={styles.subtitle}>Your digital grimoire</p>
        </div>
        <div className={styles.tabs}>
          <button className={`${styles.tab} ${mode === 'login' ? styles.tabActive : ''}`} onClick={() => setMode('login')}>Sign In</button>
          <button className={`${styles.tab} ${mode === 'register' ? styles.tabActive : ''}`} onClick={() => setMode('register')}>Register</button>
        </div>
        <form onSubmit={submit} className={styles.form}>
          <div className={styles.field}><label>Username</label><input value={form.username} onChange={set('username')} required placeholder="Your adventurer name" /></div>
          {mode === 'register' && <div className={styles.field}><label>Email</label><input type="email" value={form.email} onChange={set('email')} required placeholder="your@email.com" /></div>}
          <div className={styles.field}><label>Password</label><input type="password" value={form.password} onChange={set('password')} required placeholder="••••••••" /></div>
          {mode === 'register' && <div className={styles.field}><label>Registration Code</label><input value={form.registration_code} onChange={set('registration_code')} required placeholder="Provided by your DM" /></div>}
          {mode === 'register' && <div className={styles.field}><label>Role</label><input value={form.role} onChange={set('role')} required placeholder="player/dm" /></div>}
          {error && <p className={styles.error}>{error}</p>}
          <button className={styles.submit} type="submit" disabled={loading}>
            {loading ? 'Loading…' : mode === 'login' ? 'Enter the Realm' : 'Join the Party'}
          </button>
        </form>
      </div>
    </div>
  );
}
