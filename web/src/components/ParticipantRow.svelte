<script lang="ts">
  import Icon from '../lib/Icon.svelte';

  let { p, encId: _encId, onDmg, onHeal, onRemove }: {
    p: any;
    encId: string;
    onDmg: (n: number) => void;
    onHeal: (n: number) => void;
    onRemove: () => void;
  } = $props();

  let expanded = $state(false);
  let amount   = $state('');

  const hpPct   = $derived(Math.max(0, Math.min(100, (p.current_hp / p.max_hp) * 100)));
  const hpColor = $derived(hpPct > 50 ? '#27ae60' : hpPct > 25 ? '#f39c12' : '#c0392b');
</script>

<div class="card" class:inactive={!p.is_active}>
  <div class="init">{p.initiative}</div>
  <div class="body">
    <div class="top" onclick={() => expanded = !expanded}>
      <div class="name-block">
        <div class="name">{p.name}</div>
        {#if p.conditions?.length > 0}
          <div class="conditions">
            {#each p.conditions as c}<span class="cond">{c}</span>{/each}
          </div>
        {/if}
      </div>
      <div class="hp-area">
        <div class="hp-bar-wrap"><div class="hp-bar" style="width:{hpPct}%;background:{hpColor}"></div></div>
        <span class="hp-num" style="color:{hpColor}">{p.current_hp}</span>
        <span class="hp-max">/{p.max_hp}</span>
      </div>
      <div class="stats">
        <span>AC {p.armor_class}</span>
        {#if p.concentration}<span class="conc">⊛</span>{/if}
      </div>
      {#if expanded}
        <Icon name="chevronUp" size={14} class="chevron"/>
      {:else}
        <Icon name="chevronDown" size={14} class="chevron"/>
      {/if}
    </div>
    {#if expanded}
      <div class="expanded">
        <input class="amount-input" type="number" min="1" bind:value={amount} placeholder="Amount"/>
        <button class="dmg-btn" onclick={() => { if (amount) { onDmg(+amount); amount = ''; } }}>⚔ Damage</button>
        <button class="heal-btn" onclick={() => { if (amount) { onHeal(+amount); amount = ''; } }}>❤ Heal</button>
        <button class="remove-btn" onclick={onRemove}><Icon name="trash" size={13}/></button>
      </div>
    {/if}
  </div>
</div>

<style>
.card { display:flex;align-items:stretch;background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);overflow:hidden; }
.card.inactive { opacity:.4; }
.init { width:48px;flex-shrink:0;background:var(--stone-mid);display:flex;align-items:center;justify-content:center;font-family:var(--font-display);font-size:18px;font-weight:700;color:var(--gold); }
.body { flex:1; }
.top { display:flex;align-items:center;gap:12px;padding:12px 16px;cursor:pointer; }
.top:hover { background:rgba(255,255,255,.02); }
.name-block { flex:1; }
.name { font-size:15px;font-weight:500;color:var(--parchment); }
.conditions { display:flex;gap:4px;flex-wrap:wrap;margin-top:3px; }
.cond { font-size:10px;background:rgba(139,26,26,.2);color:var(--crimson-light);border-radius:3px;padding:1px 5px; }
.hp-area { display:flex;align-items:center;gap:6px; }
.hp-bar-wrap { width:80px;height:6px;background:var(--stone-border);border-radius:3px;overflow:hidden; }
.hp-bar { height:100%;border-radius:3px; }
.hp-num { font-weight:700;font-size:16px; }
.hp-max { font-size:12px;color:var(--ash); }
.stats { display:flex;gap:10px;font-size:12px;color:var(--ash); }
.conc { color:#9b59b6; }
:global(.chevron) { color:var(--ash);flex-shrink:0; }
.expanded { padding:12px 16px;border-top:1px solid var(--stone-border);display:flex;align-items:center;gap:10px;background:var(--stone-mid); }
.amount-input { width:80px;text-align:center;padding:6px 8px; }
.dmg-btn { background:rgba(139,26,26,.2);border:1px solid rgba(139,26,26,.3);color:var(--crimson-light);padding:6px 12px;border-radius:var(--radius-sm);font-size:12px; }
.dmg-btn:hover { background:rgba(139,26,26,.35); }
.heal-btn { background:rgba(39,174,96,.15);border:1px solid rgba(39,174,96,.3);color:#27ae60;padding:6px 12px;border-radius:var(--radius-sm);font-size:12px; }
.heal-btn:hover { background:rgba(39,174,96,.25); }
.remove-btn { display:flex;align-items:center;gap:4px;background:none;border:none;color:var(--ash);font-size:12px;margin-left:auto; }
.remove-btn:hover { color:var(--crimson-light); }
</style>
