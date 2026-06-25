<script lang="ts">
  import { onMount } from 'svelte';
  import * as api from '../lib/api';
  import { navigate } from '../lib/router.svelte';
  import Icon from '../lib/Icon.svelte';

  let { id }: { id: string } = $props();

  const ABILITIES  = ['strength','dexterity','constitution','intelligence','wisdom','charisma'] as const;
  const CONDITIONS = ['Blinded','Charmed','Deafened','Exhaustion','Frightened','Grappled','Incapacitated','Invisible','Paralyzed','Petrified','Poisoned','Prone','Restrained','Stunned','Unconscious'];
  const SKILLS = [
    {name:'Acrobatics',ab:'dexterity',f:'skill_acrobatics'},{name:'Animal Handling',ab:'wisdom',f:'skill_animal_handling'},
    {name:'Arcana',ab:'intelligence',f:'skill_arcana'},{name:'Athletics',ab:'strength',f:'skill_athletics'},
    {name:'Deception',ab:'charisma',f:'skill_deception'},{name:'History',ab:'intelligence',f:'skill_history'},
    {name:'Insight',ab:'wisdom',f:'skill_insight'},{name:'Intimidation',ab:'charisma',f:'skill_intimidation'},
    {name:'Investigation',ab:'intelligence',f:'skill_investigation'},{name:'Medicine',ab:'wisdom',f:'skill_medicine'},
    {name:'Nature',ab:'intelligence',f:'skill_nature'},{name:'Perception',ab:'wisdom',f:'skill_perception'},
    {name:'Performance',ab:'charisma',f:'skill_performance'},{name:'Persuasion',ab:'charisma',f:'skill_persuasion'},
    {name:'Religion',ab:'intelligence',f:'skill_religion'},{name:'Sleight of Hand',ab:'dexterity',f:'skill_sleight_of_hand'},
    {name:'Stealth',ab:'dexterity',f:'skill_stealth'},{name:'Survival',ab:'wisdom',f:'skill_survival'},
  ] as const;
  const RACES       = ['Human','Elf','Dwarf','Halfling','Gnome','Half-Elf','Half-Orc','Tiefling','Dragonborn','Aasimar','Tabaxi','Kenku','Firbolg'];
  const CLASSES     = ['Barbarian','Bard','Cleric','Druid','Fighter','Monk','Paladin','Ranger','Rogue','Sorcerer','Warlock','Wizard','Artificer'];
  const ALIGNMENTS  = ['Lawful Good','Neutral Good','Chaotic Good','Lawful Neutral','True Neutral','Chaotic Neutral','Lawful Evil','Neutral Evil','Chaotic Evil'];
  const SP_ABILITIES  = ['none','strength','dexterity','constitution','intelligence','wisdom','charisma'];
  const SPELL_SCHOOLS = ['Abjuration','Conjuration','Divination','Enchantment','Evocation','Illusion','Necromancy','Transmutation'];

  const mod = (s: number) => Math.floor((s - 10) / 2);
  const fmt = (n: number) => n >= 0 ? `+${n}` : `${n}`;
  const pb  = (lvl: number) => Math.ceil(lvl / 4) + 1;

  let char      = $state<any>(null);
  let inventory = $state<any[]>([]);
  let spells    = $state<any[]>([]);
  let slots     = $state<any[]>([]);
  let loading   = $state(true);
  let tab       = $state<'stats'|'inventory'|'spells'|'notes'>('stats');
  let editing   = $state(false);
  let hpInput   = $state('');
  let hpMode    = $state<'damage'|'heal'|'temp'>('damage');

  // Edit form state
  let editInfo     = $state<any>({});
  let editLevel    = $state<any>({});
  let editAbi      = $state<any>({});
  let editSkills   = $state<any>({});
  let editTraining = $state<any>({});
  let editCurrency = $state<any>({});
  let saving       = $state('');
  let saved        = $state('');
  let editError    = $state('');

  // Inventory form
  let showAddItem = $state(false);
  let itemForm    = $state({name:'',quantity:1,weight:0,value:0,description:'',is_equipped:false,requires_attunement:false});

  // Spell form
  let showAddSpell = $state(false);
  let spellForm    = $state({name:'',level:1,school:'Evocation',casting_time:'1 action',range:'60 feet',components:'V, S',duration:'Instantaneous',description:'',is_prepared:true});

  // Notes
  let notes      = $state('');
  let notesSaved = $state(false);

  async function reload() {
    const [c, inv, sp, sl] = await Promise.all([
      api.getCharacter(id), api.listInventory(id), api.listSpells(id), api.listSpellSlots(id),
    ]);
    char = c; inventory = inv||[]; spells = sp||[]; slots = sl||[];
    notes = c?.notes || '';
    initEditState(c);
  }

  function initEditState(c: any) {
    if (!c) return;
    editInfo     = { name:c.name||'', race:c.race||'Human', background:c.background||'', alignment:c.alignment||'True Neutral', inspiration:c.inspiration||false, speed:c.speed||30, personality_traits:c.personality_traits||'', ideals:c.ideals||'', bonds:c.bonds||'', flaws:c.flaws||'' };
    editLevel    = { class:c.class||'Fighter', level:c.level||1, xp:c.xp||0, max_hp:c.max_hp||10, hit_dice_type:c.hit_dice_type||8, hit_dice_remaining:c.hit_dice_remaining||1 };
    editAbi      = { strength:c.strength||10,dexterity:c.dexterity||10,constitution:c.constitution||10,intelligence:c.intelligence||10,wisdom:c.wisdom||10,charisma:c.charisma||10,save_prof_strength:c.save_prof_strength||false,save_prof_dexterity:c.save_prof_dexterity||false,save_prof_constitution:c.save_prof_constitution||false,save_prof_intelligence:c.save_prof_intelligence||false,save_prof_wisdom:c.save_prof_wisdom||false,save_prof_charisma:c.save_prof_charisma||false };
    editSkills   = { skill_acrobatics:c.skill_acrobatics||0,skill_animal_handling:c.skill_animal_handling||0,skill_arcana:c.skill_arcana||0,skill_athletics:c.skill_athletics||0,skill_deception:c.skill_deception||0,skill_history:c.skill_history||0,skill_insight:c.skill_insight||0,skill_intimidation:c.skill_intimidation||0,skill_investigation:c.skill_investigation||0,skill_medicine:c.skill_medicine||0,skill_nature:c.skill_nature||0,skill_perception:c.skill_perception||0,skill_performance:c.skill_performance||0,skill_persuasion:c.skill_persuasion||0,skill_religion:c.skill_religion||0,skill_sleight_of_hand:c.skill_sleight_of_hand||0,skill_stealth:c.skill_stealth||0,skill_survival:c.skill_survival||0 };
    editTraining = { armor_class:c.armor_class||10,attunement_slots:c.attunement_slots||3,spellcasting_ability:c.spellcasting_ability||'none',training_armor:(c.training_armor||[]).join(', '),training_weapons:(c.training_weapons||[]).join(', '),training_tools:(c.training_tools||[]).join(', '),training_languages:(c.training_languages||[]).join(', '),resistances:(c.resistances||[]).join(', '),vulnerabilities:(c.vulnerabilities||[]).join(', '),immunities:(c.immunities||[]).join(', ') };
    editCurrency = { copper:c.copper||0,silver:c.silver||0,electrum:c.electrum||0,gold:c.gold||0,platinum:c.platinum||0 };
  }

  onMount(() => reload().finally(() => loading = false));

  async function handleHP() {
    const n = parseInt(hpInput); if (!n || n <= 0) return;
    if (hpMode === 'damage') await api.applyDamage(id, n);
    else if (hpMode === 'heal') await api.applyHeal(id, n);
    else await api.addTempHP(id, n);
    hpInput = ''; await reload();
  }

  async function toggleCond(c: string) {
    const next = char.conditions.includes(c)
      ? char.conditions.filter((x:string) => x !== c)
      : [...char.conditions, c];
    await api.updateConditions(id, next); await reload();
  }

  async function handleRest(type: 'long'|'short') {
    if (type === 'long') await api.longRest(id);
    else {
      const healed = Math.max(1, char.hit_dice_type + mod(char.constitution));
      await api.shortRest(id, Math.max(0, char.hit_dice_remaining-1), Math.min(char.max_hp, char.current_hp+healed));
    }
    await reload();
  }

  const splitList = (s: string) => s.split(',').map((x:string) => x.trim()).filter(Boolean);

  async function saveSection(section: string, fn: () => Promise<any>) {
    saving = section; editError = '';
    try { await fn(); await reload(); saved = section; setTimeout(() => saved='', 2000); }
    catch(e:any) { editError = e?.data?.error || 'Failed to save'; }
    finally { saving = ''; }
  }

  async function addItem(e: Event) {
    e.preventDefault(); await api.createInventoryItem(id, itemForm);
    itemForm = {name:'',quantity:1,weight:0,value:0,description:'',is_equipped:false,requires_attunement:false};
    showAddItem = false; await reload();
  }

  async function addSpell(e: Event) {
    e.preventDefault(); await api.createSpell(id, spellForm);
    showAddSpell = false; await reload();
  }

  async function saveNotes() {
    await api.updateCharacterInfo(id, { notes });
    notesSaved = true; setTimeout(() => notesSaved=false, 2000); await reload();
  }

  $derived: var profBonus = char ? pb(char.level) : 2;
  $derived: var hpPct     = char ? Math.max(0, Math.min(100, (char.current_hp / char.max_hp) * 100)) : 0;
  $derived: var hpColor   = hpPct > 50 ? '#27ae60' : hpPct > 25 ? '#f39c12' : '#c0392b';
  $derived: var byLevel   = Array.from({length:10}, (_,i) => ({
    level: i,
    spells: spells.filter((s:any) => s.level === i),
    slot: i > 0 ? slots.find((s:any) => s.spell_level === i) : undefined,
  })).filter(g => g.spells.length > 0 || (g.slot && g.slot.total > 0));
