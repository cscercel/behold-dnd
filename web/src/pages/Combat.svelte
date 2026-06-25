<script lang="ts">
  import { onMount } from 'svelte';
  import * as api from '../lib/api';
  import Icon from '../lib/Icon.svelte';
  import ParticipantRow from '../components/ParticipantRow.svelte';

  let encounters   = $state<any[]>([]);
  let active       = $state<any>(null);
  let participants = $state<any[]>([]);
  let characters   = $state<any[]>([]);
  let selectedId   = $state('');
  let newName      = $state('');
  let loading      = $state(true);
  let showAddPart  = $state(false);
  let partForm     = $state({ character_id:'', name:'', initiative:0, max_hp:10, current_hp:10, armor_class:10, speed:30 });

  async function reload() {
    const [enc, chars] = await Promise.all([api.listEncounters(), api.listCharacters()]);
    encounters = enc || []; characters = chars || [];
    try {
      const act = await api.getActiveEncounter();
      active = act;
      if (act) {
        const parts = await api.listParticipants(act.id);
        participants = (parts || []).sort((a:any,b:any) => b.initiative - a.initiative);
      }
    } catch { active = null; participants = []; }
  }

  onMount(() => reload().finally(() => loading = false));

  async function loadEnc(id: string) {
    selectedId = id;
    const parts = await api.listParticipants(id);
    participants = (parts || []).sort((a:any,b:any) => b.initiative - a.initiative);
  }

  async function createEnc() {
    if (!newName.trim()) return;
    await api.createEncounter(newName.trim());
    newName = ''; await reload();
  }

  async function addPart(e: Event) {
    e.preventDefault();
    const encId = active?.id || selectedId; if (!encId) return;
    await api.addParticipant(encId, {
      character_id: partForm.character_id,
      initiative:   partForm.initiative,
    });
    showAddPart = false;
    partForm = { character_id:'', name:'', initiative:0, max_hp:10, current_hp:10, armor_class:10, speed:30 };
    await reload();
  }

  function selectChar(charId: string) {
    const c = characters.find((ch:any) => ch.id === charId);
    if (c) partForm = { ...partForm, character_id:charId, name:c.name, max_hp:c.max_hp, current_hp:c.current_hp, armor_class:c.armor_class, speed:c.speed };
    else   partForm = { ...partForm, character_id:charId };
  }

  $derived: var encId      = active?.id || selectedId;
  $derived: var currentEnc = active || encounters.find((e:any) => e.id === selectedId);
</script>

