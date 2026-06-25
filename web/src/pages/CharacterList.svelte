<script lang="ts">
  import { onMount } from 'svelte';
  import { listCharacters, createCharacter, deleteCharacter } from '../lib/api';
  import { getUser } from '../lib/auth.svelte';
  import { navigate } from '../lib/router.svelte';
  import Icon from '../lib/Icon.svelte';

  const CLASS_COLORS: Record<string,string> = {
    Barbarian:'#c0392b',Bard:'#8e44ad',Cleric:'#f39c12',Druid:'#27ae60',
    Fighter:'#7f8c8d',Monk:'#16a085',Paladin:'#f1c40f',Ranger:'#2ecc71',
    Rogue:'#2c3e50',Sorcerer:'#e74c3c',Warlock:'#6c3483',Wizard:'#2980b9',
  };
  const RACES    = ['Human','Elf','Dwarf','Halfling','Gnome','Half-Elf','Half-Orc','Tiefling','Dragonborn','Aasimar','Tabaxi','Kenku','Firbolg'];
  const CLASSES  = ['Barbarian','Bard','Cleric','Druid','Fighter','Monk','Paladin','Ranger','Rogue','Sorcerer','Warlock','Wizard','Artificer'];
  const ALIGNMENTS = ['Lawful Good','Neutral Good','Chaotic Good','Lawful Neutral','True Neutral','Chaotic Neutral','Lawful Evil','Neutral Evil','Chaotic Evil'];

  let chars = $state<any[]>([]);
  let loading = $state(true);
  let showCreate = $state(false);
  let createError = $state('');
  let form = $state({
    name:'',race:'Human',class:'Fighter',level:1,background:'',alignment:'True Neutral',
    strength:10,dexterity:10,constitution:10,intelligence:10,wisdom:10,charisma:10,
    max_hp:10,current_hp:10,armor_class:10,speed:30,hit_dice_type:8,hit_dice_remaining:1,
    spellcasting_ability:'none',is_npc:false,
  });

  onMount(async () => {
    try { chars = await listCharacters() || []; } finally { loading = false; }
  });

  async function handleDelete(id: string, e: MouseEvent) {
    e.stopPropagation();
    if (!confirm('Delete this character permanently?')) return;
    await deleteCharacter(id);
    chars = chars.filter(c => c.id !== id);
  }

  async function handleCreate(e: Event) {
    e.preventDefault(); createError = '';
    try {
      const user = getUser();
      const newChar = await createCharacter({
        ...form,
        owner_id: user?.id || '',
        xp:0,temp_hp:0,inspiration:false,attunement_slots:3,
        death_save_successes:0,death_save_failures:0,
        training_armor:[],training_weapons:[],training_tools:[],training_languages:[],
        copper:0,silver:0,electrum:0,gold:0,platinum:0,
        conditions:[],resistances:[],vulnerabilities:[],immunities:[],
        personality_traits:'',ideals:'',bonds:'',flaws:'',notes:'',
        save_prof_strength:false,save_prof_dexterity:false,save_prof_constitution:false,
        save_prof_intelligence:false,save_prof_wisdom:false,save_prof_charisma:false,
        skill_acrobatics:0,skill_animal_handling:0,skill_arcana:0,skill_athletics:0,
        skill_deception:0,skill_history:0,skill_insight:0,skill_intimidation:0,
        skill_investigation:0,skill_medicine:0,skill_nature:0,skill_perception:0,
        skill_performance:0,skill_persuasion:0,skill_religion:0,skill_sleight_of_hand:0,
        skill_stealth:0,skill_survival:0,
      });
      chars = [...chars, newChar];
      showCreate = false;
      navigate(`/characters/${newChar.id}`);
    } catch (err: any) { createError = err?.data?.error || 'Failed to create'; }
  }
</script>

