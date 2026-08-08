import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import * as api from '../lib/api';
import {
    IconChevronLeft, IconShield, IconBook, IconBackpack, IconStar,
    IconMoon, IconSun, IconPlus, IconTrash, IconCheck,
} from '../components/Icon';

type Tab = 'stats' | 'inventory' | 'spells' | 'notes';
const mod = (s: number) => Math.floor((s - 10) / 2);
const fmt = (n: number) => n >= 0 ? `+${n}` : `${n}`;
const pb = (lvl: number) => Math.ceil(lvl / 4) + 1;

const CONDITIONS = ['Blinded', 'Charmed', 'Deafened', 'Exhaustion', 'Frightened', 'Grappled',
    'Incapacitated', 'Invisible', 'Paralyzed', 'Petrified', 'Poisoned', 'Prone', 'Restrained', 'Stunned', 'Unconscious'];

const ABILITIES = ['strength', 'dexterity', 'constitution', 'intelligence', 'wisdom', 'charisma'] as const;

const SKILLS = [
    { name: 'Acrobatics', ab: 'dexterity', f: 'skill_acrobatics' },
    { name: 'Animal Handling', ab: 'wisdom', f: 'skill_animal_handling' },
    { name: 'Arcana', ab: 'intelligence', f: 'skill_arcana' },
    { name: 'Athletics', ab: 'strength', f: 'skill_athletics' },
    { name: 'Deception', ab: 'charisma', f: 'skill_deception' },
    { name: 'History', ab: 'intelligence', f: 'skill_history' },
    { name: 'Insight', ab: 'wisdom', f: 'skill_insight' },
    { name: 'Intimidation', ab: 'charisma', f: 'skill_intimidation' },
    { name: 'Investigation', ab: 'intelligence', f: 'skill_investigation' },
    { name: 'Medicine', ab: 'wisdom', f: 'skill_medicine' },
    { name: 'Nature', ab: 'intelligence', f: 'skill_nature' },
    { name: 'Perception', ab: 'wisdom', f: 'skill_perception' },
    { name: 'Performance', ab: 'charisma', f: 'skill_performance' },
    { name: 'Persuasion', ab: 'charisma', f: 'skill_persuasion' },
    { name: 'Religion', ab: 'intelligence', f: 'skill_religion' },
    { name: 'Sleight of Hand', ab: 'dexterity', f: 'skill_sleight_of_hand' },
    { name: 'Stealth', ab: 'dexterity', f: 'skill_stealth' },
    { name: 'Survival', ab: 'wisdom', f: 'skill_survival' },
] as const;

