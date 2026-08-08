import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import * as api from '../api';
import {
  IconChevronLeft, IconShield, IconBook, IconBackpack, IconStar,
  IconMoon, IconSun, IconPlus, IconTrash, IconCheck,
} from '../components/Icon';
import styles from './CharacterSheet.module.css';

type Tab = 'stats' | 'inventory' | 'spells' | 'notes';
const mod = (s: number) => Math.floor((s - 10) / 2);
const fmt = (n: number) => n >= 0 ? `+${n}` : `${n}`;
const pb  = (lvl: number) => Math.ceil(lvl / 4) + 1;

const CONDITIONS = ['Blinded','Charmed','Deafened','Exhaustion','Frightened','Grappled',
  'Incapacitated','Invisible','Paralyzed','Petrified','Poisoned','Prone','Restrained','Stunned','Unconscious'];

const ABILITIES = ['strength','dexterity','constitution','intelligence','wisdom','charisma'] as const;

const SKILLS = [
  { name: 'Acrobatics',     ab: 'dexterity',     f: 'skill_acrobatics' },
  { name: 'Animal Handling',ab: 'wisdom',         f: 'skill_animal_handling' },
  { name: 'Arcana',         ab: 'intelligence',   f: 'skill_arcana' },
  { name: 'Athletics',      ab: 'strength',       f: 'skill_athletics' },
  { name: 'Deception',      ab: 'charisma',       f: 'skill_deception' },
  { name: 'History',        ab: 'intelligence',   f: 'skill_history' },
  { name: 'Insight',        ab: 'wisdom',         f: 'skill_insight' },
  { name: 'Intimidation',   ab: 'charisma',       f: 'skill_intimidation' },
  { name: 'Investigation',  ab: 'intelligence',   f: 'skill_investigation' },
  { name: 'Medicine',       ab: 'wisdom',         f: 'skill_medicine' },
  { name: 'Nature',         ab: 'intelligence',   f: 'skill_nature' },
  { name: 'Perception',     ab: 'wisdom',         f: 'skill_perception' },
  { name: 'Performance',    ab: 'charisma',       f: 'skill_performance' },
  { name: 'Persuasion',     ab: 'charisma',       f: 'skill_persuasion' },
  { name: 'Religion',       ab: 'intelligence',   f: 'skill_religion' },
  { name: 'Sleight of Hand',ab: 'dexterity',      f: 'skill_sleight_of_hand' },
  { name: 'Stealth',        ab: 'dexterity',      f: 'skill_stealth' },
  { name: 'Survival',       ab: 'wisdom',         f: 'skill_survival' },
] as const;