<div class="page">
  <div class="header">
    <div>
      <h1 class="title">Your Characters</h1>
      <p class="subtitle">{chars.length} adventurer{chars.length !== 1 ? 's' : ''} in your party</p>
    </div>
    <button class="create-btn" onclick={() => showCreate = true}>
      <Icon name="plus" size={16}/> New Character
    </button>
  </div>

  {#if loading}
    <div class="loading">Loading characters…</div>
  {:else if chars.length === 0}
    <div class="empty"><Icon name="sword" size={48} class="empty-icon"/><p>No characters yet.</p></div>
  {:else}
    <div class="grid">
      {#each chars as char}
        <div class="card" onclick={() => navigate(`/characters/${char.id}`)}>
          <div class="card-accent" style="background:{CLASS_COLORS[char.class]||'#8b1a1a'}"></div>
          <div class="card-body">
            <div class="card-top">
              <div>
                <h2 class="char-name">{char.name}</h2>
                <p class="char-meta">Level {char.level} {char.race} {char.class}</p>
              </div>
              <button class="delete-btn" onclick={(e) => handleDelete(char.id, e)}>
                <Icon name="trash" size={14}/>
              </button>
            </div>
            <div class="card-stats">
              <span class="stat"><Icon name="shield" size={12}/> {char.armor_class} AC</span>
              <span class="stat">❤ {char.current_hp}/{char.max_hp}</span>
              {#if char.is_npc}<span class="npc-tag">NPC</span>{/if}
            </div>
          </div>
          <Icon name="chevronRight" size={16} class="chevron"/>
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showCreate}
  <div class="overlay" onclick={() => showCreate = false}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">Create Character</h2>
      <form onsubmit={handleCreate}>
        <div class="form-row">
          <div class="field"><label>Name *</label><input bind:value={form.name} required placeholder="Character name"/></div>
          <div class="field"><label>Level</label><input type="number" min="1" max="20" bind:value={form.level}/></div>
        </div>
        <div class="form-row">
          <div class="field"><label>Race</label><select bind:value={form.race}>{#each RACES as r}<option>{r}</option>{/each}</select></div>
          <div class="field"><label>Class</label><select bind:value={form.class}>{#each CLASSES as c}<option>{c}</option>{/each}</select></div>
        </div>
        <div class="form-row">
          <div class="field"><label>Background</label><input bind:value={form.background} placeholder="Soldier, Sage…"/></div>
          <div class="field"><label>Alignment</label><select bind:value={form.alignment}>{#each ALIGNMENTS as a}<option>{a}</option>{/each}</select></div>
        </div>
        <div class="section-label">Ability Scores</div>
        <div class="ability-grid">
          {#each ['strength','dexterity','constitution','intelligence','wisdom','charisma'] as ab}
            <div class="ability-field">
              <label>{ab.slice(0,3).toUpperCase()}</label>
              <input type="number" min="1" max="30" bind:value={form[ab as keyof typeof form]}/>
            </div>
          {/each}
        </div>
        <div class="section-label">Combat Stats</div>
        <div class="form-row-4">
          <div class="field"><label>Max HP</label><input type="number" min="1" bind:value={form.max_hp}/></div>
          <div class="field"><label>AC</label><input type="number" min="1" bind:value={form.armor_class}/></div>
          <div class="field"><label>Speed</label><input type="number" min="0" bind:value={form.speed}/></div>
          <div class="field"><label>Hit Die</label><input type="number" min="4" bind:value={form.hit_dice_type}/></div>
        </div>
        {#if createError}<p class="error">{createError}</p>{/if}
        <div class="modal-actions">
          <button type="button" class="cancel-btn" onclick={() => showCreate = false}>Cancel</button>
          <button type="submit" class="submit-btn">Create Character</button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .page { padding:32px;max-width:1100px;margin:0 auto; }
  .loading,.empty { display:flex;align-items:center;justify-content:center;height:200px;color:var(--ash); }
  .empty { flex-direction:column;gap:12px;padding:80px 20px; }
  :global(.empty-icon) { opacity:.3; }
  .header { display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:32px; }
  .title { font-family:var(--font-display);font-size:28px;font-weight:700;color:var(--parchment); }
  .subtitle { font-size:14px;color:var(--ash);margin-top:4px; }
  .create-btn { display:flex;align-items:center;gap:8px;background:var(--crimson);border:none;color:var(--parchment);padding:10px 18px;border-radius:var(--radius);font-family:var(--font-display);font-size:13px;font-weight:600;letter-spacing:.04em;transition:background var(--transition); }
  .create-btn:hover { background:var(--crimson-light); }
  .grid { display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:16px; }
  .card { display:flex;align-items:stretch;background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);cursor:pointer;overflow:hidden;transition:border-color var(--transition),transform var(--transition); }
  .card:hover { border-color:var(--gold-dim);transform:translateY(-2px); }
  .card-accent { width:4px;flex-shrink:0; }
  .card-body { flex:1;padding:16px; }
  .card-top { display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px; }
  .char-name { font-family:var(--font-display);font-size:17px;font-weight:600;color:var(--parchment); }
  .char-meta { font-size:12px;color:var(--ash);margin-top:3px; }
  .card-stats { display:flex;gap:12px;align-items:center;flex-wrap:wrap; }
  .stat { display:flex;align-items:center;gap:4px;font-size:12px;color:var(--ash-light); }
  .npc-tag { background:rgba(201,168,76,.15);color:var(--gold);border:1px solid rgba(201,168,76,.3);padding:2px 6px;border-radius:3px;font-size:10px;font-weight:600;letter-spacing:.1em; }
  .delete-btn { background:none;border:none;color:var(--stone-border);padding:4px;border-radius:3px;transition:color var(--transition);flex-shrink:0; }
  .delete-btn:hover { color:var(--crimson-light); }
  :global(.chevron) { color:var(--stone-border);align-self:center;margin-right:12px;flex-shrink:0; }
  .overlay { position:fixed;inset:0;background:rgba(0,0,0,.75);z-index:100;display:flex;align-items:center;justify-content:center;padding:20px; }
  .modal { background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:32px;width:100%;max-width:560px;max-height:90vh;overflow-y:auto; }
  .modal-title { font-family:var(--font-display);font-size:22px;font-weight:700;color:var(--parchment);margin-bottom:24px; }
  form { display:flex;flex-direction:column;gap:16px; }
  .form-row { display:grid;grid-template-columns:1fr 1fr;gap:12px; }
  .form-row-4 { display:grid;grid-template-columns:repeat(4,1fr);gap:12px; }
  .field { display:flex;flex-direction:column;gap:6px; }
  .field label { font-size:11px;font-weight:500;color:var(--ash-light);text-transform:uppercase;letter-spacing:.05em; }
  .section-label { font-size:11px;font-weight:600;color:var(--gold-dim);text-transform:uppercase;letter-spacing:.1em;border-bottom:1px solid var(--stone-border);padding-bottom:6px; }
  .ability-grid { display:grid;grid-template-columns:repeat(6,1fr);gap:8px; }
  .ability-field { display:flex;flex-direction:column;gap:4px; }
  .ability-field label { font-size:10px;font-weight:700;color:var(--ash);text-align:center;letter-spacing:.08em; }
  .ability-field input { text-align:center; }
  .error { font-size:13px;color:var(--crimson-light); }
  .modal-actions { display:flex;justify-content:flex-end;gap:10px;margin-top:8px; }
  .cancel-btn { background:var(--stone-mid);border:none;color:var(--ash-light);padding:10px 20px;border-radius:var(--radius);font-size:13px; }
  .cancel-btn:hover { color:var(--parchment); }
  .submit-btn { background:var(--crimson);border:none;color:var(--parchment);padding:10px 20px;border-radius:var(--radius);font-family:var(--font-display);font-size:13px;font-weight:600; }
  .submit-btn:hover { background:var(--crimson-light); }
</style>
