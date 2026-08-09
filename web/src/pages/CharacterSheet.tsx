import { useState, useEffect, useCallback, type ReactNode } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import * as api from '../lib/api';
import {
    IconChevronLeft, IconShield, IconBook, IconBackpack, IconStar,
    IconMoon, IconSun, IconPlus, IconTrash, IconCheck, IconZap,
} from '../components/Icon';

type Tab = 'stats' | 'inventory' | 'spells' | 'background';
const mod = (s: number) => Math.floor((s - 10) / 2);
const fmt = (n: number) => n >= 0 ? `+${n}` : `${n}`;
const pb = (lvl: number) => Math.ceil(lvl / 4) + 1;

const CONDITIONS = ['Blinded', 'Charmed', 'Deafened', 'Exhaustion', 'Frightened', 'Grappled',
    'Incapacitated', 'Invisible', 'Paralyzed', 'Petrified', 'Poisoned', 'Prone', 'Restrained', 'Stunned', 'Unconscious'];

const ABILITIES = ['strength', 'dexterity', 'constitution', 'intelligence', 'wisdom', 'charisma'] as const;

const ACTION_TYPES = [
    { value: 'none', label: 'None' },
    { value: 'action', label: 'Action' },
    { value: 'bonus_action', label: 'Bonus Action' },
    { value: 'reaction', label: 'Reaction' },
    { value: 'free', label: 'Free' },
] as const;
const ACTION_TYPE_COLOR: Record<string, string> = {
    none: 'bg-stone-mid text-ash-light border-stone-border',
    action: 'bg-crimson/15 text-crimson-light border-crimson/30',
    bonus_action: 'bg-gold/15 text-gold border-gold/30',
    reaction: 'bg-info/15 text-info border-info/30',
    free: 'bg-emerald/15 text-[#27ae60] border-emerald/30',
};

const CURRENCIES = [
    { key: 'copper', label: 'CP', color: '#b5651d' },
    { key: 'silver', label: 'SP', color: '#c0c0c0' },
    { key: 'electrum', label: 'EP', color: '#b8b8ff' },
    { key: 'gold', label: 'GP', color: '#ffd700' },
    { key: 'platinum', label: 'PP', color: '#e5e4e2' },
] as const;

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
const addForm = "bg-stone border border-stone-border rounded-lg p-5 mb-4 flex flex-col gap-3";
const formRow = "flex gap-3 flex-wrap [&>*]:flex-1 [&>*]:min-w-[120px]";
const field = "flex flex-col gap-[5px]";
const fieldLabel = "text-[11px] text-ash uppercase tracking-wide";
const modalActions = "flex justify-end gap-2";
const cancelBtn = "bg-stone-mid border-none text-ash-light px-4 py-2 rounded-md text-[13px] hover:text-parchment";
const submitBtn = "bg-crimson border-none text-parchment px-4 py-2 rounded-md font-display text-[13px] font-semibold hover:bg-crimson-light";
const tabHeader = "flex items-center justify-between mb-3";
const addBtn = "flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment";
const emptyState = "text-center py-8 text-ash text-sm";
const itemBtnDanger = "bg-transparent border-none text-stone-border p-[5px] transition-colors duration-180 hover:text-crimson-light";
const statBox = "bg-stone-mid rounded-md p-3 text-center";
const statBoxLabel = "text-[10px] text-ash uppercase tracking-wide mb-1";
const statBoxValue = "text-lg font-semibold text-parchment";
const tag = "text-[10px] px-1.5 py-px rounded-sm";

function DetailModal({ onClose, children }: { onClose: () => void; children: ReactNode }) {
    return (
        <div className="fixed inset-0 bg-black/75 z-[100] flex items-center justify-center p-5" onClick={onClose}>
            <div className="bg-stone border border-stone-border rounded-lg p-8 w-full max-w-[480px] max-h-[85vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
                {children}
            </div>
        </div>
    );
}

function Subtabs<T extends string>({ tabs, active, onChange }: { tabs: { key: T; label: string }[]; active: T; onChange: (t: T) => void }) {
    return (
        <div className="shrink-0 flex gap-1.5 mb-3">
            {tabs.map(t => (
                <button key={t.key}
                    className={`px-3 py-1.5 rounded-md border text-xs font-medium transition-all duration-180 ${active === t.key ? 'bg-stone-mid border-gold-dim text-gold' : 'bg-transparent border-stone-border text-ash hover:text-parchment'}`}
                    onClick={() => onChange(t.key)}>
                    {t.label}
                </button>
            ))}
        </div>
    );
}