export function CharacterSheet() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [char, setChar] = useState<any>(null);
  const [inventory, setInventory] = useState<any[]>([]);
  const [spells, setSpells] = useState<any[]>([]);
  const [slots, setSlots] = useState<any[]>([]);
  const [tab, setTab] = useState<Tab>('stats');
  const [loading, setLoading] = useState(true);
  const [hpInput, setHpInput] = useState('');
  const [hpMode, setHpMode] = useState<'damage' | 'heal' | 'temp'>('damage');

  const reload = useCallback(async () => {
    if (!id) return;
    const [c, inv, sp, sl] = await Promise.all([
      api.getCharacter(id), api.listInventory(id), api.listSpells(id), api.listSpellSlots(id),
    ]);
    setChar(c); setInventory(inv || []); setSpells(sp || []); setSlots(sl || []);
  }, [id]);

  useEffect(() => { reload().finally(() => setLoading(false)); }, [reload]);

  if (loading) return <div className={styles.loading}>Loading…</div>;
  if (!char)   return <div className={styles.loading}>Character not found.</div>;

  const profBonus = pb(char.level);
  const hpPct   = Math.max(0, Math.min(100, (char.current_hp / char.max_hp) * 100));
  const hpColor = hpPct > 50 ? '#27ae60' : hpPct > 25 ? '#f39c12' : '#c0392b';

  const handleHP = async () => {
    const n = parseInt(hpInput);
    if (!n || n <= 0 || !id) return;
    if (hpMode === 'damage') await api.applyDamage(id, n);
    else if (hpMode === 'heal') await api.applyHeal(id, n);
    else await api.addTempHP(id, n);
    setHpInput(''); reload();
  };

  const toggleCond = async (c: string) => {
    if (!id) return;
    const next = char.conditions.includes(c)
      ? char.conditions.filter((x: string) => x !== c)
      : [...char.conditions, c];
    await api.updateConditions(id, next); reload();
  };

  const handleRest = async (type: 'long' | 'short') => {
    if (!id) return;
    if (type === 'long') {
      await api.longRest(id);
    } else {
      const healed = Math.max(1, char.hit_dice_type + mod(char.constitution));
      await api.shortRest(id, Math.max(0, char.hit_dice_remaining - 1), Math.min(char.max_hp, char.current_hp + healed));
    }
    reload();
  };

  const TABS: [Tab, React.ReactNode, string][] = [
    ['stats',     <IconShield size={15} />,   'Stats'],
    ['inventory', <IconBackpack size={15} />, 'Inventory'],
    ['spells',    <IconBook size={15} />,     'Spells'],
    ['notes',     <IconStar size={15} />,     'Notes'],
  ];

  return (
    <div className={styles.page}>
      {/* Header */}
      <div className={styles.sheetHeader}>
        <button className={styles.back} onClick={() => navigate('/characters')}>
          <IconChevronLeft size={16} /> Characters
        </button>
        <div className={styles.headerInfo}>
          <h1 className={styles.charName}>{char.name}</h1>
          <p className={styles.charMeta}>Level {char.level} {char.race} {char.class} · {char.background} · {char.alignment}</p>
        </div>
        <div className={styles.headerBadges}>
          {char.inspiration && <span className={styles.badge} style={{ background: 'rgba(201,168,76,0.2)', color: 'var(--gold)' }}>✦ Inspired</span>}
          {char.is_npc      && <span className={styles.badge} style={{ background: 'rgba(139,26,26,0.2)', color: 'var(--crimson-light)' }}>NPC</span>}
          <span className={styles.badge}>XP {char.xp?.toLocaleString()}</span>
        </div>
      </div>

      {/* Combat bar */}
      <div className={styles.combatBar}>
        <div className={styles.hpSection}>
          <div className={styles.hpLabel}>HIT POINTS</div>
          <div className={styles.hpTrack}>
            <div className={styles.hpBar} style={{ width: `${hpPct}%`, background: hpColor }} />
          </div>
          <div className={styles.hpNumbers}>
            <span style={{ color: hpColor, fontWeight: 700, fontSize: 20 }}>{char.current_hp}</span>
            <span className={styles.hpMax}>/ {char.max_hp}</span>
            {char.temp_hp > 0 && <span className={styles.tempHp}>+{char.temp_hp} temp</span>}
          </div>
          <div className={styles.hpControls}>
            <div className={styles.hpModeRow}>
              {(['damage', 'heal', 'temp'] as const).map(m => (
                <button key={m} className={`${styles.hpModeBtn} ${hpMode === m ? styles.hpModeActive : ''}`} onClick={() => setHpMode(m)}>
                  {m === 'damage' ? '⚔ Dmg' : m === 'heal' ? '❤ Heal' : '🛡 Temp'}
                </button>
              ))}
            </div>
            <div className={styles.hpInputRow}>
              <input className={styles.hpInput} type="number" min={1} value={hpInput}
                onChange={e => setHpInput(e.target.value)} placeholder="0"
                onKeyDown={e => e.key === 'Enter' && handleHP()} />
              <button className={styles.hpApply} onClick={handleHP}>Apply</button>
            </div>
          </div>
        </div>

        <div className={styles.combatStats}>
          {[
            { l: 'AC',    v: char.armor_class },
            { l: 'INIT',  v: fmt(mod(char.dexterity)) },
            { l: 'SPEED', v: `${char.speed}ft` },
            { l: 'PROF',  v: `+${profBonus}` },
            { l: `d${char.hit_dice_type}`, v: `${char.hit_dice_remaining} left` },
          ].map(s => (
            <div key={s.l} className={styles.combatStat}>
              <div className={styles.combatStatVal}>{s.v}</div>
              <div className={styles.combatStatLabel}>{s.l}</div>
            </div>
          ))}
        </div>

        <div className={styles.restSection}>
          <button className={styles.restBtn} onClick={() => handleRest('short')}><IconMoon size={14} /> Short Rest</button>
          <button className={styles.restBtn} onClick={() => handleRest('long')}><IconSun size={14} /> Long Rest</button>
        </div>

        {char.current_hp <= 0 && (
          <div className={styles.deathSaves}>
            <div className={styles.deathLabel}>DEATH SAVES</div>
            {(['Successes', 'Failures'] as const).map(label => (
              <div key={label} className={styles.deathRow}>
                <span className={styles.deathName}>{label}</span>
                {[0, 1, 2].map(i => {
                  const count = label === 'Successes' ? char.death_save_successes : char.death_save_failures;
                  const cls = i < count ? (label === 'Successes' ? styles.deathSuccess : styles.deathFailure) : '';
                  return <button key={i} className={`${styles.deathDot} ${cls}`}
                    onClick={() => api.recordDeathSave(id!, label === 'Successes').then(reload)} />;
                })}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Tabs */}
      <div className={styles.tabs}>
        {TABS.map(([t, icon, label]) => (
          <button key={t} className={`${styles.tabBtn} ${tab === t ? styles.tabActive : ''}`} onClick={() => setTab(t)}>
            {icon} {label}
          </button>
        ))}
      </div>

      {/* Tab content */}
      <div className={styles.content}>
        {tab === 'stats'     && <StatsTab char={char} profBonus={profBonus} toggleCond={toggleCond} />}
        {tab === 'inventory' && <InventoryTab id={id!} items={inventory} char={char} reload={reload} />}
        {tab === 'spells'    && <SpellsTab id={id!} spells={spells} slots={slots} reload={reload} />}
        {tab === 'notes'     && <NotesTab char={char} id={id!} reload={reload} />}
      </div>
    </div>
  );
}

function StatsTab({ char, profBonus, toggleCond }: { char: any; profBonus: number; toggleCond: (c: string) => void }) {
  return (
    <div className={styles.statsGrid}>
      <div className={styles.panel}>
        <h3 className={styles.panelTitle}>Ability Scores</h3>
        <div className={styles.abilityGrid}>
          {ABILITIES.map(ab => {
            const score = char[ab];
            const m = mod(score);
            const saveProf = char[`save_prof_${ab}`] as boolean;
            return (
              <div key={ab} className={styles.abilityBox}>
                <div className={styles.abilityName}>{ab.slice(0, 3).toUpperCase()}</div>
                <div className={styles.abilityScore}>{score}</div>
                <div className={styles.abilityMod}>{fmt(m)}</div>
                <div className={`${styles.saveBadge} ${saveProf ? styles.saveProf : ''}`}>
                  {fmt(m + (saveProf ? profBonus : 0))} save
                </div>
              </div>
            );
          })}
        </div>
      </div>

      <div className={styles.panel}>
        <h3 className={styles.panelTitle}>Skills</h3>
        <div className={styles.skillList}>
          {SKILLS.map(sk => {
            const val = char[sk.f] as number;
            const total = mod(char[sk.ab]) + val * profBonus;
            return (
              <div key={sk.f} className={styles.skillRow}>
                <div className={`${styles.skillDot} ${val === 2 ? styles.expert : val === 1 ? styles.proficient : ''}`} />
                <span className={styles.skillName}>{sk.name}</span>
                <span className={styles.skillMod}>{fmt(total)}</span>
                <span className={styles.skillAbility}>{sk.ab.slice(0, 3)}</span>
              </div>
            );
          })}
        </div>
      </div>

      <div className={styles.panel}>
        <h3 className={styles.panelTitle}>Conditions</h3>
        <div className={styles.conditionGrid}>
          {CONDITIONS.map(c => (
            <button key={c}
              className={`${styles.conditionPill} ${char.conditions?.includes(c) ? styles.condActive : ''}`}
              onClick={() => toggleCond(c)}>
              {c}
            </button>
          ))}
        </div>
      </div>

      <div className={styles.panel}>
        <h3 className={styles.panelTitle}>Currency</h3>
        <div className={styles.currencyRow}>
          {([['CP', char.copper, '#b5651d'], ['SP', char.silver, '#c0c0c0'], ['EP', char.electrum, '#b8b8ff'], ['GP', char.gold, '#ffd700'], ['PP', char.platinum, '#e5e4e2']] as const).map(([l, v, color]) => (
            <div key={l} className={styles.coin}>
              <div className={styles.coinCircle} style={{ borderColor: color, color }}>{l}</div>
              <div className={styles.coinVal}>{v}</div>
            </div>
          ))}
        </div>
      </div>

      <div className={`${styles.panel} ${styles.panelFull}`}>
        <h3 className={styles.panelTitle}>Character Traits</h3>
        <div className={styles.traitsGrid}>
          {[['Personality', char.personality_traits], ['Ideals', char.ideals], ['Bonds', char.bonds], ['Flaws', char.flaws]].map(([l, v]) => (
            <div key={l as string} className={styles.traitBox}>
              <div className={styles.traitLabel}>{l as string}</div>
              <div className={styles.traitText}>{(v as string) || '—'}</div>
            </div>
          ))}
        </div>
      </div>

      <div className={`${styles.panel} ${styles.panelFull}`}>
        <h3 className={styles.panelTitle}>Proficiencies & Training</h3>
        <div className={styles.traitsGrid}>
          {[['Armor', char.training_armor], ['Weapons', char.training_weapons], ['Tools', char.training_tools], ['Languages', char.training_languages]].map(([l, v]) => (
            <div key={l as string} className={styles.traitBox}>
              <div className={styles.traitLabel}>{l as string}</div>
              <div className={styles.traitText}>{((v as string[]) || []).join(', ') || '—'}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function InventoryTab({ id, items, char, reload }: { id: string; items: any[]; char: any; reload: () => void }) {
  const [showAdd, setShowAdd] = useState(false);
  const [form, setForm] = useState({ name: '', quantity: 1, weight: 0, value: 0, description: '', is_equipped: false, requires_attunement: false });

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    await api.createInventoryItem(id, form);
    setShowAdd(false);
    setForm({ name: '', quantity: 1, weight: 0, value: 0, description: '', is_equipped: false, requires_attunement: false });
    reload();
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.tabHeader}>
        <div className={styles.tabMeta}>
          <span>{items.reduce((s, i) => s + i.weight * i.quantity, 0)} lbs</span>
          <span>{items.filter(i => i.is_attuned).length}/{char.attunement_slots} attuned</span>
        </div>
        <button className={styles.addBtn} onClick={() => setShowAdd(true)}><IconPlus size={14} /> Add Item</button>
      </div>

      {showAdd && (
        <form onSubmit={submit} className={styles.addForm}>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Name *</label><input value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} required /></div>
            <div className={styles.field}><label>Qty</label><input type="number" min={1} value={form.quantity} onChange={e => setForm(f => ({ ...f, quantity: +e.target.value }))} /></div>
            <div className={styles.field}><label>Weight (lb)</label><input type="number" min={0} value={form.weight} onChange={e => setForm(f => ({ ...f, weight: +e.target.value }))} /></div>
            <div className={styles.field}><label>Value (gp)</label><input type="number" min={0} value={form.value} onChange={e => setForm(f => ({ ...f, value: +e.target.value }))} /></div>
          </div>
          <div className={styles.field}><label>Description</label><input value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} /></div>
          <label className={styles.check}>
            <input type="checkbox" checked={form.requires_attunement} onChange={e => setForm(f => ({ ...f, requires_attunement: e.target.checked }))} />
            Requires Attunement
          </label>
          <div className={styles.modalActions}>
            <button type="button" className={styles.cancelBtn} onClick={() => setShowAdd(false)}>Cancel</button>
            <button type="submit" className={styles.submitBtn}>Add</button>
          </div>
        </form>
      )}

      <div className={styles.itemList}>
        {items.map(item => (
          <div key={item.id} className={`${styles.itemRow} ${item.is_equipped ? styles.itemEquipped : ''}`}>
            <div className={styles.itemMain}>
              <div className={styles.itemName}>
                {item.name}
                {item.is_equipped && <span className={styles.equippedTag}>Equipped</span>}
                {item.is_attuned  && <span className={styles.attunedTag}>Attuned</span>}
              </div>
              {item.description && <div className={styles.itemDesc}>{item.description}</div>}
              <div className={styles.itemMeta}>
                <span>Qty: {item.quantity}</span>
                {item.weight > 0 && <span>{item.weight}lb</span>}
                {item.value > 0  && <span>{item.value}gp</span>}
              </div>
            </div>
            <div className={styles.itemActions}>
              <button className={styles.itemBtn}
                onClick={() => api.updateInventoryItem(id, item.id, { ...item, is_equipped: !item.is_equipped }).then(reload)}>
                {item.is_equipped ? '⚔' : '○'}
              </button>
              {item.requires_attunement && (
                <button className={styles.itemBtn}
                  onClick={() => (item.is_attuned ? api.unattuneItem : api.attuneItem)(id, item.id).then(reload)}>
                  ✦
                </button>
              )}
              <button className={styles.itemBtnDanger} onClick={() => api.deleteInventoryItem(id, item.id).then(reload)}>
                <IconTrash size={13} />
              </button>
            </div>
          </div>
        ))}
        {items.length === 0 && <div className={styles.emptyState}>No items yet. Add something!</div>}
      </div>
    </div>
  );
}

const SCHOOLS = ['Abjuration','Conjuration','Divination','Enchantment','Evocation','Illusion','Necromancy','Transmutation'];

function SpellsTab({ id, spells, slots, reload }: { id: string; spells: any[]; slots: any[]; reload: () => void }) {
  const [showAdd, setShowAdd] = useState(false);
  const [form, setForm] = useState({
    name: '', level: 1, school: 'Evocation', casting_time: '1 action',
    range: '60 feet', components: 'V, S', duration: 'Instantaneous', description: '', is_prepared: true,
  });

  const byLevel = Array.from({ length: 10 }, (_, i) => ({
    level: i,
    spells: spells.filter(s => s.level === i),
    slot: i > 0 ? slots.find(s => s.spell_level === i) : undefined,
  })).filter(g => g.spells.length > 0 || (g.slot && g.slot.total > 0));

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    await api.createSpell(id, form);
    setShowAdd(false); reload();
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.tabHeader}>
        <div className={styles.tabMeta}><span>{spells.filter(s => s.is_prepared).length} prepared</span></div>
        <button className={styles.addBtn} onClick={() => setShowAdd(true)}><IconPlus size={14} /> Add Spell</button>
      </div>

      {showAdd && (
        <form onSubmit={submit} className={styles.addForm}>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Name *</label><input value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} required /></div>
            <div className={styles.field}><label>Level</label>
              <select value={form.level} onChange={e => setForm(f => ({ ...f, level: +e.target.value }))}>
                {[0,1,2,3,4,5,6,7,8,9].map(l => <option key={l} value={l}>{l === 0 ? 'Cantrip' : `Level ${l}`}</option>)}
              </select>
            </div>
            <div className={styles.field}><label>School</label>
              <select value={form.school} onChange={e => setForm(f => ({ ...f, school: e.target.value }))}>
                {SCHOOLS.map(s => <option key={s}>{s}</option>)}
              </select>
            </div>
          </div>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Cast Time</label><input value={form.casting_time} onChange={e => setForm(f => ({ ...f, casting_time: e.target.value }))} /></div>
            <div className={styles.field}><label>Range</label><input value={form.range} onChange={e => setForm(f => ({ ...f, range: e.target.value }))} /></div>
            <div className={styles.field}><label>Duration</label><input value={form.duration} onChange={e => setForm(f => ({ ...f, duration: e.target.value }))} /></div>
          </div>
          <div className={styles.field}><label>Components</label><input value={form.components} onChange={e => setForm(f => ({ ...f, components: e.target.value }))} /></div>
          <div className={styles.field}><label>Description</label><textarea value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} rows={3} /></div>
          <div className={styles.modalActions}>
            <button type="button" className={styles.cancelBtn} onClick={() => setShowAdd(false)}>Cancel</button>
            <button type="submit" className={styles.submitBtn}>Add Spell</button>
          </div>
        </form>
      )}

      {byLevel.map(group => (
        <div key={group.level} className={styles.spellGroup}>
          <div className={styles.spellGroupHeader}>
            <span>{group.level === 0 ? 'Cantrips' : `Level ${group.level}`}</span>
            {group.slot && (
              <div className={styles.slotTrack}>
                {Array.from({ length: group.slot.total }, (_, i) => (
                  <button key={i}
                    className={`${styles.slotDot} ${i < group.slot!.used ? styles.slotUsed : ''}`}
                    onClick={() => api.useSpellSlot(id, group.level).then(reload)} />
                ))}
                <span className={styles.slotLabel}>{group.slot.total - group.slot.used}/{group.slot.total}</span>
              </div>
            )}
          </div>
          {group.spells.map(spell => (
            <div key={spell.id} className={`${styles.spellRow} ${!spell.is_prepared && spell.level > 0 ? styles.unprepared : ''}`}>
              <button className={`${styles.prepareBtn} ${spell.is_prepared ? styles.prepared : ''}`}
                onClick={() => api.toggleSpellPrepared(id, spell.id).then(reload)}>
                <IconCheck size={12} />
              </button>
              <div className={styles.spellMain}>
                <div className={styles.spellName}>{spell.name}</div>
                <div className={styles.spellMeta}>{spell.school} · {spell.casting_time} · {spell.range}</div>
                {spell.description && <div className={styles.spellDesc}>{spell.description}</div>}
              </div>
              <button className={styles.itemBtnDanger} onClick={() => api.deleteSpell(id, spell.id).then(reload)}>
                <IconTrash size={13} />
              </button>
            </div>
          ))}
        </div>
      ))}
      {byLevel.length === 0 && <div className={styles.emptyState}>No spells yet. Add some!</div>}
    </div>
  );
}

function NotesTab({ char, id, reload }: { char: any; id: string; reload: () => void }) {
  const [notes, setNotes] = useState(char.notes || '');
  const [saved, setSaved] = useState(false);

  const save = async () => {
    await api.updateCharacterInfo(id, { ...char, notes });
    setSaved(true); setTimeout(() => setSaved(false), 2000); reload();
  };

  return (
    <div className={styles.tabContent}>
      <div className={styles.tabHeader}>
        <span />
        <button className={styles.addBtn} onClick={save}>{saved ? '✓ Saved' : 'Save Notes'}</button>
      </div>
      <textarea className={styles.notesArea} value={notes} onChange={e => setNotes(e.target.value)}
        placeholder="Session notes, backstory, important NPCs…" />
    </div>
  );
}
