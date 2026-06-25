<script lang="ts">
  import { onMount } from 'svelte';
  import { initAuth, getUser, isLoading } from './lib/auth.svelte';
  import { getPath, matchPath } from './lib/router.svelte';
  import Layout from './components/Layout.svelte';
  import Login from './pages/Login.svelte';
  import CharacterList from './pages/CharacterList.svelte';
  import CharacterSheet from './pages/CharacterSheet.svelte';
  import Combat from './pages/Combat.svelte';

  onMount(initAuth);
</script>

{#if isLoading()}
  <div style="display:flex;align-items:center;justify-content:center;height:100vh;color:var(--ash)">
    Loading…
  </div>
{:else if !getUser()}
  <Login />
{:else}
  {@const charParams = matchPath('/characters/:id')}
  <Layout>
    {#if charParams}
      <CharacterSheet id={charParams.id} />
    {:else if getPath() === '/combat' && getUser()?.role === 'dm'}
      <Combat />
    {:else}
      <CharacterList />
    {/if}
  </Layout>
{/if}