{#if loading}
  <div class="loading">Loading…</div>
{:else}
<div class="page">

  <!-- Sidebar -->
  <div class="sidebar">
    <h2 class="side-title">Encounters</h2>
    <div class="new-enc">
      <input bind:value={newName} placeholder="Encounter name…"
        onkeydown={(e) => e.key==='Enter' && createEnc()}/>
      <button class="create-btn" onclick={createEnc}><Icon name="plus" size={14}/></button>
    </div>
    <div class="enc-list">
      {#each encounters as enc}
        <div class="enc-item"
          class:active={enc.id===encId}
          class:live={enc.is_active}
          onclick={() => loadEnc(enc.id)}>
          <div>
            <div class="enc-name">{enc.name}</div>
            <div class="enc-meta">{enc.is_active ? `Round ${enc.round} · LIVE` : 'Inactive'}</div>
          </div>
          {#if enc.is_active}<div class="live-dot"></div>{/if}
        </div>
      {/each}
      {#if encounters.length === 0}
        <div class="empty-enc">No encounters yet.</div>
      {/if}
    </div>
  </div>

  <!-- Main -->
  <div class="main">
    {#if currentEnc}
      <div class="enc-header">
        <div>
          <h1 class="enc-title">{currentEnc.name}</h1>
          {#if currentEnc.is_active}
            <span class="round-badge">Round {currentEnc.round}</span>
          {/if}
        </div>
        <div class="enc-actions">
          {#if !currentEnc.is_active}
            <button class="action-btn" onclick={() => api.startEncounter(currentEnc.id).then(reload)}>
              <Icon name="play" size={15}/> Start
            </button>
          {:else}
            <button class="action-btn" onclick={() => api.nextRound(currentEnc.id).then(reload)}>
              <Icon name="skip" size={15}/> Next Round
            </button>
            <button class="action-btn danger" onclick={() => api.endEncounter(currentEnc.id).then(reload)}>
              <Icon name="stop" size={15}/> End
            </button>
          {/if}
          <button class="action-btn" onclick={() => showAddPart = !showAddPart}>
            <Icon name="plus" size={15}/> Add Combatant
          </button>
        </div>
      </div>

      {#if showAddPart}
        <form class="add-form" onsubmit={addPart}>
          <div class="form-row">
            <div class="field">
              <label>Character</label>
              <select value={partForm.character_id}
                onchange={(e) => selectChar((e.target as HTMLSelectElement).value)}>
                <option value="">— Select a character —</option>
                {#each characters as c}
                  <option value={c.id}>{c.name}</option>
                {/each}
              </select>
            </div>
            <div class="field">
              <label>Initiative</label>
              <input type="number" bind:value={partForm.initiative}/>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="cancel-btn" onclick={() => showAddPart=false}>Cancel</button>
            <button type="submit" class="submit-btn" disabled={!partForm.character_id}>Add to Battle</button>
          </div>
        </form>
      {/if}

      <div class="tracker">
        {#if participants.length === 0}
          <div class="empty-tracker">
            <Icon name="sword" size={40} class="empty-icon"/>
            <p>No combatants yet.</p>
          </div>
        {/if}
        {#each participants as p}
          <ParticipantRow
            {p} {encId}
            onDmg={(n) => api.participantDamage(encId, p.id, n).then(reload)}
            onHeal={(n) => api.participantHeal(encId, p.id, n).then(reload)}
            onRemove={() => api.removeParticipant(encId, p.id).then(reload)}
          />
        {/each}
      </div>

    {:else}
      <div class="no-enc">
        <Icon name="sword" size={64} class="no-enc-icon"/>
        <p>Select or create an encounter to begin</p>
      </div>
    {/if}
  </div>

</div>
{/if}

<style>
.loading { display:flex;align-items:center;justify-content:center;height:200px;color:var(--ash); }
.page { display:flex;height:100vh;overflow:hidden; }
.sidebar { width:260px;flex-shrink:0;background:var(--stone);border-right:1px solid var(--stone-border);display:flex;flex-direction:column;padding:20px;gap:12px;overflow-y:auto; }
.side-title { font-family:var(--font-display);font-size:15px;font-weight:700;color:var(--parchment); }
.new-enc { display:flex;gap:6px; }
.new-enc input { flex:1; }
.create-btn { background:var(--crimson);border:none;color:var(--parchment);padding:8px 10px;border-radius:var(--radius-sm); }
.create-btn:hover { background:var(--crimson-light); }
.enc-list { display:flex;flex-direction:column;gap:4px;flex:1; }
.enc-item { padding:10px 12px;border-radius:var(--radius);border:1px solid transparent;cursor:pointer;transition:all var(--transition);display:flex;align-items:center;justify-content:space-between; }
.enc-item:hover { background:var(--stone-mid); }
.enc-item.active { background:var(--stone-mid);border-color:var(--stone-border); }
.enc-item.live { border-color:var(--crimson); }
.enc-name { font-size:13px;font-weight:500;color:var(--parchment); }
.enc-meta { font-size:11px;color:var(--ash);margin-top:2px; }
.live-dot { width:8px;height:8px;border-radius:50%;background:var(--crimson-light);box-shadow:0 0 6px var(--crimson-glow); }
.empty-enc { font-size:13px;color:var(--ash);text-align:center;padding:20px 0; }
.main { flex:1;overflow-y:auto;padding:24px 32px; }
.enc-header { display:flex;align-items:center;justify-content:space-between;margin-bottom:20px;flex-wrap:wrap;gap:12px; }
.enc-title { font-family:var(--font-display);font-size:24px;font-weight:700;color:var(--parchment); }
.round-badge { font-size:12px;background:rgba(139,26,26,.2);color:var(--crimson-light);border:1px solid rgba(139,26,26,.3);padding:3px 10px;border-radius:12px;margin-top:4px;display:inline-block; }
.enc-actions { display:flex;gap:8px;flex-wrap:wrap; }
.action-btn { display:flex;align-items:center;gap:6px;background:var(--stone);border:1px solid var(--stone-border);color:var(--ash-light);padding:8px 14px;border-radius:var(--radius);font-size:13px;transition:all var(--transition); }
.action-btn:hover { border-color:var(--gold-dim);color:var(--parchment); }
.action-btn.danger { background:rgba(139,26,26,.15);border-color:rgba(139,26,26,.3);color:var(--crimson-light); }
.action-btn.danger:hover { background:rgba(139,26,26,.25); }
.add-form { background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:20px;margin-bottom:20px;display:flex;flex-direction:column;gap:12px; }
.form-row { display:flex;gap:12px;flex-wrap:wrap; }
.form-row > * { flex:1;min-width:120px; }
.field { display:flex;flex-direction:column;gap:5px; }
.field label { font-size:11px;color:var(--ash);text-transform:uppercase;letter-spacing:.05em; }
.modal-actions { display:flex;justify-content:flex-end;gap:8px; }
.cancel-btn { background:var(--stone-mid);border:none;color:var(--ash-light);padding:8px 16px;border-radius:var(--radius);font-size:13px; }
.submit-btn { background:var(--crimson);border:none;color:var(--parchment);padding:8px 16px;border-radius:var(--radius);font-family:var(--font-display);font-size:13px;font-weight:600; }
.submit-btn:disabled { opacity:.4;cursor:not-allowed; }
.tracker { display:flex;flex-direction:column;gap:6px; }
.empty-tracker { text-align:center;padding:60px 20px;color:var(--ash); }
:global(.empty-icon) { margin:0 auto 12px;display:block;opacity:.3; }
.no-enc { display:flex;flex-direction:column;align-items:center;justify-content:center;height:60vh;color:var(--ash);gap:12px; }
:global(.no-enc-icon) { opacity:.2; }
</style>
