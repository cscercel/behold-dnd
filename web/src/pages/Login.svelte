<script lang="ts">
  import { login as apiLogin, register as apiRegister } from '../lib/api';
  import { authLogin } from '../lib/auth.svelte';
  import { navigate } from '../lib/router.svelte';
  import Icon from '../lib/Icon.svelte';

  let mode    = $state<'login'|'register'>('login');
  let form    = $state({ username:'', email:'', password:'', registration_code:'' });
  let error   = $state('');
  let loading = $state(false);

  async function submit(e: Event) {
    e.preventDefault(); error = ''; loading = true;
    try {
      const res = mode === 'login'
        ? await apiLogin(form.email, form.password)
        : await apiRegister(form.username, form.email, form.password, form.registration_code);
      await authLogin(res.token);
      navigate('/characters');
    } catch (err: any) {
      error = err?.data?.error || 'Something went wrong';
    } finally { loading = false; }
  }
</script>

<div class="page">
  <div class="card">
    <div class="header">
      <Icon name="skull" size={40} class="icon" />
      <h1 class="title">Behold</h1>
      <p class="subtitle">Your digital grimoire</p>
    </div>
    <div class="tabs">
      <button class="tab" class:active={mode==='login'} onclick={() => mode='login'}>Sign In</button>
      <button class="tab" class:active={mode==='register'} onclick={() => mode='register'}>Register</button>
    </div>
    <form onsubmit={submit}>
      {#if mode === 'register'}
        <div class="field"><label>Username</label><input bind:value={form.username} required placeholder="Your adventurer name"/></div>
      {/if}
      <div class="field"><label>Email</label><input type="email" bind:value={form.email} required placeholder="your@email.com"/></div>
      <div class="field"><label>Password</label><input type="password" bind:value={form.password} required placeholder="••••••••"/></div>
      {#if mode === 'register'}
        <div class="field"><label>Registration Code</label><input bind:value={form.registration_code} required placeholder="Provided by your DM"/></div>
      {/if}
      {#if error}<p class="error">{error}</p>{/if}
      <button class="submit" type="submit" disabled={loading}>
        {loading ? 'Loading…' : mode === 'login' ? 'Enter the Realm' : 'Join the Party'}
      </button>
    </form>
  </div>
</div>

<style>
  .page { min-height:100vh;display:flex;align-items:center;justify-content:center;background:var(--ink);background-image:radial-gradient(ellipse at 20% 50%,rgba(139,26,26,.08) 0%,transparent 60%),radial-gradient(ellipse at 80% 20%,rgba(201,168,76,.05) 0%,transparent 50%);padding:20px; }
  .card { width:100%;max-width:400px;background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:40px;box-shadow:0 24px 64px rgba(0,0,0,.6); }
  .header { text-align:center;margin-bottom:28px; }
  :global(.icon) { color:var(--crimson-light);margin-bottom:12px; }
  .title { font-family:var(--font-display);font-size:32px;font-weight:900;color:var(--parchment);letter-spacing:.1em; }
  .subtitle { font-family:var(--font-body);font-size:15px;color:var(--ash);margin-top:4px; }
  .tabs { display:flex;gap:4px;background:var(--stone-mid);border-radius:var(--radius);padding:4px;margin-bottom:24px; }
  .tab { flex:1;background:none;border:none;color:var(--ash);padding:8px;border-radius:var(--radius-sm);font-size:13px;font-weight:500;transition:all var(--transition); }
  .tab.active { background:var(--stone-light);color:var(--parchment); }
  form { display:flex;flex-direction:column;gap:16px; }
  .field { display:flex;flex-direction:column;gap:6px; }
  .field label { font-size:12px;font-weight:500;color:var(--ash-light);text-transform:uppercase;letter-spacing:.05em; }
  .error { font-size:13px;color:var(--crimson-light);background:rgba(139,26,26,.15);border:1px solid rgba(139,26,26,.3);border-radius:var(--radius-sm);padding:8px 12px; }
  .submit { background:var(--crimson);border:none;color:var(--parchment);padding:12px;border-radius:var(--radius);font-family:var(--font-display);font-size:14px;font-weight:600;letter-spacing:.05em;transition:background var(--transition);margin-top:4px; }
  .submit:hover:not(:disabled) { background:var(--crimson-light); }
  .submit:disabled { opacity:.5;cursor:not-allowed; }
</style>
