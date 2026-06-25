<script lang="ts">
  import { getUser, authLogout } from '../lib/auth.svelte';
  import { navigate, getPath } from '../lib/router.svelte';
  import Icon from '../lib/Icon.svelte';

  let { children } = $props();

  function logout() { authLogout(); navigate('/'); }
  function active(path: string) { return getPath().startsWith(path); }
</script>

<div class="shell">
  <aside class="sidebar">
    <div class="logo">
      <Icon name="skull" size={28} class="logo-icon" />
      <span class="logo-text">Behold</span>
    </div>
    <nav class="nav">
      <a href="/characters" class="nav-item" class:active={active('/characters')}
        onclick={(e) => { e.preventDefault(); navigate('/characters'); }}>
        <Icon name="users" size={18}/> <span>Characters</span>
      </a>
      {#if getUser()?.role === 'dm'}
        <a href="/combat" class="nav-item" class:active={active('/combat')}
          onclick={(e) => { e.preventDefault(); navigate('/combat'); }}>
          <Icon name="sword" size={18}/> <span>Combat</span>
        </a>
      {/if}
    </nav>
    <div class="sidebar-footer">
      <div class="user-info">
        <div class="user-avatar">{getUser()?.username?.[0]?.toUpperCase()}</div>
        <div>
          <div class="user-name">{getUser()?.username}</div>
          <div class="user-role">{getUser()?.role === 'dm' ? 'Dungeon Master' : 'Player'}</div>
        </div>
      </div>
      <button class="logout-btn" onclick={logout} title="Log out">
        <Icon name="logOut" size={16}/>
      </button>
    </div>
  </aside>
  <main class="main">{@render children()}</main>
</div>

<style>
  .shell { display:flex;height:100vh;overflow:hidden; }
  .sidebar { width:220px;flex-shrink:0;background:var(--stone);border-right:1px solid var(--stone-border);display:flex;flex-direction:column; }
  .logo { display:flex;align-items:center;gap:10px;padding:20px 20px 16px;border-bottom:1px solid var(--stone-border); }
  :global(.logo-icon) { color:var(--crimson-light); }
  .logo-text { font-family:var(--font-display);font-size:20px;font-weight:700;color:var(--parchment);letter-spacing:0.05em; }
  .nav { flex:1;padding:16px 12px;display:flex;flex-direction:column;gap:4px; }
  .nav-item { display:flex;align-items:center;gap:10px;padding:10px 12px;border-radius:var(--radius);color:var(--ash-light);font-size:13px;font-weight:500;transition:all var(--transition); }
  .nav-item:hover { background:var(--stone-mid);color:var(--parchment); }
  .nav-item.active { background:var(--stone-mid);color:var(--gold);border-left:3px solid var(--gold); }
  .sidebar-footer { padding:16px;border-top:1px solid var(--stone-border);display:flex;align-items:center;gap:10px; }
  .user-info { display:flex;align-items:center;gap:10px;flex:1;min-width:0; }
  .user-avatar { width:32px;height:32px;border-radius:50%;background:var(--crimson);display:flex;align-items:center;justify-content:center;font-family:var(--font-display);font-size:13px;font-weight:700;color:var(--parchment);flex-shrink:0; }
  .user-name { font-size:13px;font-weight:500;color:var(--parchment);white-space:nowrap;overflow:hidden;text-overflow:ellipsis; }
  .user-role { font-size:11px;color:var(--ash); }
  .logout-btn { background:none;border:none;color:var(--ash);padding:6px;border-radius:var(--radius-sm);display:flex;align-items:center;transition:color var(--transition); }
  .logout-btn:hover { color:var(--crimson-light); }
  .main { flex:1;overflow-y:auto;background:var(--ink); }
</style>