// Shared utility class fragments reused across tabs
const panel = "bg-stone border border-stone-border rounded-lg p-5";
const panelTitle = "font-display text-[13px] font-semibold text-gold tracking-wide uppercase mb-4";
const addForm = "bg-stone border border-stone-border rounded-lg p-5 mb-5 flex flex-col gap-3";
const formRow = "flex gap-3 flex-wrap [&>*]:flex-1 [&>*]:min-w-[120px]";
const field = "flex flex-col gap-[5px]";
const fieldLabel = "text-[11px] text-ash uppercase tracking-wide";
const modalActions = "flex justify-end gap-2";
const cancelBtn = "bg-stone-mid border-none text-ash-light px-4 py-2 rounded-md text-[13px] hover:text-parchment";
const submitBtn = "bg-crimson border-none text-parchment px-4 py-2 rounded-md font-display text-[13px] font-semibold hover:bg-crimson-light";
const tabHeader = "flex items-center justify-between mb-4";
const addBtn = "flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment";
const emptyState = "text-center py-10 text-ash text-sm";
const itemBtnDanger = "bg-transparent border-none text-stone-border p-[5px] transition-colors duration-180 hover:text-crimson-light";

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

    if (loading) return <div className="flex items-center justify-center h-[300px] text-ash">Loading…</div>;
    if (!char) return <div className="flex items-center justify-center h-[300px] text-ash">Character not found.</div>;

    const profBonus = pb(char.level);
    const hpPct = Math.max(0, Math.min(100, (char.current_hp / char.max_hp) * 100));
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
        ['stats', <IconShield size={15} />, 'Stats'],
        ['inventory', <IconBackpack size={15} />, 'Inventory'],
        ['spells', <IconBook size={15} />, 'Spells'],
        ['notes', <IconStar size={15} />, 'Notes'],
    ];

    return (
        <div className="min-h-screen">
            {/* Header */}
            <div className="flex items-center gap-5 py-5 px-8 bg-stone border-b border-stone-border flex-wrap">
                <button className="flex items-center gap-1 bg-transparent border-none text-ash text-[13px] transition-colors duration-180 hover:text-parchment" onClick={() => navigate('/characters')}>
                    <IconChevronLeft size={16} /> Characters
                </button>
                <div className="flex-1">
                    <h1 className="font-display text-2xl font-bold text-parchment">{char.name}</h1>
                    <p className="text-[13px] text-ash mt-[3px]">Level {char.level} {char.race} {char.class} · {char.background} · {char.alignment}</p>
                </div>
                <div className="flex gap-2 flex-wrap">
                    {char.inspiration && <span className="text-[11px] px-2 py-[3px] rounded-full bg-gold/20 text-gold">✦ Inspired</span>}
                    {char.is_npc && <span className="text-[11px] px-2 py-[3px] rounded-full bg-crimson/20 text-crimson-light">NPC</span>}
                    <span className="text-[11px] px-2 py-[3px] rounded-full bg-stone-mid text-ash-light">XP {char.xp?.toLocaleString()}</span>
                </div>
            </div>

            {/* Combat bar */}
            <div className="py-5 px-8 bg-stone-mid border-b border-stone-border flex gap-8 items-start flex-wrap">
                <div className="min-w-[220px]">
                    <div className="text-[10px] font-bold tracking-widest text-ash mb-1.5">HIT POINTS</div>
                    <div className="h-1.5 bg-stone-border rounded-sm mb-2 overflow-hidden">
                        <div className="h-full rounded-sm transition-[width,background] duration-[0.4s] ease-out" style={{ width: `${hpPct}%`, background: hpColor }} />
                    </div>
                    <div className="flex items-baseline gap-1.5 mb-3">
                        <span style={{ color: hpColor, fontWeight: 700, fontSize: 20 }}>{char.current_hp}</span>
                        <span className="text-ash text-sm">/ {char.max_hp}</span>
                        {char.temp_hp > 0 && <span className="text-[13px] text-info">+{char.temp_hp} temp</span>}
                    </div>
                    <div className="flex flex-col gap-1.5">
                        <div className="flex gap-1">
                            {(['damage', 'heal', 'temp'] as const).map(m => (
                                <button key={m}
                                    className={`bg-stone border border-stone-border text-ash px-2.5 py-[5px] rounded-sm text-[11px] transition-all duration-180 ${hpMode === m ? 'bg-stone-light border-gold-dim text-parchment' : ''}`}
                                    onClick={() => setHpMode(m)}>
                                    {m === 'damage' ? '⚔ Dmg' : m === 'heal' ? '❤ Heal' : '🛡 Temp'}
                                </button>
                            ))}
                        </div>
                        <div className="flex gap-1.5">
                            <input className="w-[70px] text-center px-2 py-1.5" type="number" min={1} value={hpInput}
                                onChange={e => setHpInput(e.target.value)} placeholder="0"
                                onKeyDown={e => e.key === 'Enter' && handleHP()} />
                            <button className="bg-crimson border-none text-parchment px-3.5 py-1.5 rounded-sm text-[13px] transition-colors duration-180 hover:bg-crimson-light" onClick={handleHP}>Apply</button>
                        </div>
                    </div>
                </div>

                <div className="flex gap-4 flex-wrap items-center">
                    {[
                        { l: 'AC', v: char.armor_class },
                        { l: 'INIT', v: fmt(mod(char.dexterity)) },
                        { l: 'SPEED', v: `${char.speed}ft` },
                        { l: 'PROF', v: `+${profBonus}` },
                        { l: `d${char.hit_dice_type}`, v: `${char.hit_dice_remaining} left` },
                    ].map(s => (
                        <div key={s.l} className="text-center">
                            <div className="font-display text-[22px] font-bold text-parchment">{s.v}</div>
                            <div className="text-[10px] text-ash tracking-wide mt-0.5">{s.l}</div>
                        </div>
                    ))}
                </div>

                <div className="flex flex-col gap-1.5">
                    <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-xs transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => handleRest('short')}><IconMoon size={14} /> Short Rest</button>
                    <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-xs transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => handleRest('long')}><IconSun size={14} /> Long Rest</button>
                </div>

                {char.current_hp <= 0 && (
                    <div>
                        <div className="text-[10px] font-bold tracking-wide text-crimson-light mb-2">DEATH SAVES</div>
                        {(['Successes', 'Failures'] as const).map(label => (
                            <div key={label} className="flex items-center gap-2 mb-1.5">
                                <span className="text-[11px] text-ash w-[70px]">{label}</span>
                                {[0, 1, 2].map(i => {
                                    const count = label === 'Successes' ? char.death_save_successes : char.death_save_failures;
                                    const cls = i < count ? (label === 'Successes' ? 'bg-emerald border-emerald' : 'bg-crimson border-crimson') : '';
                                    return <button key={i} className={`w-[18px] h-[18px] rounded-full border-2 border-stone-border bg-transparent transition-all duration-180 ${cls}`}
                                        onClick={() => api.recordDeathSave(id!, label === 'Successes').then(reload)} />;
                                })}
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Tabs */}
            <div className="flex gap-0.5 px-8 bg-stone border-b border-stone-border">
                {TABS.map(([t, icon, label]) => (
                    <button key={t} className={`flex items-center gap-1.5 bg-transparent border-none text-ash px-[18px] py-3.5 text-[13px] font-medium border-b-2 border-transparent transition-all duration-180 hover:text-parchment ${tab === t ? 'text-gold border-gold' : ''}`} onClick={() => setTab(t)}>
                        {icon} {label}
                    </button>
                ))}
            </div>

            {/* Tab content */}
            <div className="py-6 px-8">
                {tab === 'stats' && <StatsTab char={char} profBonus={profBonus} toggleCond={toggleCond} />}
                {tab === 'inventory' && <InventoryTab id={id!} items={inventory} char={char} reload={reload} />}
                {tab === 'spells' && <SpellsTab id={id!} spells={spells} slots={slots} reload={reload} />}
                {tab === 'notes' && <NotesTab char={char} id={id!} reload={reload} />}
            </div>
        </div>
    );
}

function StatsTab({ char, profBonus, toggleCond }: { char: any; profBonus: number; toggleCond: (c: string) => void }) {
    return (
        <div className="grid grid-cols-[auto_1fr_1fr] gap-5 max-[900px]:grid-cols-1">
            <div className={panel}>
                <h3 className={panelTitle}>Ability Scores</h3>
                <div className="grid grid-cols-3 gap-2.5">
                    {ABILITIES.map(ab => {
                        const score = char[ab];
                        const m = mod(score);
                        const saveProf = char[`save_prof_${ab}`] as boolean;
                        return (
                            <div key={ab} className="flex flex-col items-center bg-stone-mid border border-stone-border rounded-md py-2.5 px-2">
                                <div className="text-[9px] font-bold tracking-widest text-ash mb-1">{ab.slice(0, 3).toUpperCase()}</div>
                                <div className="font-display text-2xl font-bold text-parchment leading-none">{score}</div>
                                <div className="text-base font-semibold text-gold">{fmt(m)}</div>
                                <div className={`text-[10px] text-ash mt-1 ${saveProf ? 'text-gold-dim' : ''}`}>
                                    {fmt(m + (saveProf ? profBonus : 0))} save
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            <div className={panel}>
                <h3 className={panelTitle}>Skills</h3>
                <div className="flex flex-col gap-[3px]">
                    {SKILLS.map(sk => {
                        const val = char[sk.f] as number;
                        const total = mod(char[sk.ab]) + val * profBonus;
                        return (
                            <div key={sk.f} className="flex items-center gap-2 py-1 px-1.5 rounded-sm hover:bg-stone-mid">
                                <div className={`w-2.5 h-2.5 rounded-full border-2 border-stone-border shrink-0 ${val === 2 ? 'bg-gold border-gold' : val === 1 ? 'bg-gold-dim border-gold-dim' : ''}`} />
                                <span className="flex-1 text-[13px] text-ash-light">{sk.name}</span>
                                <span className="font-semibold text-[13px] text-parchment w-[30px] text-right">{fmt(total)}</span>
                                <span className="text-[10px] text-ash w-6 text-center">{sk.ab.slice(0, 3)}</span>
                            </div>
                        );
                    })}
                </div>
            </div>

            <div className={panel}>
                <h3 className={panelTitle}>Conditions</h3>
                <div className="flex flex-wrap gap-1.5">
                    {CONDITIONS.map(c => (
                        <button key={c}
                            className={`bg-stone-mid border border-stone-border text-ash px-2.5 py-[5px] rounded-full text-xs transition-all duration-180 hover:border-crimson hover:text-parchment ${char.conditions?.includes(c) ? 'bg-crimson/20 border-crimson text-crimson-light' : ''}`}
                            onClick={() => toggleCond(c)}>
                            {c}
                        </button>
                    ))}
                </div>
            </div>

            <div className={panel}>
                <h3 className={panelTitle}>Currency</h3>
                <div className="flex gap-3 justify-center flex-wrap">
                    {([['CP', char.copper, '#b5651d'], ['SP', char.silver, '#c0c0c0'], ['EP', char.electrum, '#b8b8ff'], ['GP', char.gold, '#ffd700'], ['PP', char.platinum, '#e5e4e2']] as const).map(([l, v, color]) => (
                        <div key={l} className="flex flex-col items-center gap-1.5">
                            <div className="w-10 h-10 rounded-full border-2 flex items-center justify-center text-[10px] font-bold tracking-wide" style={{ borderColor: color, color }}>{l}</div>
                            <div className="text-[15px] font-semibold text-parchment">{v}</div>
                        </div>
                    ))}
                </div>
            </div>

            <div className={`${panel} col-span-full`}>
                <h3 className={panelTitle}>Character Traits</h3>
                <div className="grid grid-cols-2 gap-4">
                    {[['Personality', char.personality_traits], ['Ideals', char.ideals], ['Bonds', char.bonds], ['Flaws', char.flaws]].map(([l, v]) => (
                        <div key={l as string} className="bg-stone-mid rounded-md p-3.5">
                            <div className="text-[10px] font-bold text-gold-dim tracking-wide uppercase mb-1.5">{l as string}</div>
                            <div className="font-body text-[15px] text-parchment leading-relaxed">{(v as string) || '—'}</div>
                        </div>
                    ))}
                </div>
            </div>

            <div className={`${panel} col-span-full`}>
                <h3 className={panelTitle}>Proficiencies & Training</h3>
                <div className="grid grid-cols-2 gap-4">
                    {[['Armor', char.training_armor], ['Weapons', char.training_weapons], ['Tools', char.training_tools], ['Languages', char.training_languages]].map(([l, v]) => (
                        <div key={l as string} className="bg-stone-mid rounded-md p-3.5">
                            <div className="text-[10px] font-bold text-gold-dim tracking-wide uppercase mb-1.5">{l as string}</div>
                            <div className="font-body text-[15px] text-parchment leading-relaxed">{((v as string[]) || []).join(', ') || '—'}</div>
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
        <div>
            <div className={tabHeader}>
                <div className="flex gap-4 text-[13px] text-ash">
                    <span>{items.reduce((s, i) => s + i.weight * i.quantity, 0)} lbs</span>
                    <span>{items.filter(i => i.is_attuned).length}/{char.attunement_slots} attuned</span>
                </div>
                <button className={addBtn} onClick={() => setShowAdd(true)}><IconPlus size={14} /> Add Item</button>
            </div>

            {showAdd && (
                <form onSubmit={submit} className={addForm}>
                    <div className={formRow}>
                        <div className={field}><label className={fieldLabel}>Name *</label><input value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} required /></div>
                        <div className={field}><label className={fieldLabel}>Qty</label><input type="number" min={1} value={form.quantity} onChange={e => setForm(f => ({ ...f, quantity: +e.target.value }))} /></div>
                        <div className={field}><label className={fieldLabel}>Weight (lb)</label><input type="number" min={0} value={form.weight} onChange={e => setForm(f => ({ ...f, weight: +e.target.value }))} /></div>
                        <div className={field}><label className={fieldLabel}>Value (gp)</label><input type="number" min={0} value={form.value} onChange={e => setForm(f => ({ ...f, value: +e.target.value }))} /></div>
                    </div>
                    <div className={field}><label className={fieldLabel}>Description</label><input value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} /></div>
                    <label className="flex items-center gap-1.5 text-[13px] text-ash-light cursor-pointer">
                        <input className="w-auto" type="checkbox" checked={form.requires_attunement} onChange={e => setForm(f => ({ ...f, requires_attunement: e.target.checked }))} />
                        Requires Attunement
                    </label>
                    <div className={modalActions}>
                        <button type="button" className={cancelBtn} onClick={() => setShowAdd(false)}>Cancel</button>
                        <button type="submit" className={submitBtn}>Add</button>
                    </div>
                </form>
            )}

            <div className="flex flex-col gap-1.5">
                {items.map(item => (
                    <div key={item.id} className={`flex items-center gap-3 bg-stone border rounded-md py-3 px-4 ${item.is_equipped ? 'border-gold-dim' : 'border-stone-border'}`}>
                        <div className="flex-1">
                            <div className="text-sm font-medium text-parchment flex items-center gap-2">
                                {item.name}
                                {item.is_equipped && <span className="text-[10px] bg-gold/15 text-gold border border-gold/30 px-1.5 py-px rounded-sm">Equipped</span>}
                                {item.is_attuned && <span className="text-[10px] bg-info/15 text-info border border-info/30 px-1.5 py-px rounded-sm">Attuned</span>}
                            </div>
                            {item.description && <div className="text-xs text-ash mt-[3px]">{item.description}</div>}
                            <div className="flex gap-2.5 text-[11px] text-ash mt-1">
                                <span>Qty: {item.quantity}</span>
                                {item.weight > 0 && <span>{item.weight}lb</span>}
                                {item.value > 0 && <span>{item.value}gp</span>}
                            </div>
                        </div>
                        <div className="flex gap-1.5">
                            <button className="bg-transparent border border-stone-border text-ash px-2.5 py-[5px] rounded-sm text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-gold"
                                onClick={() => api.updateInventoryItem(id, item.id, { ...item, is_equipped: !item.is_equipped }).then(reload)}>
                                {item.is_equipped ? '⚔' : '○'}
                            </button>
                            {item.requires_attunement && (
                                <button className="bg-transparent border border-stone-border text-ash px-2.5 py-[5px] rounded-sm text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-gold"
                                    onClick={() => (item.is_attuned ? api.unattuneItem : api.attuneItem)(id, item.id).then(reload)}>
                                    ✦
                                </button>
                            )}
                            <button className={itemBtnDanger} onClick={() => api.deleteInventoryItem(id, item.id).then(reload)}>
                                <IconTrash size={13} />
                            </button>
                        </div>
                    </div>
                ))}
                {items.length === 0 && <div className={emptyState}>No items yet. Add something!</div>}
            </div>
        </div>
    );
}

const SCHOOLS = ['Abjuration', 'Conjuration', 'Divination', 'Enchantment', 'Evocation', 'Illusion', 'Necromancy', 'Transmutation'];

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
        <div>
            <div className={tabHeader}>
                <div className="flex gap-4 text-[13px] text-ash"><span>{spells.filter(s => s.is_prepared).length} prepared</span></div>
                <button className={addBtn} onClick={() => setShowAdd(true)}><IconPlus size={14} /> Add Spell</button>
            </div>

            {showAdd && (
                <form onSubmit={submit} className={addForm}>
                    <div className={formRow}>
                        <div className={field}><label className={fieldLabel}>Name *</label><input value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} required /></div>
                        <div className={field}><label className={fieldLabel}>Level</label>
                            <select value={form.level} onChange={e => setForm(f => ({ ...f, level: +e.target.value }))}>
                                {[0, 1, 2, 3, 4, 5, 6, 7, 8, 9].map(l => <option key={l} value={l}>{l === 0 ? 'Cantrip' : `Level ${l}`}</option>)}
                            </select>
                        </div>
                        <div className={field}><label className={fieldLabel}>School</label>
                            <select value={form.school} onChange={e => setForm(f => ({ ...f, school: e.target.value }))}>
                                {SCHOOLS.map(s => <option key={s}>{s}</option>)}
                            </select>
                        </div>
                    </div>
                    <div className={formRow}>
                        <div className={field}><label className={fieldLabel}>Cast Time</label><input value={form.casting_time} onChange={e => setForm(f => ({ ...f, casting_time: e.target.value }))} /></div>
                        <div className={field}><label className={fieldLabel}>Range</label><input value={form.range} onChange={e => setForm(f => ({ ...f, range: e.target.value }))} /></div>
                        <div className={field}><label className={fieldLabel}>Duration</label><input value={form.duration} onChange={e => setForm(f => ({ ...f, duration: e.target.value }))} /></div>
                    </div>
                    <div className={field}><label className={fieldLabel}>Components</label><input value={form.components} onChange={e => setForm(f => ({ ...f, components: e.target.value }))} /></div>
                    <div className={field}><label className={fieldLabel}>Description</label><textarea value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} rows={3} /></div>
                    <div className={modalActions}>
                        <button type="button" className={cancelBtn} onClick={() => setShowAdd(false)}>Cancel</button>
                        <button type="submit" className={submitBtn}>Add Spell</button>
                    </div>
                </form>
            )}

            {byLevel.map(group => (
                <div key={group.level} className="mb-5">
                    <div className="flex items-center justify-between py-2 border-b border-stone-border mb-2 font-display text-[13px] font-semibold text-gold tracking-wide">
                        <span>{group.level === 0 ? 'Cantrips' : `Level ${group.level}`}</span>
                        {group.slot && (
                            <div className="flex items-center gap-1">
                                {Array.from({ length: group.slot.total }, (_, i) => (
                                    <button key={i}
                                        className={`w-3.5 h-3.5 rounded-full border-2 border-gold-dim bg-gold/15 transition-all duration-180 hover:border-gold ${i < group.slot!.used ? 'bg-transparent border-stone-border opacity-40' : ''}`}
                                        onClick={() => api.useSpellSlot(id, group.level).then(reload)} />
                                ))}
                                <span className="text-xs text-ash ml-1.5">{group.slot.total - group.slot.used}/{group.slot.total}</span>
                            </div>
                        )}
                    </div>
                    {group.spells.map(spell => (
                        <div key={spell.id} className={`flex items-start gap-2.5 py-2.5 px-3 rounded-md mb-1 bg-stone border border-stone-border ${!spell.is_prepared && spell.level > 0 ? 'opacity-50' : ''}`}>
                            <button className={`w-[22px] h-[22px] rounded-full border-2 border-stone-border bg-transparent flex items-center justify-center transition-all duration-180 shrink-0 mt-0.5 hover:border-gold-dim ${spell.is_prepared ? 'bg-gold-dim border-gold text-ink' : 'text-transparent'}`}
                                onClick={() => api.toggleSpellPrepared(id, spell.id).then(reload)}>
                                <IconCheck size={12} />
                            </button>
                            <div className="flex-1">
                                <div className="text-sm font-medium text-parchment">{spell.name}</div>
                                <div className="text-xs text-ash mt-0.5">{spell.school} · {spell.casting_time} · {spell.range}</div>
                                {spell.description && <div className="font-body text-sm text-ash-light mt-1.5 leading-relaxed">{spell.description}</div>}
                            </div>
                            <button className={itemBtnDanger} onClick={() => api.deleteSpell(id, spell.id).then(reload)}>
                                <IconTrash size={13} />
                            </button>
                        </div>
                    ))}
                </div>
            ))}
            {byLevel.length === 0 && <div className={emptyState}>No spells yet. Add some!</div>}
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
        <div>
            <div className={tabHeader}>
                <span />
                <button className={addBtn} onClick={save}>{saved ? '✓ Saved' : 'Save Notes'}</button>
            </div>
            <textarea
                className="w-full min-h-[400px] bg-stone border border-stone-border rounded-lg p-5 font-body text-base text-parchment leading-[1.7] resize-y"
                value={notes} onChange={e => setNotes(e.target.value)}
                placeholder="Session notes, backstory, important NPCs…" />
        </div>
    );
}