export function CharacterSheet() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [char, setChar] = useState<any>(null);
    const [inventory, setInventory] = useState<any[]>([]);
    const [spells, setSpells] = useState<any[]>([]);
    const [slots, setSlots] = useState<any[]>([]);
    const [features, setFeatures] = useState<any[]>([]);
    const [tab, setTab] = useState<Tab>('stats');
    const [loading, setLoading] = useState(true);
    const [hpInput, setHpInput] = useState('');
    const [hpMode, setHpMode] = useState<'damage' | 'heal' | 'temp'>('damage');
    const [showConditions, setShowConditions] = useState(false);
    const [showManageLevel, setShowManageLevel] = useState(false);

    const reload = useCallback(async () => {
        if (!id) return;
        const [c, inv, sp, sl, ft] = await Promise.all([
            api.getCharacter(id), api.listInventory(id), api.listSpells(id), api.listSpellSlots(id), api.listFeatures(id),
        ]);
        setChar(c); setInventory(inv || []); setSpells(sp || []); setSlots(sl || []); setFeatures(ft || []);
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

    const toggleInspiration = async () => {
        if (!id) return;
        const next = !char.inspiration;
        setChar((c: any) => ({ ...c, inspiration: next }));
        try {
            await api.updateCharacterInfo(id, { ...char, inspiration: next });
        } catch (err) {
            console.error('Failed to toggle inspiration', err);
            setChar((c: any) => ({ ...c, inspiration: !next }));
        }
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
        ['background', <IconStar size={15} />, 'Background'],
    ];

    return (
        <div className="h-full flex flex-col overflow-hidden">
            {/* Header */}
            <div className="shrink-0 flex items-center gap-4 py-2 px-8 bg-stone border-b border-stone-border flex-wrap">
                <button className="flex items-center gap-1 bg-transparent border-none text-ash text-[13px] transition-colors duration-180 hover:text-parchment" onClick={() => navigate('/characters')}>
                    <IconChevronLeft size={16} /> Characters
                </button>
                <div className="flex-1">
                    <h1 className="font-display text-lg font-bold text-parchment leading-tight">{char.name}</h1>
                    <p className="text-[11px] text-ash">Level {char.level} {char.race} {char.class} · {char.background} · {char.alignment}</p>
                </div>
                <div className="flex items-center gap-2 flex-wrap">
                    <button
                        className={`flex items-center gap-1.5 bg-transparent border-none transition-colors duration-180 ${char.inspiration ? 'text-gold hover:text-gold-light' : 'text-stone-border hover:text-ash'}`}
                        onClick={toggleInspiration} title="Toggle inspiration">
                        <span className="text-xl leading-none">★</span>
                        {char.inspiration && <span className="font-display text-xs font-semibold tracking-wide">Inspired</span>}
                    </button>
                    {char.is_npc && <span className="text-[11px] px-2 py-[3px] rounded-full bg-crimson/20 text-crimson-light">NPC</span>}
                    <span className="text-[11px] px-2 py-[3px] rounded-full bg-stone-mid text-ash-light">XP {char.xp?.toLocaleString()}</span>
                    <button className="flex items-center gap-1.5 bg-stone-mid border border-stone-border text-ash-light px-3 py-1.5 rounded-md text-xs font-medium transition-all duration-180 hover:border-gold-dim hover:text-gold" onClick={() => setShowManageLevel(true)}>
                        <IconPlus size={13} /> Manage Level
                    </button>
                </div>
            </div>

            {/* Combat bar */}
            <div className="shrink-0 py-2.5 px-8 bg-stone-mid border-b border-stone-border flex items-center justify-between gap-6 flex-wrap">
                <div className="flex items-center gap-6 flex-1 min-w-0 flex-wrap">
                    <div className="w-[170px] shrink-0">
                        <div className="text-[10px] font-bold tracking-widest text-ash mb-1">HIT POINTS</div>
                        <div className="w-full h-1.5 bg-stone-border rounded-sm mb-1 overflow-hidden">
                            <div className="h-full rounded-sm transition-[width,background] duration-[0.4s] ease-out" style={{ width: `${hpPct}%`, background: hpColor }} />
                        </div>
                        <div className="flex items-baseline gap-1.5">
                            <span style={{ color: hpColor, fontWeight: 700, fontSize: 18 }}>{char.current_hp}</span>
                            <span className="text-ash text-sm">/ {char.max_hp}</span>
                            {char.temp_hp > 0 && <span className="text-[12px] text-info">+{char.temp_hp}</span>}
                        </div>
                    </div>

                    <div className="flex items-center gap-1.5 shrink-0">
                        {(['damage', 'heal', 'temp'] as const).map(m => (
                            <button key={m}
                                className={`bg-stone border border-stone-border text-ash px-2.5 py-[5px] rounded-sm text-[11px] transition-all duration-180 ${hpMode === m ? 'bg-stone-light border-gold-dim text-parchment' : ''}`}
                                onClick={() => setHpMode(m)}>
                                {m === 'damage' ? '⚔ Dmg' : m === 'heal' ? '❤ Heal' : '🛡 Temp'}
                            </button>
                        ))}
                        <input className="w-[60px] text-center px-2 py-1.5" type="number" min={1} value={hpInput}
                            onChange={e => setHpInput(e.target.value)} placeholder="0"
                            onKeyDown={e => e.key === 'Enter' && handleHP()} />
                        <button className="bg-crimson border-none text-parchment px-3.5 py-1.5 rounded-sm text-[13px] transition-colors duration-180 hover:bg-crimson-light" onClick={handleHP}>Apply</button>
                    </div>

                    <div className="flex items-center gap-6 flex-1 justify-center flex-wrap">
                        {[
                            { l: 'AC', v: char.armor_class },
                            { l: 'INIT', v: fmt(mod(char.dexterity)) },
                            { l: 'SPEED', v: `${char.speed}ft` },
                            { l: 'PROF', v: `+${profBonus}` },
                            { l: `d${char.hit_dice_type}`, v: `${char.hit_dice_remaining} left` },
                        ].map(s => (
                            <div key={s.l} className="text-center">
                                <div className="font-display text-lg font-bold text-parchment">{s.v}</div>
                                <div className="text-[10px] text-ash tracking-wide mt-0.5">{s.l}</div>
                            </div>
                        ))}
                    </div>

                    {char.current_hp <= 0 && (
                        <div className="flex items-center gap-4 shrink-0">
                            <span className="text-[10px] font-bold tracking-wide text-crimson-light">DEATH SAVES</span>
                            {(['Successes', 'Failures'] as const).map(label => (
                                <div key={label} className="flex items-center gap-1.5">
                                    <span className="text-[10px] text-ash">{label === 'Successes' ? 'Succ' : 'Fail'}</span>
                                    {[0, 1, 2].map(i => {
                                        const count = label === 'Successes' ? char.death_save_successes : char.death_save_failures;
                                        const cls = i < count ? (label === 'Successes' ? 'bg-emerald border-emerald' : 'bg-crimson border-crimson') : '';
                                        return <button key={i} className={`w-[16px] h-[16px] rounded-full border-2 border-stone-border bg-transparent transition-all duration-180 ${cls}`}
                                            onClick={() => api.recordDeathSave(id!, label === 'Successes').then(reload)} />;
                                    })}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="flex items-center gap-1.5 shrink-0">
                    <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3 py-2 rounded-md text-xs transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => handleRest('short')}><IconMoon size={14} /> Short Rest</button>
                    <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3 py-2 rounded-md text-xs transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => handleRest('long')}><IconSun size={14} /> Long Rest</button>
                    <div className="relative">
                        <button
                            className={`flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3 py-2 rounded-md text-xs transition-all duration-180 hover:border-gold-dim hover:text-parchment ${char.conditions?.length ? 'border-crimson/50 text-crimson-light' : ''}`}
                            onClick={() => setShowConditions(v => !v)}>
                            ⚠ Conditions{char.conditions?.length ? ` (${char.conditions.length})` : ''}
                        </button>
                        {showConditions && (
                            <>
                                <div className="fixed inset-0 z-40" onClick={() => setShowConditions(false)} />
                                <div className="absolute top-full right-0 mt-1.5 z-50 bg-stone border border-stone-border rounded-md p-3 w-[260px] shadow-[0_12px_32px_rgba(0,0,0,0.5)] flex flex-wrap gap-1.5">
                                    {CONDITIONS.map(c => (
                                        <button key={c}
                                            className={`bg-stone-mid border border-stone-border text-ash px-2.5 py-[5px] rounded-full text-xs transition-all duration-180 hover:border-crimson hover:text-parchment ${char.conditions?.includes(c) ? 'bg-crimson/20 border-crimson text-crimson-light' : ''}`}
                                            onClick={() => toggleCond(c)}>
                                            {c}
                                        </button>
                                    ))}
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </div>

            {/* Tabs */}
            <div className="shrink-0 flex gap-0.5 px-8 bg-stone border-b border-stone-border">
                {TABS.map(([t, icon, label]) => (
                    <button key={t} className={`flex items-center gap-1.5 bg-transparent border-none text-ash px-[18px] py-2.5 text-[13px] font-medium border-b-2 border-transparent transition-all duration-180 hover:text-parchment ${tab === t ? 'text-gold border-gold' : ''}`} onClick={() => setTab(t)}>
                        {icon} {label}
                    </button>
                ))}
            </div>

            {/* Tab content — the only region that may ever scroll; the page itself never does */}
            <div className="flex-1 min-h-0 overflow-y-auto py-4 px-8">
                {tab === 'stats' && <StatsTab char={char} profBonus={profBonus} toggleCond={toggleCond} features={features} id={id!} reload={reload} />}
                {tab === 'inventory' && <InventoryTab id={id!} items={inventory} char={char} reload={reload} />}
                {tab === 'spells' && <SpellsTab id={id!} spells={spells} slots={slots} reload={reload} />}
                {tab === 'background' && <BackgroundTab char={char} id={id!} reload={reload} />}
            </div>

            {showManageLevel && (
                <ManageLevelModal char={char} id={id!} reload={reload} onClose={() => setShowManageLevel(false)} />
            )}
        </div>
    );
}

function ManageLevelModal({ char, id, reload, onClose }: { char: any; id: string; reload: () => void; onClose: () => void }) {
    const [level, setLevel] = useState<number>(char.level);
    const [maxHp, setMaxHp] = useState<number>(char.max_hp);
    const [abilities, setAbilities] = useState<Record<string, number>>(() =>
        Object.fromEntries(ABILITIES.map(a => [a, char[a]])));
    const [showAbilities, setShowAbilities] = useState(false);
    const [showAddFeature, setShowAddFeature] = useState(false);
    const [featureForm, setFeatureForm] = useState({ name: '', action_type: 'action', source: '', description: '' });
    const [confirming, setConfirming] = useState(false);
    const [saving, setSaving] = useState(false);

    const abilityChanges = ABILITIES.filter(a => abilities[a] !== char[a]);
    const dirty = level !== char.level || maxHp !== char.max_hp || abilityChanges.length > 0;

    const submitFeature = async (e: React.FormEvent) => {
        e.preventDefault();
        await api.createFeature(id, featureForm);
        setFeatureForm({ name: '', action_type: 'action', source: '', description: '' });
        setShowAddFeature(false);
        reload();
    };

    const applyChanges = async () => {
        setSaving(true);
        try {
            await api.updateLevel(id, { level, max_hp: maxHp });
            if (abilityChanges.length > 0) {
                await api.updateAbilityScores(id, Object.fromEntries(abilityChanges.map(a => [a, abilities[a]])));
            }
            await reload();
            onClose();
        } catch (err) {
            console.error('Failed to save level changes', err);
            setSaving(false);
            setConfirming(false);
        }
    };

    return (
        <DetailModal onClose={onClose}>
            {!confirming ? (
                <>
                    <h2 className="font-display text-xl font-bold text-parchment mb-5">Manage Level</h2>

                    <div className={`${field} mb-4`}>
                        <label className={fieldLabel}>Level</label>
                        <div className="flex items-center gap-3">
                            <button type="button" className="w-8 h-8 rounded-md bg-stone-mid border border-stone-border text-ash transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => setLevel(l => Math.max(1, l - 1))}>−</button>
                            <span className="font-display text-2xl font-bold text-parchment w-10 text-center">{level}</span>
                            <button type="button" className="w-8 h-8 rounded-md bg-stone-mid border border-stone-border text-ash transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => setLevel(l => Math.min(20, l + 1))}>+</button>
                        </div>
                    </div>

                    <div className={`${field} mb-5`}>
                        <label className={fieldLabel}>Max HP</label>
                        <div className="flex items-center gap-3">
                            <button type="button" className="w-8 h-8 rounded-md bg-stone-mid border border-stone-border text-ash transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => setMaxHp(h => Math.max(1, h - 1))}>−</button>
                            <input className="w-20 text-center" type="number" min={1} value={maxHp} onChange={e => setMaxHp(Math.max(1, +e.target.value))} />
                            <button type="button" className="w-8 h-8 rounded-md bg-stone-mid border border-stone-border text-ash transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => setMaxHp(h => h + 1)}>+</button>
                        </div>
                    </div>

                    <div className="flex gap-2 mb-6 flex-wrap">
                        <button type="button" className={addBtn} onClick={() => setShowAbilities(true)}>Change Ability Scores</button>
                        <button type="button" className={addBtn} onClick={() => setShowAddFeature(true)}><IconPlus size={14} /> Add Feature</button>
                    </div>

                    <div className={modalActions}>
                        <button className={cancelBtn} onClick={onClose}>Cancel</button>
                        <button className={submitBtn} disabled={!dirty} onClick={() => setConfirming(true)}>Save Changes</button>
                    </div>
                </>
            ) : (
                <>
                    <h2 className="font-display text-xl font-bold text-parchment mb-3">Confirm Changes</h2>
                    <div className="bg-stone-mid rounded-md p-4 mb-5 flex flex-col gap-1.5 text-sm">
                        {level !== char.level && <div className="text-ash-light">Level: <span className="text-parchment font-medium">{char.level} → {level}</span></div>}
                        {maxHp !== char.max_hp && <div className="text-ash-light">Max HP: <span className="text-parchment font-medium">{char.max_hp} → {maxHp}</span></div>}
                        {abilityChanges.map(a => (
                            <div key={a} className="text-ash-light">{a.slice(0, 3).toUpperCase()}: <span className="text-parchment font-medium">{char[a]} → {abilities[a]}</span></div>
                        ))}
                    </div>
                    <p className="text-[13px] text-ash mb-5">Are you sure you want to apply these changes?</p>
                    <div className={modalActions}>
                        <button className={cancelBtn} onClick={() => setConfirming(false)} disabled={saving}>Cancel</button>
                        <button className={submitBtn} onClick={applyChanges} disabled={saving}>{saving ? 'Saving…' : 'Confirm'}</button>
                    </div>
                </>
            )}

            {showAbilities && (
                <DetailModal onClose={() => setShowAbilities(false)}>
                    <h2 className="font-display text-lg font-bold text-parchment mb-4">Ability Scores</h2>
                    <div className="grid grid-cols-3 gap-3 mb-5">
                        {ABILITIES.map(a => (
                            <div key={a} className={field}>
                                <label className={fieldLabel}>{a.slice(0, 3).toUpperCase()}</label>
                                <input className="text-center" type="number" min={1} max={30} value={abilities[a]}
                                    onChange={e => setAbilities(v => ({ ...v, [a]: +e.target.value }))} />
                            </div>
                        ))}
                    </div>
                    <div className={modalActions}>
                        <button className={submitBtn} onClick={() => setShowAbilities(false)}>Done</button>
                    </div>
                </DetailModal>
            )}

            {showAddFeature && (
                <DetailModal onClose={() => setShowAddFeature(false)}>
                    <h2 className="font-display text-lg font-bold text-parchment mb-4">Add Feature</h2>
                    <form onSubmit={submitFeature} className="flex flex-col gap-3">
                        <div className={field}><label className={fieldLabel}>Name *</label><input value={featureForm.name} onChange={e => setFeatureForm(f => ({ ...f, name: e.target.value }))} required /></div>
                        <div className={field}><label className={fieldLabel}>Action Type</label>
                            <select value={featureForm.action_type} onChange={e => setFeatureForm(f => ({ ...f, action_type: e.target.value }))}>
                                {ACTION_TYPES.map(a => <option key={a.value} value={a.value}>{a.label}</option>)}
                            </select>
                        </div>
                        <div className={field}><label className={fieldLabel}>Source</label><input value={featureForm.source} onChange={e => setFeatureForm(f => ({ ...f, source: e.target.value }))} placeholder="Class, race, feat…" /></div>
                        <div className={field}><label className={fieldLabel}>Description</label><textarea value={featureForm.description} onChange={e => setFeatureForm(f => ({ ...f, description: e.target.value }))} rows={3} /></div>
                        <div className={modalActions}>
                            <button type="button" className={cancelBtn} onClick={() => setShowAddFeature(false)}>Cancel</button>
                            <button type="submit" className={submitBtn}>Add</button>
                        </div>
                    </form>
                </DetailModal>
            )}
        </DetailModal>
    );
}

function StatsTab({ char, profBonus, toggleCond, features, id, reload }: {
    char: any; profBonus: number; toggleCond: (c: string) => void;
    features: any[]; id: string; reload: () => void;
}) {
    const [box, setBox] = useState<'actions' | 'features' | 'conditions'>('actions');
    const [showAdd, setShowAdd] = useState(false);
    const [selected, setSelected] = useState<any>(null);
    const [form, setForm] = useState({ name: '', action_type: 'action', source: '', description: '' });

    const actionFeatures = features.filter(f => f.action_type !== 'none');
    const passiveFeatures = features.filter(f => f.action_type === 'none');

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        await api.createFeature(id, form);
        setShowAdd(false);
        setForm({ name: '', action_type: 'action', source: '', description: '' });
        reload();
    };

    const FeatureRow = ({ feature }: { feature: any }) => (
        <div className="flex items-center gap-2 py-1.5 px-2 rounded-sm bg-stone-mid cursor-pointer transition-colors duration-180 hover:bg-stone-light"
            onClick={() => setSelected(feature)}>
            {feature.action_type !== 'none' && (
                <span className={`${tag} border shrink-0 ${ACTION_TYPE_COLOR[feature.action_type]}`}>
                    {ACTION_TYPES.find(a => a.value === feature.action_type)?.label}
                </span>
            )}
            <span className="text-[13px] text-parchment truncate flex-1">{feature.name}</span>
        </div>
    );

    return (
        <div className="flex flex-col gap-4">
            <div className="grid grid-cols-[1fr_1fr_1.3fr] gap-4 max-[1050px]:grid-cols-1">
                <div className="flex flex-col gap-4">
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
                        <h3 className={panelTitle}>Proficiencies & Training</h3>
                        <div className="flex flex-col gap-2.5">
                            {[['Armor', char.training_armor], ['Weapons', char.training_weapons], ['Tools', char.training_tools], ['Languages', char.training_languages]].map(([l, v]) => (
                                <div key={l as string} className="bg-stone-mid rounded-md p-3">
                                    <div className="text-[10px] font-bold text-gold-dim tracking-wide uppercase mb-1">{l as string}</div>
                                    <div className="font-body text-sm text-parchment leading-relaxed">{((v as string[]) || []).join(', ') || '—'}</div>
                                </div>
                            ))}
                        </div>
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
                                    <span className="flex-1 text-[13px] text-ash-light">
                                        {sk.name} <span className="text-ash">({sk.ab.slice(0, 3).toUpperCase()})</span>
                                    </span>
                                    <span className="font-semibold text-[13px] text-parchment w-[30px] text-right">{fmt(total)}</span>
                                </div>
                            );
                        })}
                    </div>
                </div>

                <div className={panel}>
                    <div className="flex items-center justify-between mb-3">
                        <h3 className="font-display text-[13px] font-semibold text-gold tracking-wide uppercase flex items-center gap-1.5">
                            <IconZap size={13} /> Actions & Features
                        </h3>
                        <button className="bg-transparent border-none text-ash transition-colors duration-180 hover:text-gold" onClick={() => setShowAdd(true)} title="Add feature">
                            <IconPlus size={15} />
                        </button>
                    </div>
                    <Subtabs
                        tabs={[
                            { key: 'actions', label: `Actions (${actionFeatures.length})` },
                            { key: 'features', label: `Features (${passiveFeatures.length})` },
                            { key: 'conditions', label: 'Conditions' },
                        ]}
                        active={box} onChange={setBox} />
                    {box === 'actions' && (
                        <div className="flex flex-col gap-1">
                            {actionFeatures.map(f => <FeatureRow key={f.id} feature={f} />)}
                            {actionFeatures.length === 0 && <div className={emptyState}>No actions yet.</div>}
                        </div>
                    )}
                    {box === 'features' && (
                        <div className="flex flex-col gap-1">
                            {passiveFeatures.map(f => <FeatureRow key={f.id} feature={f} />)}
                            {passiveFeatures.length === 0 && <div className={emptyState}>No features yet.</div>}
                        </div>
                    )}
                    {box === 'conditions' && (
                        <div className="flex flex-wrap content-start gap-1.5">
                            {CONDITIONS.map(c => (
                                <button key={c}
                                    className={`bg-stone-mid border border-stone-border text-ash px-2.5 py-[5px] rounded-full text-xs transition-all duration-180 hover:border-crimson hover:text-parchment ${char.conditions?.includes(c) ? 'bg-crimson/20 border-crimson text-crimson-light' : ''}`}
                                    onClick={() => toggleCond(c)}>
                                    {c}
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            {showAdd && (
                <DetailModal onClose={() => setShowAdd(false)}>
                    <h2 className="font-display text-xl font-bold text-parchment mb-4">Add Feature</h2>
                    <form onSubmit={submit} className="flex flex-col gap-3">
                        <div className={field}><label className={fieldLabel}>Name *</label><input value={form.name} onChange={e => setForm(f => ({ ...f, name: e.target.value }))} required /></div>
                        <div className={field}><label className={fieldLabel}>Action Type</label>
                            <select value={form.action_type} onChange={e => setForm(f => ({ ...f, action_type: e.target.value }))}>
                                {ACTION_TYPES.map(a => <option key={a.value} value={a.value}>{a.label}</option>)}
                            </select>
                        </div>
                        <div className={field}><label className={fieldLabel}>Source</label><input value={form.source} onChange={e => setForm(f => ({ ...f, source: e.target.value }))} placeholder="Class, race, feat…" /></div>
                        <div className={field}><label className={fieldLabel}>Description</label><textarea value={form.description} onChange={e => setForm(f => ({ ...f, description: e.target.value }))} rows={3} /></div>
                        <div className={modalActions}>
                            <button type="button" className={cancelBtn} onClick={() => setShowAdd(false)}>Cancel</button>
                            <button type="submit" className={submitBtn}>Add</button>
                        </div>
                    </form>
                </DetailModal>
            )}

            {selected && (
                <DetailModal onClose={() => setSelected(null)}>
                    <div className="flex items-start justify-between mb-1">
                        <h2 className="font-display text-xl font-bold text-parchment pr-4">{selected.name}</h2>
                        <button className="text-ash text-xl leading-none hover:text-parchment" onClick={() => setSelected(null)}>×</button>
                    </div>
                    <div className="flex gap-2 flex-wrap mb-4">
                        <span className={`${tag} border ${ACTION_TYPE_COLOR[selected.action_type] || ACTION_TYPE_COLOR.none}`}>
                            {ACTION_TYPES.find(a => a.value === selected.action_type)?.label || selected.action_type}
                        </span>
                        {selected.source && <span className={`${tag} bg-stone-mid text-ash-light border border-stone-border`}>{selected.source}</span>}
                    </div>
                    {selected.description && <p className="font-body text-[15px] text-ash-light leading-relaxed mb-5">{selected.description}</p>}
                    <div className={modalActions}>
                        <button className={cancelBtn} onClick={() => { api.deleteFeature(id, selected.id).then(reload); setSelected(null); }}>Delete</button>
                    </div>
                </DetailModal>
            )}
        </div>
    );
}

function InventoryTab({ id, items, char, reload }: { id: string; items: any[]; char: any; reload: () => void }) {
    const [sub, setSub] = useState<'equipped' | 'backpack' | 'attuned'>('equipped');
    const [showAdd, setShowAdd] = useState(false);
    const [selected, setSelected] = useState<any>(null);
    const [form, setForm] = useState({ name: '', quantity: 1, weight: 0, value: 0, description: '', is_equipped: false, requires_attunement: false });

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        await api.createInventoryItem(id, form);
        setShowAdd(false);
        setForm({ name: '', quantity: 1, weight: 0, value: 0, description: '', is_equipped: false, requires_attunement: false });
        reload();
    };

    const equipped = items.filter(i => i.is_equipped);
    const backpack = items.filter(i => !i.is_equipped);
    const attuned = items.filter(i => i.is_attuned);
    const list = sub === 'equipped' ? equipped : sub === 'attuned' ? attuned : backpack;

    const adjustCurrency = (key: string, delta: number) => {
        const next = Math.max(0, (char[key] || 0) + delta);
        api.updateCurrency(id, { [key]: next }).then(reload);
    };

    return (
        <div className="h-full grid grid-cols-[1fr_220px] gap-4 max-[850px]:grid-cols-1 max-[850px]:overflow-y-auto">
            <div className="flex flex-col min-h-0">
                <Subtabs
                    tabs={[
                        { key: 'equipped', label: `Equipped (${equipped.length})` },
                        { key: 'backpack', label: `Backpack (${backpack.length})` },
                        { key: 'attuned', label: `Attuned (${attuned.length})` },
                    ]}
                    active={sub} onChange={setSub} />

                <div className={tabHeader}>
                    <span className="text-[13px] text-ash">{list.reduce((s, i) => s + i.weight * i.quantity, 0)} lbs</span>
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

                <div className="flex flex-col gap-1.5 flex-1 min-h-0 overflow-y-auto">
                    {list.map(item => (
                        <div key={item.id}
                            className={`flex items-center gap-3 bg-stone border rounded-md py-3 px-4 cursor-pointer transition-colors duration-180 hover:border-gold-dim ${item.is_equipped ? 'border-gold-dim' : 'border-stone-border'}`}
                            onClick={() => setSelected(item)}>
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
                            <div className="flex gap-1.5" onClick={e => e.stopPropagation()}>
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
                    {list.length === 0 && <div className={emptyState}>Nothing here.</div>}
                </div>
            </div>

            <div className={`${panel} h-fit`}>
                <h3 className={panelTitle}>Currency</h3>
                <div className="flex flex-col gap-2">
                    {CURRENCIES.map(c => (
                        <div key={c.key} className="flex items-center justify-between gap-2 bg-stone-mid rounded-md px-2.5 py-2">
                            <span className="text-[11px] font-bold w-7" style={{ color: c.color }}>{c.label}</span>
                            <div className="flex items-center gap-1.5">
                                <button className="w-6 h-6 rounded-sm bg-stone border border-stone-border text-ash text-sm leading-none transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => adjustCurrency(c.key, -1)}>−</button>
                                <span className="w-9 text-center text-sm font-semibold text-parchment">{char[c.key]}</span>
                                <button className="w-6 h-6 rounded-sm bg-stone border border-stone-border text-ash text-sm leading-none transition-colors duration-180 hover:text-parchment hover:border-gold-dim" onClick={() => adjustCurrency(c.key, 1)}>+</button>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {selected && (
                <DetailModal onClose={() => setSelected(null)}>
                    <div className="flex items-start justify-between mb-1">
                        <h2 className="font-display text-xl font-bold text-parchment pr-4">{selected.name}</h2>
                        <button className="text-ash text-xl leading-none hover:text-parchment" onClick={() => setSelected(null)}>×</button>
                    </div>
                    <div className="flex gap-2 flex-wrap mb-4">
                        {selected.is_equipped && <span className={`${tag} bg-gold/15 text-gold border border-gold/30`}>Equipped</span>}
                        {selected.is_attuned && <span className={`${tag} bg-info/15 text-info border border-info/30`}>Attuned</span>}
                        {selected.requires_attunement && !selected.is_attuned && <span className={`${tag} bg-stone-mid text-ash border border-stone-border`}>Requires Attunement</span>}
                    </div>
                    <div className="grid grid-cols-3 gap-3 mb-4">
                        <div className={statBox}><div className={statBoxLabel}>Quantity</div><div className={statBoxValue}>{selected.quantity}</div></div>
                        <div className={statBox}><div className={statBoxLabel}>Weight</div><div className={statBoxValue}>{selected.weight}lb</div></div>
                        <div className={statBox}><div className={statBoxLabel}>Value</div><div className={statBoxValue}>{selected.value}gp</div></div>
                    </div>
                    {selected.description && <p className="font-body text-[15px] text-ash-light leading-relaxed mb-5">{selected.description}</p>}
                    <div className={modalActions}>
                        <button className={cancelBtn} onClick={() => { api.deleteInventoryItem(id, selected.id).then(reload); setSelected(null); }}>Delete</button>
                        <button className={submitBtn} onClick={() => { api.updateInventoryItem(id, selected.id, { ...selected, is_equipped: !selected.is_equipped }).then(reload); setSelected(null); }}>
                            {selected.is_equipped ? 'Unequip' : 'Equip'}
                        </button>
                    </div>
                </DetailModal>
            )}
        </div>
    );
}

const SCHOOLS = ['Abjuration', 'Conjuration', 'Divination', 'Enchantment', 'Evocation', 'Illusion', 'Necromancy', 'Transmutation'];

function SpellsTab({ id, spells, slots, reload }: { id: string; spells: any[]; slots: any[]; reload: () => void }) {
    const [sub, setSub] = useState<'prepared' | 'unprepared'>('prepared');
    const [showAdd, setShowAdd] = useState(false);
    const [selected, setSelected] = useState<any>(null);
    const [form, setForm] = useState({
        name: '', level: 1, school: 'Evocation', casting_time: '1 action',
        range: '60 feet', components: 'V, S', duration: 'Instantaneous', description: '', is_prepared: true,
    });

    const prepared = spells.filter(s => s.is_prepared);
    const unprepared = spells.filter(s => !s.is_prepared);
    const filtered = sub === 'prepared' ? prepared : unprepared;

    const byLevel = Array.from({ length: 10 }, (_, i) => ({
        level: i,
        spells: filtered.filter(s => s.level === i),
        slot: i > 0 ? slots.find(s => s.spell_level === i) : undefined,
    })).filter(g => g.spells.length > 0);

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        await api.createSpell(id, form);
        setShowAdd(false); reload();
    };

    return (
        <div className="h-full flex flex-col min-h-0">
            <Subtabs
                tabs={[
                    { key: 'prepared', label: `Prepared (${prepared.length})` },
                    { key: 'unprepared', label: `Unprepared (${unprepared.length})` },
                ]}
                active={sub} onChange={setSub} />

            <div className={tabHeader}>
                <span className="text-[13px] text-ash">{filtered.length} spell{filtered.length !== 1 ? 's' : ''}</span>
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

            <div className="flex-1 min-h-0 overflow-y-auto">
                {byLevel.map(group => (
                    <div key={group.level} className="mb-4">
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
                            <div key={spell.id}
                                className="flex items-start gap-2.5 py-2.5 px-3 rounded-md mb-1 bg-stone border border-stone-border cursor-pointer transition-colors duration-180 hover:border-gold-dim"
                                onClick={() => setSelected(spell)}>
                                <button className={`w-[22px] h-[22px] rounded-full border-2 border-stone-border bg-transparent flex items-center justify-center transition-all duration-180 shrink-0 mt-0.5 hover:border-gold-dim ${spell.is_prepared ? 'bg-gold-dim border-gold text-ink' : 'text-transparent'}`}
                                    onClick={e => { e.stopPropagation(); api.toggleSpellPrepared(id, spell.id).then(reload); }}>
                                    <IconCheck size={12} />
                                </button>
                                <div className="flex-1">
                                    <div className="text-sm font-medium text-parchment">{spell.name}</div>
                                    <div className="text-xs text-ash mt-0.5">{spell.school} · {spell.casting_time} · {spell.range}</div>
                                </div>
                                <button className={itemBtnDanger} onClick={e => { e.stopPropagation(); api.deleteSpell(id, spell.id).then(reload); }}>
                                    <IconTrash size={13} />
                                </button>
                            </div>
                        ))}
                    </div>
                ))}
                {byLevel.length === 0 && <div className={emptyState}>No {sub} spells.</div>}
            </div>

            {selected && (
                <DetailModal onClose={() => setSelected(null)}>
                    <div className="flex items-start justify-between mb-1">
                        <h2 className="font-display text-xl font-bold text-parchment pr-4">{selected.name}</h2>
                        <button className="text-ash text-xl leading-none hover:text-parchment" onClick={() => setSelected(null)}>×</button>
                    </div>
                    <div className="flex gap-2 flex-wrap mb-4">
                        <span className={`${tag} bg-stone-mid text-ash-light border border-stone-border`}>{selected.level === 0 ? 'Cantrip' : `Level ${selected.level}`}</span>
                        <span className={`${tag} bg-stone-mid text-ash-light border border-stone-border`}>{selected.school}</span>
                        {selected.is_prepared && <span className={`${tag} bg-gold/15 text-gold border border-gold/30`}>Prepared</span>}
                    </div>
                    <div className="grid grid-cols-2 gap-3 mb-4">
                        <div className={statBox}><div className={statBoxLabel}>Casting Time</div><div className={statBoxValue}>{selected.casting_time}</div></div>
                        <div className={statBox}><div className={statBoxLabel}>Range</div><div className={statBoxValue}>{selected.range}</div></div>
                        <div className={statBox}><div className={statBoxLabel}>Components</div><div className={statBoxValue}>{selected.components}</div></div>
                        <div className={statBox}><div className={statBoxLabel}>Duration</div><div className={statBoxValue}>{selected.duration}</div></div>
                    </div>
                    {selected.description && <p className="font-body text-[15px] text-ash-light leading-relaxed mb-5">{selected.description}</p>}
                    <div className={modalActions}>
                        <button className={cancelBtn} onClick={() => { api.deleteSpell(id, selected.id).then(reload); setSelected(null); }}>Delete</button>
                        <button className={submitBtn} onClick={() => { api.toggleSpellPrepared(id, selected.id).then(reload); setSelected(null); }}>
                            {selected.is_prepared ? 'Unprepare' : 'Prepare'}
                        </button>
                    </div>
                </DetailModal>
            )}
        </div>
    );
}

function BackgroundTab({ char, id, reload }: { char: any; id: string; reload: () => void }) {
    const [sub, setSub] = useState<'traits' | 'notes'>('traits');
    const [notes, setNotes] = useState(char.notes || '');
    const [saved, setSaved] = useState(false);

    const save = async () => {
        await api.updateCharacterInfo(id, { ...char, notes });
        setSaved(true); setTimeout(() => setSaved(false), 2000); reload();
    };

    return (
        <div className="h-full flex flex-col min-h-0">
            <Subtabs tabs={[{ key: 'traits', label: 'Character Traits' }, { key: 'notes', label: 'Notes' }]} active={sub} onChange={setSub} />

            {sub === 'traits' && (
                <div className={`${panel} overflow-y-auto`}>
                    <h3 className={panelTitle}>Character Traits</h3>
                    <div className="grid grid-cols-2 gap-4 max-[600px]:grid-cols-1">
                        {[['Personality', char.personality_traits], ['Ideals', char.ideals], ['Bonds', char.bonds], ['Flaws', char.flaws]].map(([l, v]) => (
                            <div key={l as string} className="bg-stone-mid rounded-md p-3.5">
                                <div className="text-[10px] font-bold text-gold-dim tracking-wide uppercase mb-1.5">{l as string}</div>
                                <div className="font-body text-[15px] text-parchment leading-relaxed">{(v as string) || '—'}</div>
                            </div>
                        ))}
                    </div>
                </div>
            )}

            {sub === 'notes' && (
                <div className="flex flex-col flex-1 min-h-0">
                    <div className={tabHeader}>
                        <span />
                        <button className={addBtn} onClick={save}>{saved ? '✓ Saved' : 'Save Notes'}</button>
                    </div>
                    <textarea
                        className="w-full flex-1 min-h-0 bg-stone border border-stone-border rounded-lg p-5 font-body text-base text-parchment leading-[1.7] resize-none"
                        value={notes} onChange={e => setNotes(e.target.value)}
                        placeholder="Session notes, backstory, important NPCs…" />
                </div>
            )}
        </div>
    );
}