</script>

{#if loading}
  <div class="loading">Loading…</div>
{:else if !char}
  <div class="loading">Character not found.</div>
{:else}
<div class="page">

  <!-- Header -->
  <div class="sheet-header">
    <button class="back" onclick={() => navigate('/characters')}>
      <Icon name="chevronLeft" size={16}/> Characters
    </button>
    <div class="header-info">
      <h1 class="char-name">{char.name}</h1>
      <p class="char-meta">Level {char.level} {char.race} {char.class} · {char.background} · {char.alignment}</p>
    </div>
    <div class="header-badges">
      {#if char.inspiration}<span class="badge gold">✦ Inspired</span>{/if}
      {#if char.is_npc}<span class="badge crimson">NPC</span>{/if}
      <span class="badge">{char.xp?.toLocaleString()} XP</span>
      <button class="edit-btn" class:active={editing} onclick={() => editing = !editing}>
        {editing ? '✓ Done Editing' : '✎ Edit'}
      </button>
    </div>
  </div>

  <!-- Combat Bar -->
  <div class="combat-bar">
    <div class="hp-section">
      <div class="hp-label">HIT POINTS</div>
      <div class="hp-track"><div class="hp-bar" style="width:{hpPct}%;background:{hpColor}"></div></div>
      <div class="hp-numbers">
        <span style="color:{hpColor};font-weight:700;font-size:20px">{char.current_hp}</span>
        <span class="hp-max">/ {char.max_hp}</span>
        {#if char.temp_hp > 0}<span class="temp-hp">+{char.temp_hp} temp</span>{/if}
      </div>
      <div class="hp-controls">
        <div class="hp-mode-row">
          {#each ['damage','heal','temp'] as m}
            <button class="hp-mode-btn" class:active={hpMode===m} onclick={() => hpMode=m as any}>
              {m==='damage' ? '⚔ Dmg' : m==='heal' ? '❤ Heal' : '🛡 Temp'}
            </button>
          {/each}
        </div>
        <div class="hp-input-row">
          <input class="hp-input" type="number" min="1" bind:value={hpInput} placeholder="0"
            onkeydown={(e) => e.key==='Enter' && handleHP()}/>
          <button class="hp-apply" onclick={handleHP}>Apply</button>
        </div>
      </div>
    </div>

    <div class="combat-stats">
      {#each [{l:'AC',v:char.armor_class},{l:'INIT',v:fmt(mod(char.dexterity))},{l:'SPEED',v:`${char.speed}ft`},{l:'PROF',v:`+${profBonus}`},{l:`d${char.hit_dice_type}`,v:`${char.hit_dice_remaining} left`}] as s}
        <div class="combat-stat">
          <div class="cs-val">{s.v}</div>
          <div class="cs-label">{s.l}</div>
        </div>
      {/each}
    </div>

    <div class="rest-section">
      <button class="rest-btn" onclick={() => handleRest('short')}><Icon name="moon" size={14}/> Short Rest</button>
      <button class="rest-btn" onclick={() => handleRest('long')}><Icon name="sun" size={14}/> Long Rest</button>
    </div>

    {#if char.current_hp <= 0}
      <div class="death-saves">
        <div class="death-label">DEATH SAVES</div>
        {#each ['Successes','Failures'] as label}
          <div class="death-row">
            <span class="death-name">{label}</span>
            {#each [0,1,2] as i}
              {@const count = label==='Successes' ? char.death_save_successes : char.death_save_failures}
              <button class="death-dot"
                class:success={label==='Successes' && i<count}
                class:failure={label==='Failures' && i<count}
                onclick={() => api.recordDeathSave(id, label==='Successes').then(reload)}>
              </button>
            {/each}
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Tabs -->
  <div class="tabs">
    {#each [['stats','shield','Stats'],['inventory','backpack','Inventory'],['spells','book','Spells'],['notes','star','Notes']] as [t,icon,label]}
      <button class="tab-btn" class:active={tab===t} onclick={() => tab=t as any}>
        <Icon name={icon} size={15}/> {label}
      </button>
    {/each}
  </div>

  <!-- Content -->
  <div class="content">

    {#if tab === 'stats'}
      {#if editing}
        <!-- EDIT MODE -->
        <div class="edit-grid">
          {#if editError}<p class="edit-error">{editError}</p>{/if}

          <!-- Info -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Character Info</h3>
              <button class="save-btn" class:saved={saved==='info'} disabled={saving==='info'}
                onclick={() => saveSection('info', () => api.updateCharacterInfo(id, editInfo))}>
                {saving==='info' ? 'Saving…' : saved==='info' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <div class="edit-form-grid">
              <div class="field"><label>Name</label><input bind:value={editInfo.name}/></div>
              <div class="field"><label>Race</label><select bind:value={editInfo.race}>{#each RACES as r}<option>{r}</option>{/each}</select></div>
              <div class="field"><label>Background</label><input bind:value={editInfo.background}/></div>
              <div class="field"><label>Alignment</label><select bind:value={editInfo.alignment}>{#each ALIGNMENTS as a}<option>{a}</option>{/each}</select></div>
              <div class="field"><label>Speed (ft)</label><input type="number" bind:value={editInfo.speed}/></div>
              <div class="field" style="align-self:center">
                <label class="check-label"><input type="checkbox" bind:checked={editInfo.inspiration}/> Inspiration</label>
              </div>
              <div class="field full"><label>Personality Traits</label><textarea rows="2" bind:value={editInfo.personality_traits}></textarea></div>
              <div class="field full"><label>Ideals</label><textarea rows="2" bind:value={editInfo.ideals}></textarea></div>
              <div class="field full"><label>Bonds</label><textarea rows="2" bind:value={editInfo.bonds}></textarea></div>
              <div class="field full"><label>Flaws</label><textarea rows="2" bind:value={editInfo.flaws}></textarea></div>
            </div>
          </div>

          <!-- Level & Class -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Level & Class</h3>
              <button class="save-btn" class:saved={saved==='level'} disabled={saving==='level'}
                onclick={() => saveSection('level', () => api.updateLevel(id, editLevel))}>
                {saving==='level' ? 'Saving…' : saved==='level' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <div class="edit-form-grid">
              <div class="field"><label>Class</label><select bind:value={editLevel.class}>{#each CLASSES as c}<option>{c}</option>{/each}</select></div>
              <div class="field"><label>Level</label><input type="number" min="1" max="20" bind:value={editLevel.level}/></div>
              <div class="field"><label>XP</label><input type="number" min="0" bind:value={editLevel.xp}/></div>
              <div class="field"><label>Max HP</label><input type="number" min="1" bind:value={editLevel.max_hp}/></div>
              <div class="field"><label>Hit Die (d?)</label><input type="number" min="4" bind:value={editLevel.hit_dice_type}/></div>
              <div class="field"><label>Hit Dice Remaining</label><input type="number" min="0" bind:value={editLevel.hit_dice_remaining}/></div>
            </div>
          </div>

          <!-- Ability Scores -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Ability Scores & Saves</h3>
              <button class="save-btn" class:saved={saved==='abilities'} disabled={saving==='abilities'}
                onclick={() => saveSection('abilities', () => api.updateAbilityScores(id, editAbi))}>
                {saving==='abilities' ? 'Saving…' : saved==='abilities' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <div class="ability-edit-grid">
              {#each ABILITIES as ab}
                <div class="ability-edit-box">
                  <label>{ab.slice(0,3).toUpperCase()}</label>
                  <input type="number" min="1" max="30" bind:value={editAbi[ab]}/>
                  <label class="check-label">
                    <input type="checkbox" bind:checked={editAbi[`save_prof_${ab}`]}/> Save
                  </label>
                </div>
              {/each}
            </div>
          </div>

          <!-- Skills -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Skills</h3>
              <button class="save-btn" class:saved={saved==='skills'} disabled={saving==='skills'}
                onclick={() => saveSection('skills', () => api.updateSkills(id, editSkills))}>
                {saving==='skills' ? 'Saving…' : saved==='skills' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <p class="edit-hint">0 = none · 1 = proficient · 2 = expertise</p>
            <div class="skill-edit-grid">
              {#each SKILLS as sk}
                <div class="skill-edit-row">
                  <span class="skill-edit-name">{sk.name}</span>
                  <select class="skill-edit-select" bind:value={editSkills[sk.f]}>
                    <option value={0}>—</option>
                    <option value={1}>Prof</option>
                    <option value={2}>Expert</option>
                  </select>
                </div>
              {/each}
            </div>
          </div>

          <!-- Training & Defenses -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Training & Defenses</h3>
              <button class="save-btn" class:saved={saved==='training'} disabled={saving==='training'}
                onclick={() => saveSection('training', () => api.updateTraining(id, {
                  ...editTraining,
                  training_armor:     splitList(editTraining.training_armor),
                  training_weapons:   splitList(editTraining.training_weapons),
                  training_tools:     splitList(editTraining.training_tools),
                  training_languages: splitList(editTraining.training_languages),
                  resistances:        splitList(editTraining.resistances),
                  vulnerabilities:    splitList(editTraining.vulnerabilities),
                  immunities:         splitList(editTraining.immunities),
                }))}>
                {saving==='training' ? 'Saving…' : saved==='training' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <p class="edit-hint">Separate multiple values with commas</p>
            <div class="edit-form-grid">
              <div class="field"><label>Armor Class</label><input type="number" min="1" bind:value={editTraining.armor_class}/></div>
              <div class="field"><label>Attunement Slots</label><input type="number" min="0" max="5" bind:value={editTraining.attunement_slots}/></div>
              <div class="field"><label>Spellcasting Ability</label>
                <select bind:value={editTraining.spellcasting_ability}>
                  {#each SP_ABILITIES as s}<option value={s}>{s==='none' ? 'None' : s[0].toUpperCase()+s.slice(1)}</option>{/each}
                </select>
              </div>
              <div class="field"><label>Armor Profs</label><input bind:value={editTraining.training_armor} placeholder="Light, Medium…"/></div>
              <div class="field"><label>Weapon Profs</label><input bind:value={editTraining.training_weapons} placeholder="Simple, Martial…"/></div>
              <div class="field"><label>Tool Profs</label><input bind:value={editTraining.training_tools}/></div>
              <div class="field"><label>Languages</label><input bind:value={editTraining.training_languages} placeholder="Common, Elvish…"/></div>
              <div class="field"><label>Resistances</label><input bind:value={editTraining.resistances}/></div>
              <div class="field"><label>Vulnerabilities</label><input bind:value={editTraining.vulnerabilities}/></div>
              <div class="field"><label>Immunities</label><input bind:value={editTraining.immunities}/></div>
            </div>
          </div>

          <!-- Currency -->
          <div class="edit-panel">
            <div class="edit-panel-header">
              <h3 class="panel-title">Currency</h3>
              <button class="save-btn" class:saved={saved==='currency'} disabled={saving==='currency'}
                onclick={() => saveSection('currency', () => api.updateCurrency(id, editCurrency))}>
                {saving==='currency' ? 'Saving…' : saved==='currency' ? '✓ Saved' : 'Save'}
              </button>
            </div>
            <div class="currency-edit-grid">
              {#each ['copper','silver','electrum','gold','platinum'] as coin}
                <div class="field">
                  <label>{coin[0].toUpperCase()+coin.slice(1)}</label>
                  <input type="number" min="0" bind:value={editCurrency[coin]}/>
                </div>
              {/each}
            </div>
          </div>
        </div>

      {:else}
        <!-- VIEW MODE -->
        <div class="stats-grid">
          <div class="panel">
            <h3 class="panel-title">Ability Scores</h3>
            <div class="ability-grid">
              {#each ABILITIES as ab}
                {@const score = char[ab]}
                {@const m = mod(score)}
                {@const saveProf = char[`save_prof_${ab}`]}
                <div class="ability-box">
                  <div class="abi-name">{ab.slice(0,3).toUpperCase()}</div>
                  <div class="abi-score">{score}</div>
                  <div class="abi-mod">{fmt(m)}</div>
                  <div class="save-badge" class:prof={saveProf}>{fmt(m+(saveProf?profBonus:0))} save</div>
                </div>
              {/each}
            </div>
          </div>

          <div class="panel">
            <h3 class="panel-title">Skills</h3>
            <div class="skill-list">
              {#each SKILLS as sk}
                {@const val = char[sk.f]||0}
                {@const total = mod(char[sk.ab]) + val*profBonus}
                <div class="skill-row">
                  <div class="skill-dot" class:proficient={val===1} class:expert={val===2}></div>
                  <span class="skill-name">{sk.name}</span>
                  <span class="skill-mod">{fmt(total)}</span>
                  <span class="skill-ab">{sk.ab.slice(0,3)}</span>
                </div>
              {/each}
            </div>
          </div>

          <div class="panel">
            <h3 class="panel-title">Conditions</h3>
            <div class="condition-grid">
              {#each CONDITIONS as c}
                <button class="cond-pill" class:active={char.conditions?.includes(c)} onclick={() => toggleCond(c)}>{c}</button>
              {/each}
            </div>
          </div>

          <div class="panel">
            <h3 class="panel-title">Currency</h3>
            <div class="currency-row">
              {#each [['CP',char.copper,'#b5651d'],['SP',char.silver,'#c0c0c0'],['EP',char.electrum,'#b8b8ff'],['GP',char.gold,'#ffd700'],['PP',char.platinum,'#e5e4e2']] as [l,v,color]}
                <div class="coin">
                  <div class="coin-circle" style="border-color:{color};color:{color}">{l}</div>
                  <div class="coin-val">{v}</div>
                </div>
              {/each}
            </div>
          </div>

          <div class="panel panel-full">
            <h3 class="panel-title">Character Traits</h3>
            <div class="traits-grid">
              {#each [['Personality',char.personality_traits],['Ideals',char.ideals],['Bonds',char.bonds],['Flaws',char.flaws]] as [l,v]}
                <div class="trait-box">
                  <div class="trait-label">{l}</div>
                  <div class="trait-text">{v||'—'}</div>
                </div>
              {/each}
            </div>
          </div>

          <div class="panel panel-full">
            <h3 class="panel-title">Proficiencies & Training</h3>
            <div class="traits-grid">
              {#each [['Armor',char.training_armor],['Weapons',char.training_weapons],['Tools',char.training_tools],['Languages',char.training_languages]] as [l,v]}
                <div class="trait-box">
                  <div class="trait-label">{l}</div>
                  <div class="trait-text">{(v||[]).join(', ')||'—'}</div>
                </div>
              {/each}
            </div>
          </div>
        </div>
      {/if}

    {:else if tab === 'inventory'}
      <div class="tab-content">
        <div class="tab-header">
          <div class="tab-meta">
            <span>{inventory.reduce((s:number,i:any) => s+i.weight*i.quantity, 0)} lbs</span>
            <span>{inventory.filter((i:any) => i.is_attuned).length}/{char.attunement_slots} attuned</span>
          </div>
          <button class="add-btn" onclick={() => showAddItem=true}><Icon name="plus" size={14}/> Add Item</button>
        </div>
        {#if showAddItem}
          <form class="add-form" onsubmit={addItem}>
            <div class="form-row">
              <div class="field"><label>Name *</label><input bind:value={itemForm.name} required/></div>
              <div class="field"><label>Qty</label><input type="number" min="1" bind:value={itemForm.quantity}/></div>
              <div class="field"><label>Weight (lb)</label><input type="number" min="0" bind:value={itemForm.weight}/></div>
              <div class="field"><label>Value (gp)</label><input type="number" min="0" bind:value={itemForm.value}/></div>
            </div>
            <div class="field"><label>Description</label><input bind:value={itemForm.description}/></div>
            <label class="check-label">
              <input type="checkbox" bind:checked={itemForm.requires_attunement}/> Requires Attunement
            </label>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" onclick={() => showAddItem=false}>Cancel</button>
              <button type="submit" class="submit-btn">Add</button>
            </div>
          </form>
        {/if}
        <div class="item-list">
          {#each inventory as item}
            <div class="item-row" class:equipped={item.is_equipped}>
              <div class="item-main">
                <div class="item-name">
                  {item.name}
                  {#if item.is_equipped}<span class="equipped-tag">Equipped</span>{/if}
                  {#if item.is_attuned}<span class="attuned-tag">Attuned</span>{/if}
                </div>
                {#if item.description}<div class="item-desc">{item.description}</div>{/if}
                <div class="item-meta">
                  <span>Qty: {item.quantity}</span>
                  {#if item.weight>0}<span>{item.weight}lb</span>{/if}
                  {#if item.value>0}<span>{item.value}gp</span>{/if}
                </div>
              </div>
              <div class="item-actions">
                <button class="item-btn"
                  onclick={() => api.updateInventoryItem(id, item.id, {...item, is_equipped:!item.is_equipped}).then(reload)}>
                  {item.is_equipped ? '⚔' : '○'}
                </button>
                {#if item.requires_attunement}
                  <button class="item-btn"
                    onclick={() => (item.is_attuned ? api.unattuneItem : api.attuneItem)(id, item.id).then(reload)}>
                    ✦
                  </button>
                {/if}
                <button class="item-btn-danger" onclick={() => api.deleteInventoryItem(id, item.id).then(reload)}>
                  <Icon name="trash" size={13}/>
                </button>
              </div>
            </div>
          {/each}
          {#if inventory.length===0}<div class="empty-state">No items yet.</div>{/if}
        </div>
      </div>

    {:else if tab === 'spells'}
      <div class="tab-content">
        <div class="tab-header">
          <div class="tab-meta"><span>{spells.filter((s:any) => s.is_prepared).length} prepared</span></div>
          <button class="add-btn" onclick={() => showAddSpell=true}><Icon name="plus" size={14}/> Add Spell</button>
        </div>
        {#if showAddSpell}
          <form class="add-form" onsubmit={addSpell}>
            <div class="form-row">
              <div class="field"><label>Name *</label><input bind:value={spellForm.name} required/></div>
              <div class="field"><label>Level</label>
                <select bind:value={spellForm.level}>
                  {#each [0,1,2,3,4,5,6,7,8,9] as l}
                    <option value={l}>{l===0 ? 'Cantrip' : `Level ${l}`}</option>
                  {/each}
                </select>
              </div>
              <div class="field"><label>School</label>
                <select bind:value={spellForm.school}>
                  {#each SPELL_SCHOOLS as s}<option>{s}</option>{/each}
                </select>
              </div>
            </div>
            <div class="form-row">
              <div class="field"><label>Cast Time</label><input bind:value={spellForm.casting_time}/></div>
              <div class="field"><label>Range</label><input bind:value={spellForm.range}/></div>
              <div class="field"><label>Duration</label><input bind:value={spellForm.duration}/></div>
            </div>
            <div class="field"><label>Components</label><input bind:value={spellForm.components}/></div>
            <div class="field"><label>Description</label><textarea rows="3" bind:value={spellForm.description}></textarea></div>
            <div class="modal-actions">
              <button type="button" class="cancel-btn" onclick={() => showAddSpell=false}>Cancel</button>
              <button type="submit" class="submit-btn">Add Spell</button>
            </div>
          </form>
        {/if}
        {#each byLevel as group}
          <div class="spell-group">
            <div class="spell-group-header">
              <span>{group.level===0 ? 'Cantrips' : `Level ${group.level}`}</span>
              {#if group.slot}
                <div class="slot-track">
                  {#each Array.from({length:group.slot.total}, (_,i) => i) as i}
                    <button class="slot-dot" class:used={i<group.slot.used}
                      onclick={() => api.useSpellSlot(id, group.level).then(reload)}>
                    </button>
                  {/each}
                  <span class="slot-label">{group.slot.total-group.slot.used}/{group.slot.total}</span>
                </div>
              {/if}
            </div>
            {#each group.spells as spell}
              <div class="spell-row" class:unprepared={!spell.is_prepared && spell.level>0}>
                <button class="prepare-btn" class:prepared={spell.is_prepared}
                  onclick={() => api.toggleSpellPrepared(id, spell.id).then(reload)}>
                  <Icon name="check" size={12}/>
                </button>
                <div class="spell-main">
                  <div class="spell-name">{spell.name}</div>
                  <div class="spell-meta">{spell.school} · {spell.casting_time} · {spell.range}</div>
                  {#if spell.description}<div class="spell-desc">{spell.description}</div>{/if}
                </div>
                <button class="item-btn-danger" onclick={() => api.deleteSpell(id, spell.id).then(reload)}>
                  <Icon name="trash" size={13}/>
                </button>
              </div>
            {/each}
          </div>
        {/each}
        {#if byLevel.length===0}<div class="empty-state">No spells yet.</div>{/if}
      </div>

    {:else if tab === 'notes'}
      <div class="tab-content">
        <div class="tab-header">
          <span></span>
          <button class="add-btn" onclick={saveNotes}>{notesSaved ? '✓ Saved' : 'Save Notes'}</button>
        </div>
        <textarea class="notes-area" bind:value={notes} placeholder="Session notes, backstory, important NPCs…"></textarea>
      </div>
    {/if}

  </div>
</div>
{/if}

<style>
.loading { display:flex;align-items:center;justify-content:center;height:300px;color:var(--ash); }
.page { min-height:100vh; }
.sheet-header { display:flex;align-items:center;gap:20px;padding:20px 32px;background:var(--stone);border-bottom:1px solid var(--stone-border);flex-wrap:wrap; }
.back { display:flex;align-items:center;gap:4px;background:none;border:none;color:var(--ash);font-size:13px;transition:color var(--transition); }
.back:hover { color:var(--parchment); }
.header-info { flex:1; }
.char-name { font-family:var(--font-display);font-size:24px;font-weight:700;color:var(--parchment); }
.char-meta { font-size:13px;color:var(--ash);margin-top:3px; }
.header-badges { display:flex;gap:8px;flex-wrap:wrap;align-items:center; }
.badge { font-size:11px;padding:3px 8px;border-radius:12px;background:var(--stone-mid);color:var(--ash-light); }
.badge.gold { background:rgba(201,168,76,.2);color:var(--gold); }
.badge.crimson { background:rgba(139,26,26,.2);color:var(--crimson-light); }
.edit-btn { display:flex;align-items:center;gap:6px;background:var(--stone);border:1px solid var(--stone-border);color:var(--ash-light);padding:7px 14px;border-radius:var(--radius);font-size:13px;font-weight:500;transition:all var(--transition); }
.edit-btn:hover { border-color:var(--gold-dim);color:var(--parchment); }
.edit-btn.active { background:rgba(201,168,76,.15);border-color:var(--gold-dim);color:var(--gold);font-weight:600; }
.combat-bar { padding:20px 32px;background:var(--stone-mid);border-bottom:1px solid var(--stone-border);display:flex;gap:32px;align-items:flex-start;flex-wrap:wrap; }
.hp-section { min-width:220px; }
.hp-label { font-size:10px;font-weight:700;letter-spacing:.12em;color:var(--ash);margin-bottom:6px; }
.hp-track { height:6px;background:var(--stone-border);border-radius:3px;margin-bottom:8px;overflow:hidden; }
.hp-bar { height:100%;border-radius:3px;transition:width .4s ease,background .4s ease; }
.hp-numbers { display:flex;align-items:baseline;gap:6px;margin-bottom:12px; }
.hp-max { color:var(--ash);font-size:14px; }
.temp-hp { font-size:13px;color:#3498db; }
.hp-controls { display:flex;flex-direction:column;gap:6px; }
.hp-mode-row { display:flex;gap:4px; }
.hp-mode-btn { background:var(--stone);border:1px solid var(--stone-border);color:var(--ash);padding:5px 10px;border-radius:var(--radius-sm);font-size:11px;transition:all var(--transition); }
.hp-mode-btn.active { background:var(--stone-light);border-color:var(--gold-dim);color:var(--parchment); }
.hp-input-row { display:flex;gap:6px; }
.hp-input { width:70px;text-align:center;padding:6px 8px; }
.hp-apply { background:var(--crimson);border:none;color:var(--parchment);padding:6px 14px;border-radius:var(--radius-sm);font-size:13px; }
.hp-apply:hover { background:var(--crimson-light); }
.combat-stats { display:flex;gap:16px;flex-wrap:wrap;align-items:center; }
.combat-stat { text-align:center; }
.cs-val { font-family:var(--font-display);font-size:22px;font-weight:700;color:var(--parchment); }
.cs-label { font-size:10px;color:var(--ash);letter-spacing:.08em;margin-top:2px; }
.rest-section { display:flex;flex-direction:column;gap:6px; }
.rest-btn { display:flex;align-items:center;gap:6px;background:var(--stone);border:1px solid var(--stone-border);color:var(--ash-light);padding:8px 14px;border-radius:var(--radius);font-size:12px;transition:all var(--transition); }
.rest-btn:hover { border-color:var(--gold-dim);color:var(--parchment); }
.death-saves { }
.death-label { font-size:10px;font-weight:700;letter-spacing:.1em;color:var(--crimson-light);margin-bottom:8px; }
.death-row { display:flex;align-items:center;gap:8px;margin-bottom:6px; }
.death-name { font-size:11px;color:var(--ash);width:70px; }
.death-dot { width:18px;height:18px;border-radius:50%;border:2px solid var(--stone-border);background:transparent;transition:all var(--transition); }
.death-dot.success { background:var(--emerald);border-color:var(--emerald); }
.death-dot.failure { background:var(--crimson);border-color:var(--crimson); }
.tabs { display:flex;gap:2px;padding:0 32px;background:var(--stone);border-bottom:1px solid var(--stone-border); }
.tab-btn { display:flex;align-items:center;gap:6px;background:none;border:none;color:var(--ash);padding:14px 18px;font-size:13px;font-weight:500;border-bottom:2px solid transparent;transition:all var(--transition); }
.tab-btn:hover { color:var(--parchment); }
.tab-btn.active { color:var(--gold);border-bottom-color:var(--gold); }
.content { padding:24px 32px; }
.stats-grid { display:grid;grid-template-columns:auto 1fr 1fr;gap:20px; }
.panel { background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:20px; }
.panel-full { grid-column:1/-1; }
.panel-title { font-family:var(--font-display);font-size:13px;font-weight:600;color:var(--gold);letter-spacing:.08em;text-transform:uppercase;margin-bottom:16px; }
.ability-grid { display:grid;grid-template-columns:repeat(3,1fr);gap:10px; }
.ability-box { display:flex;flex-direction:column;align-items:center;background:var(--stone-mid);border:1px solid var(--stone-border);border-radius:var(--radius);padding:10px 8px; }
.abi-name { font-size:9px;font-weight:700;letter-spacing:.12em;color:var(--ash);margin-bottom:4px; }
.abi-score { font-family:var(--font-display);font-size:24px;font-weight:700;color:var(--parchment);line-height:1; }
.abi-mod { font-size:16px;font-weight:600;color:var(--gold); }
.save-badge { font-size:10px;color:var(--ash);margin-top:4px; }
.save-badge.prof { color:var(--gold-dim); }
.skill-list { display:flex;flex-direction:column;gap:3px; }
.skill-row { display:flex;align-items:center;gap:8px;padding:4px 6px;border-radius:3px; }
.skill-row:hover { background:var(--stone-mid); }
.skill-dot { width:10px;height:10px;border-radius:50%;border:2px solid var(--stone-border);flex-shrink:0; }
.skill-dot.proficient { background:var(--gold-dim);border-color:var(--gold-dim); }
.skill-dot.expert { background:var(--gold);border-color:var(--gold); }
.skill-name { flex:1;font-size:13px;color:var(--ash-light); }
.skill-mod { font-weight:600;font-size:13px;color:var(--parchment);width:30px;text-align:right; }
.skill-ab { font-size:10px;color:var(--ash);width:24px;text-align:center; }
.condition-grid { display:flex;flex-wrap:wrap;gap:6px; }
.cond-pill { background:var(--stone-mid);border:1px solid var(--stone-border);color:var(--ash);padding:5px 10px;border-radius:12px;font-size:12px;transition:all var(--transition); }
.cond-pill:hover { border-color:var(--crimson);color:var(--parchment); }
.cond-pill.active { background:rgba(139,26,26,.2);border-color:var(--crimson);color:var(--crimson-light); }
.currency-row { display:flex;gap:12px;justify-content:center;flex-wrap:wrap; }
.coin { display:flex;flex-direction:column;align-items:center;gap:6px; }
.coin-circle { width:40px;height:40px;border-radius:50%;border:2px solid;display:flex;align-items:center;justify-content:center;font-size:10px;font-weight:700; }
.coin-val { font-size:15px;font-weight:600;color:var(--parchment); }
.traits-grid { display:grid;grid-template-columns:1fr 1fr;gap:16px; }
.trait-box { background:var(--stone-mid);border-radius:var(--radius);padding:14px; }
.trait-label { font-size:10px;font-weight:700;color:var(--gold-dim);letter-spacing:.1em;text-transform:uppercase;margin-bottom:6px; }
.trait-text { font-family:var(--font-body);font-size:15px;color:var(--parchment);line-height:1.5; }
.edit-grid { display:flex;flex-direction:column;gap:20px; }
.edit-error { font-size:13px;color:var(--crimson-light);background:rgba(139,26,26,.15);border:1px solid rgba(139,26,26,.3);border-radius:var(--radius-sm);padding:8px 12px; }
.edit-panel { background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:20px; }
.edit-panel-header { display:flex;align-items:center;justify-content:space-between;margin-bottom:16px; }
.edit-hint { font-size:11px;color:var(--ash);margin-bottom:12px; }
.save-btn { background:var(--crimson);border:none;color:var(--parchment);padding:7px 16px;border-radius:var(--radius);font-family:var(--font-display);font-size:12px;font-weight:600;letter-spacing:.04em;transition:background var(--transition); }
.save-btn:hover:not(:disabled) { background:var(--crimson-light); }
.save-btn:disabled { opacity:.5;cursor:not-allowed; }
.save-btn.saved { background:var(--emerald); }
.edit-form-grid { display:grid;grid-template-columns:1fr 1fr;gap:12px; }
.field { display:flex;flex-direction:column;gap:5px; }
.field label { font-size:11px;color:var(--ash);text-transform:uppercase;letter-spacing:.05em; }
.field textarea { resize:vertical; }
.field.full { grid-column:1/-1; }
.check-label { display:flex;align-items:center;gap:6px;font-size:13px;color:var(--ash-light);cursor:pointer;text-transform:none;letter-spacing:0; }
.check-label input { width:auto; }
.ability-edit-grid { display:grid;grid-template-columns:repeat(3,1fr);gap:12px; }
.ability-edit-box { display:flex;flex-direction:column;gap:6px;background:var(--stone-mid);border-radius:var(--radius);padding:10px; }
.ability-edit-box label { font-size:10px;font-weight:700;color:var(--ash);text-transform:uppercase;letter-spacing:.1em; }
.ability-edit-box input[type="number"] { text-align:center;font-family:var(--font-display);font-size:18px;font-weight:700; }
.skill-edit-grid { display:grid;grid-template-columns:1fr 1fr;gap:6px; }
.skill-edit-row { display:flex;align-items:center;justify-content:space-between;gap:8px;background:var(--stone-mid);border-radius:var(--radius-sm);padding:6px 10px; }
.skill-edit-name { font-size:13px;color:var(--ash-light);flex:1; }
.skill-edit-select { width:auto;padding:4px 6px;font-size:12px; }
.currency-edit-grid { display:grid;grid-template-columns:repeat(5,1fr);gap:10px; }
.tab-content { }
.tab-header { display:flex;align-items:center;justify-content:space-between;margin-bottom:16px; }
.tab-meta { display:flex;gap:16px;font-size:13px;color:var(--ash); }
.add-btn { display:flex;align-items:center;gap:6px;background:var(--stone);border:1px solid var(--stone-border);color:var(--ash-light);padding:8px 14px;border-radius:var(--radius);font-size:13px;transition:all var(--transition); }
.add-btn:hover { border-color:var(--gold-dim);color:var(--parchment); }
.add-form { background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:20px;margin-bottom:20px;display:flex;flex-direction:column;gap:12px; }
.form-row { display:flex;gap:12px;flex-wrap:wrap; }
.form-row > * { flex:1;min-width:120px; }
.modal-actions { display:flex;justify-content:flex-end;gap:8px; }
.cancel-btn { background:var(--stone-mid);border:none;color:var(--ash-light);padding:8px 16px;border-radius:var(--radius);font-size:13px; }
.submit-btn { background:var(--crimson);border:none;color:var(--parchment);padding:8px 16px;border-radius:var(--radius);font-family:var(--font-display);font-size:13px;font-weight:600; }
.submit-btn:hover { background:var(--crimson-light); }
.empty-state { text-align:center;padding:40px;color:var(--ash);font-size:14px; }
.item-list { display:flex;flex-direction:column;gap:6px; }
.item-row { display:flex;align-items:center;gap:12px;background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius);padding:12px 16px; }
.item-row.equipped { border-color:var(--gold-dim); }
.item-main { flex:1; }
.item-name { font-size:14px;font-weight:500;color:var(--parchment);display:flex;align-items:center;gap:8px; }
.equipped-tag { font-size:10px;background:rgba(201,168,76,.15);color:var(--gold);border:1px solid rgba(201,168,76,.3);padding:1px 6px;border-radius:3px; }
.attuned-tag { font-size:10px;background:rgba(52,152,219,.15);color:#3498db;border:1px solid rgba(52,152,219,.3);padding:1px 6px;border-radius:3px; }
.item-desc { font-size:12px;color:var(--ash);margin-top:3px; }
.item-meta { display:flex;gap:10px;font-size:11px;color:var(--ash);margin-top:4px; }
.item-actions { display:flex;gap:6px; }
.item-btn { background:none;border:1px solid var(--stone-border);color:var(--ash);padding:5px 9px;border-radius:var(--radius-sm);font-size:13px;transition:all var(--transition); }
.item-btn:hover { border-color:var(--gold-dim);color:var(--gold); }
.item-btn-danger { background:none;border:none;color:var(--stone-border);padding:5px;transition:color var(--transition); }
.item-btn-danger:hover { color:var(--crimson-light); }
.spell-group { margin-bottom:20px; }
.spell-group-header { display:flex;align-items:center;justify-content:space-between;padding:8px 0;border-bottom:1px solid var(--stone-border);margin-bottom:8px;font-family:var(--font-display);font-size:13px;font-weight:600;color:var(--gold);letter-spacing:.06em; }
.slot-track { display:flex;align-items:center;gap:4px; }
.slot-dot { width:14px;height:14px;border-radius:50%;border:2px solid var(--gold-dim);background:rgba(201,168,76,.15);transition:all var(--transition); }
.slot-dot:hover { border-color:var(--gold); }
.slot-dot.used { background:transparent;border-color:var(--stone-border);opacity:.4; }
.slot-label { font-size:12px;color:var(--ash);margin-left:6px; }
.spell-row { display:flex;align-items:flex-start;gap:10px;padding:10px 12px;border-radius:var(--radius);margin-bottom:4px;background:var(--stone);border:1px solid var(--stone-border); }
.spell-row.unprepared { opacity:.5; }
.prepare-btn { width:22px;height:22px;border-radius:50%;border:2px solid var(--stone-border);background:transparent;display:flex;align-items:center;justify-content:center;color:transparent;transition:all var(--transition);flex-shrink:0;margin-top:2px; }
.prepare-btn:hover { border-color:var(--gold-dim); }
.prepare-btn.prepared { background:var(--gold-dim);border-color:var(--gold);color:var(--ink); }
.spell-main { flex:1; }
.spell-name { font-size:14px;font-weight:500;color:var(--parchment); }
.spell-meta { font-size:12px;color:var(--ash);margin-top:2px; }
.spell-desc { font-family:var(--font-body);font-size:14px;color:var(--ash-light);margin-top:6px;line-height:1.5; }
.notes-area { width:100%;min-height:400px;background:var(--stone);border:1px solid var(--stone-border);border-radius:var(--radius-lg);padding:20px;font-family:var(--font-body);font-size:16px;color:var(--parchment);line-height:1.7;resize:vertical; }
@media (max-width:900px) {
  .stats-grid { grid-template-columns:1fr; }
  .panel-full { grid-column:1; }
  .combat-bar { gap:20px; }
  .content { padding:16px; }
  .sheet-header { padding:16px; }
  .edit-form-grid { grid-template-columns:1fr; }
  .field.full { grid-column:1; }
  .ability-edit-grid { grid-template-columns:repeat(2,1fr); }
  .skill-edit-grid { grid-template-columns:1fr; }
  .currency-edit-grid { grid-template-columns:repeat(3,1fr); }
}
</style>
