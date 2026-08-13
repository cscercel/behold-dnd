import { useState, useEffect } from 'react';
import * as api from '../lib/api';
import { IconPlus, IconSword, IconSkip, IconStop, IconPlay, IconTrash, IconChevronDown, IconChevronUp, IconZap } from '../components/Icon';

const field = "flex flex-col gap-1.5";
const fieldLabel = "text-[11px] text-ash uppercase tracking-wide";

const CONDITIONS = ['Blinded', 'Charmed', 'Deafened', 'Exhaustion', 'Frightened', 'Grappled',
    'Incapacitated', 'Invisible', 'Paralyzed', 'Petrified', 'Poisoned', 'Prone', 'Restrained', 'Stunned', 'Unconscious'];

const rollD20 = () => Math.floor(Math.random() * 20) + 1;

export function Combat() {
    const [encounters, setEncounters] = useState<any[]>([]);
    const [active, setActive] = useState<any>(null);
    const [participants, setParticipants] = useState<any[]>([]);
    const [characters, setCharacters] = useState<any[]>([]);
    const [selectedId, setSelectedId] = useState('');
    const [newName, setNewName] = useState('');
    const [loading, setLoading] = useState(true);
    const [showAddPart, setShowAddPart] = useState(false);
    const [showRoll, setShowRoll] = useState(false);
    const [partForm, setPartForm] = useState({ character_id: '', name: '', initiative_bonus: 0, max_hp: 10, current_hp: 10, armor_class: 10, speed: 30 });

    const reload = async () => {
        const [enc, chars] = await Promise.all([api.listEncounters(), api.listCharacters()]);
        setEncounters(enc || []); setCharacters(chars || []);
        let act: any = null;
        try { act = await api.getActiveEncounter(); } catch { act = null; }
        setActive(act);
        const targetId = act?.id || selectedId;
        if (targetId) {
            const parts = await api.listParticipants(targetId);
            setParticipants((parts || []).sort((a: any, b: any) => b.initiative - a.initiative));
        } else {
            setParticipants([]);
        }
    };

    useEffect(() => { reload().finally(() => setLoading(false)); }, []);

    const loadEnc = async (id: string) => {
        setSelectedId(id);
        const parts = await api.listParticipants(id);
        setParticipants((parts || []).sort((a: any, b: any) => b.initiative - a.initiative));
    };

    const createEnc = () => {
        if (!newName.trim()) return;
        api.createEncounter(newName.trim()).then(() => { setNewName(''); reload(); });
    };

    const addPart = async (e: React.FormEvent) => {
        e.preventDefault();
        const encId = active?.id || selectedId;
        if (!encId) return;
        await api.addParticipant(encId, { ...partForm, initiative: Math.max(1, partForm.initiative_bonus) });
        setShowAddPart(false);
        setPartForm({ character_id: '', name: '', initiative_bonus: 0, max_hp: 10, current_hp: 10, armor_class: 10, speed: 30 });
        reload();
    };

    const selectChar = (charId: string) => {
        const c = characters.find((ch: any) => ch.id === charId);
        if (c) setPartForm(f => ({ ...f, character_id: charId, name: c.name, max_hp: c.max_hp, current_hp: c.current_hp, armor_class: c.armor_class, speed: c.speed }));
        else setPartForm(f => ({ ...f, character_id: charId }));
    };

    const encId = active?.id || selectedId;
    const currentEnc = active || encounters.find((e: any) => e.id === selectedId);

    if (loading) return <div className="flex items-center justify-center h-[200px] text-ash">Loading…</div>;

    return (
        <div className="flex h-screen overflow-hidden">
            <div className="w-[260px] shrink-0 bg-stone border-r border-stone-border flex flex-col p-5 gap-3 overflow-y-auto">
                <h2 className="font-display text-[15px] font-bold text-parchment">Encounters</h2>
                <div className="flex gap-1.5">
                    <input className="flex-1" value={newName} onChange={e => setNewName(e.target.value)}
                        placeholder="Encounter name…" onKeyDown={e => e.key === 'Enter' && createEnc()} />
                    <button className="bg-crimson border-none text-parchment px-2.5 py-2 rounded-sm hover:bg-crimson-light" onClick={createEnc}><IconPlus size={14} /></button>
                </div>
                <div className="flex flex-col gap-1 flex-1">
                    {encounters.map((enc: any) => (
                        <div key={enc.id}
                            className={`px-3 py-2.5 rounded-md border cursor-pointer transition-all duration-180 flex items-center justify-between hover:bg-stone-mid ${enc.id === encId ? 'bg-stone-mid border-stone-border' : 'border-transparent'} ${enc.is_active ? 'border-crimson' : ''}`}
                            onClick={() => loadEnc(enc.id)}>
                            <div>
                                <div className="text-[13px] font-medium text-parchment">{enc.name}</div>
                                <div className="text-[11px] text-ash mt-0.5">{enc.is_active ? `Round ${enc.round} · LIVE` : 'Inactive'}</div>
                            </div>
                            {enc.is_active && <div className="w-2 h-2 rounded-full bg-crimson-light shadow-[0_0_6px_#c0392b]" />}
                        </div>
                    ))}
                    {encounters.length === 0 && <div className="text-[13px] text-ash text-center py-5">No encounters yet.</div>}
                </div>
            </div>

            <div className="flex-1 overflow-y-auto py-6 px-8">
                {currentEnc ? (
                    <>
                        <div className="flex items-center justify-between mb-5 flex-wrap gap-3">
                            <div>
                                <h1 className="font-display text-2xl font-bold text-parchment">{currentEnc.name}</h1>
                                {currentEnc.is_active && <span className="text-xs bg-crimson/20 text-crimson-light border border-crimson/30 px-2.5 py-0.5 rounded-full mt-1 inline-block">Round {currentEnc.round}</span>}
                            </div>
                            <div className="flex gap-2 flex-wrap">
                                {!currentEnc.is_active
                                    ? <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => api.startEncounter(currentEnc.id).then(reload)}><IconPlay size={15} /> Start</button>
                                    : <>
                                        <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => api.nextRound(currentEnc.id).then(reload)}><IconSkip size={15} /> Next Round</button>
                                        <button className="flex items-center gap-1.5 bg-crimson/15 border border-crimson/30 text-crimson-light px-3.5 py-2 rounded-md text-[13px] hover:bg-crimson/25" onClick={() => api.endEncounter(currentEnc.id).then(reload)}><IconStop size={15} /> End</button>
                                    </>
                                }
                                {participants.length > 0 && (
                                    <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => setShowRoll(true)}>
                                        <IconZap size={15} /> Roll Initiative
                                    </button>
                                )}
                                <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash-light px-3.5 py-2 rounded-md text-[13px] transition-all duration-180 hover:border-gold-dim hover:text-parchment" onClick={() => setShowAddPart(v => !v)}><IconPlus size={15} /> Add Combatant</button>
                            </div>
                        </div>

                        {showAddPart && (
                            <form onSubmit={addPart} className="bg-stone border border-stone-border rounded-lg p-5 mb-5 flex flex-col gap-3">
                                <div className="flex gap-3 flex-wrap [&>*]:flex-1 [&>*]:min-w-[120px]">
                                    <div className={field}>
                                        <label className={fieldLabel}>From Characters</label>
                                        <select value={partForm.character_id} onChange={e => selectChar(e.target.value)}>
                                            <option value="">— Custom NPC —</option>
                                            {characters.map((c: any) => <option key={c.id} value={c.id}>{c.name}</option>)}
                                        </select>
                                    </div>
                                    <div className={field}><label className={fieldLabel}>Name *</label><input value={partForm.name} onChange={e => setPartForm(f => ({ ...f, name: e.target.value }))} required /></div>
                                    <div className={field}>
                                        <label className={fieldLabel}>Initiative Bonus</label>
                                        <input type="number" value={partForm.initiative_bonus} onChange={e => setPartForm(f => ({ ...f, initiative_bonus: +e.target.value }))} />
                                    </div>
                                    <div className={field}><label className={fieldLabel}>Max HP</label><input type="number" min={1} value={partForm.max_hp} onChange={e => setPartForm(f => ({ ...f, max_hp: +e.target.value, current_hp: +e.target.value }))} /></div>
                                    <div className={field}><label className={fieldLabel}>AC</label><input type="number" min={1} value={partForm.armor_class} onChange={e => setPartForm(f => ({ ...f, armor_class: +e.target.value }))} /></div>
                                </div>
                                <p className="text-[11px] text-ash">The bonus is added to a d20 roll — yours if a player, auto-rolled otherwise — when you click "Roll Initiative".</p>
                                <div className="flex justify-end gap-2">
                                    <button type="button" className="bg-stone-mid border-none text-ash-light px-4 py-2 rounded-md text-[13px]" onClick={() => setShowAddPart(false)}>Cancel</button>
                                    <button type="submit" className="bg-crimson border-none text-parchment px-4 py-2 rounded-md font-display text-[13px] font-semibold">Add to Battle</button>
                                </div>
                            </form>
                        )}

                        <div className="flex flex-col gap-1.5">
                            {participants.length === 0 && (
                                <div className="text-center py-[60px] px-5 text-ash"><IconSword size={40} className="mx-auto mb-3 block opacity-30" /><p>No combatants yet.</p></div>
                            )}
                            {participants.map((p: any) => (
                                <ParticipantRow key={p.id} p={p} encId={encId}
                                    onDmg={n => api.participantDamage(encId, p.id, n).then(reload)}
                                    onHeal={n => api.participantHeal(encId, p.id, n).then(reload)}
                                    onTempHP={n => api.participantTempHP(encId, p.id, n).then(reload)}
                                    onConditions={cs => api.updateParticipantConditions(encId, p.id, cs).then(reload)}
                                    onToggleConcentration={() => api.toggleParticipantConcentration(encId, p.id).then(reload)}
                                    onDeactivate={() => api.deactivateParticipant(encId, p.id).then(reload)}
                                    onRemove={() => api.removeParticipant(encId, p.id).then(reload)}
                                />
                            ))}
                        </div>
                    </>
                ) : (
                    <div className="flex flex-col items-center justify-center h-[60vh] text-ash gap-3">
                        <IconSword size={64} className="opacity-20" />
                        <p>Select or create an encounter to begin</p>
                    </div>
                )}
            </div>

            {showRoll && (
                <RollInitiativeModal
                    encId={encId}
                    participants={participants}
                    characters={characters}
                    onClose={() => setShowRoll(false)}
                    onDone={reload}
                />
            )}
        </div>
    );
}

function RollInitiativeModal({ encId, participants, characters, onClose, onDone }: {
    encId: string; participants: any[]; characters: any[];
    onClose: () => void; onDone: () => void;
}) {
    const isPlayer = (p: any) => {
        const c = characters.find((ch: any) => ch.id === p.character_id);
        return !!c && !c.is_npc;
    };
    const players = participants.filter(isPlayer);
    const npcs = participants.filter(p => !isPlayer(p));

    const [rolls, setRolls] = useState<Record<string, string>>(() =>
        Object.fromEntries(players.map(p => [p.id, ''])));
    const [rolling, setRolling] = useState(false);

    const allFilled = players.every(p => rolls[p.id]?.trim() !== '');

    const submit = async () => {
        setRolling(true);
        try {
            await Promise.all(participants.map(p => {
                const bonus = p.initiative_bonus ?? 0;
                const d20 = isPlayer(p) ? Math.max(1, Math.min(20, parseInt(rolls[p.id]) || 1)) : rollD20();
                const total = Math.max(1, d20 + bonus);
                return api.updateParticipantInitiative(encId, p.id, total);
            }));
            onDone();
            onClose();
        } finally {
            setRolling(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/75 z-[100] flex items-center justify-center p-5" onClick={onClose}>
            <div className="bg-stone border border-stone-border rounded-lg p-7 w-full max-w-[440px] max-h-[85vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
                <h2 className="font-display text-lg font-bold text-parchment mb-1">Roll Initiative</h2>
                <p className="text-[13px] text-ash mb-5">Enter each player's physical d20 roll. Non-player combatants are rolled automatically.</p>

                {players.length > 0 && (
                    <div className="flex flex-col gap-2.5 mb-5">
                        {players.map(p => (
                            <div key={p.id} className="flex items-center gap-3">
                                <span className="flex-1 text-sm text-parchment">{p.name}</span>
                                <span className="text-[11px] text-ash">+{p.initiative_bonus ?? 0}</span>
                                <input className="w-16 text-center" type="number" min={1} max={20}
                                    value={rolls[p.id]} placeholder="d20"
                                    onChange={e => setRolls(r => ({ ...r, [p.id]: e.target.value }))} />
                            </div>
                        ))}
                    </div>
                )}

                {npcs.length > 0 && (
                    <div className="bg-stone-mid rounded-md p-3 mb-5">
                        <div className="text-[11px] text-ash-light">
                            {npcs.length} combatant{npcs.length !== 1 ? 's' : ''} will be rolled automatically: {npcs.map(n => n.name).join(', ')}
                        </div>
                    </div>
                )}

                {players.length === 0 && npcs.length === 0 && (
                    <div className="text-center py-6 text-ash text-sm">No combatants to roll.</div>
                )}

                <div className="flex justify-end gap-2">
                    <button className="bg-stone-mid border-none text-ash-light px-4 py-2 rounded-md text-[13px] hover:text-parchment" onClick={onClose} disabled={rolling}>Cancel</button>
                    <button className="bg-crimson border-none text-parchment px-4 py-2 rounded-md font-display text-[13px] font-semibold hover:bg-crimson-light disabled:opacity-50 disabled:cursor-not-allowed"
                        onClick={submit} disabled={rolling || !allFilled}>
                        {rolling ? 'Rolling…' : 'Roll'}
                    </button>
                </div>
            </div>
        </div>
    );
}

function ParticipantRow({ p, onDmg, onHeal, onTempHP, onConditions, onToggleConcentration, onDeactivate, onRemove }: {
    p: any; encId: string;
    onDmg: (n: number) => void;
    onHeal: (n: number) => void;
    onTempHP: (n: number) => void;
    onConditions: (conditions: string[]) => void;
    onToggleConcentration: () => void;
    onDeactivate: () => void;
    onRemove: () => void;
}) {
    const [expanded, setExpanded] = useState(false);
    const [amount, setAmount] = useState('');
    const [tempAmount, setTempAmount] = useState('');
    const hpPct = Math.max(0, Math.min(100, (p.current_hp / p.max_hp) * 100));
    const hpColor = hpPct > 50 ? '#27ae60' : hpPct > 25 ? '#f39c12' : '#c0392b';

    const toggleCondition = (c: string) => {
        const next = (p.conditions || []).includes(c)
            ? p.conditions.filter((x: string) => x !== c)
            : [...(p.conditions || []), c];
        onConditions(next);
    };

    return (
        <div className={`flex items-stretch bg-stone border border-stone-border rounded-lg overflow-hidden ${!p.is_active ? 'opacity-40' : ''}`}>
            <div className="w-12 shrink-0 bg-stone-mid flex flex-col items-center justify-center">
                <span className="font-display text-lg font-bold text-gold leading-none">{p.initiative}</span>
                {!!p.initiative_bonus && <span className="text-[10px] text-ash mt-0.5">+{p.initiative_bonus}</span>}
            </div>
            <div className="flex-1">
                <div className="flex items-center gap-3 px-4 py-3 cursor-pointer hover:bg-white/[0.02]" onClick={() => setExpanded(v => !v)}>
                    <div>
                        <div className="text-[15px] font-medium text-parchment flex-1">{p.name}</div>
                        {p.conditions?.length > 0 && (
                            <div className="flex gap-1 flex-wrap mt-0.5">
                                {p.conditions.map((c: string) => <span key={c} className="text-[10px] bg-crimson/20 text-crimson-light rounded-sm px-1.5 py-0.5">{c}</span>)}
                            </div>
                        )}
                    </div>
                    <div className="flex items-center gap-1.5 flex-1">
                        <div className="w-20 h-1.5 bg-stone-border rounded-sm overflow-hidden">
                            <div style={{ width: `${hpPct}%`, background: hpColor, height: '100%', borderRadius: 3 }} />
                        </div>
                        <span className="font-bold text-base" style={{ color: hpColor }}>{p.current_hp}</span>
                        <span className="text-xs text-ash">/{p.max_hp}</span>
                        {p.temp_hp > 0 && <span className="text-xs text-info">+{p.temp_hp}</span>}
                    </div>
                    <div className="flex gap-2.5 text-xs text-ash items-center">
                        <span>AC {p.armor_class}</span>
                        {p.concentration && <span className="text-[#9b59b6]" title="Concentrating">⊛</span>}
                        {!p.is_active && <span className="text-crimson-light" title="Inactive">✕</span>}
                    </div>
                    {expanded ? <IconChevronUp size={14} className="text-ash shrink-0" /> : <IconChevronDown size={14} className="text-ash shrink-0" />}
                </div>
                {expanded && (
                    <div className="px-4 py-3 border-t border-stone-border flex flex-col gap-3 bg-stone-mid">
                        <div className="flex items-center gap-2.5 flex-wrap">
                            <input className="w-20 text-center px-2 py-1.5" type="number" min={1} value={amount}
                                onChange={e => setAmount(e.target.value)} placeholder="Amount" />
                            <button className="bg-crimson/20 border border-crimson/30 text-crimson-light px-3 py-1.5 rounded-sm text-xs hover:bg-crimson/35" onClick={() => { if (amount) { onDmg(+amount); setAmount(''); } }}>⚔ Damage</button>
                            <button className="bg-emerald/15 border border-emerald/30 text-[#27ae60] px-3 py-1.5 rounded-sm text-xs hover:bg-emerald/25" onClick={() => { if (amount) { onHeal(+amount); setAmount(''); } }}>❤ Heal</button>

                            <input className="w-20 text-center px-2 py-1.5" type="number" min={0} value={tempAmount}
                                onChange={e => setTempAmount(e.target.value)} placeholder="Temp HP" />
                            <button className="bg-info/15 border border-info/30 text-info px-3 py-1.5 rounded-sm text-xs hover:bg-info/25" onClick={() => { onTempHP(tempAmount === '' ? 0 : +tempAmount); setTempAmount(''); }}>🛡 Set Temp HP</button>
                        </div>

                        <div>
                            <div className="text-[10px] text-ash uppercase tracking-wide mb-1.5">Conditions</div>
                            <div className="flex flex-wrap gap-1.5">
                                {CONDITIONS.map(c => (
                                    <button key={c}
                                        className={`bg-stone border border-stone-border text-ash px-2.5 py-[5px] rounded-full text-xs transition-all duration-180 hover:border-crimson hover:text-parchment ${p.conditions?.includes(c) ? 'bg-crimson/20 border-crimson text-crimson-light' : ''}`}
                                        onClick={() => toggleCondition(c)}>
                                        {c}
                                    </button>
                                ))}
                            </div>
                        </div>

                        <div className="flex items-center gap-2 flex-wrap">
                            <button
                                className={`flex items-center gap-1.5 border px-3 py-1.5 rounded-sm text-xs transition-all duration-180 ${p.concentration ? 'bg-[#9b59b6]/20 border-[#9b59b6]/40 text-[#c99ee0]' : 'bg-stone border-stone-border text-ash hover:text-parchment'}`}
                                onClick={onToggleConcentration}>
                                ⊛ Concentration
                            </button>
                            {p.is_active && (
                                <button className="flex items-center gap-1.5 bg-stone border border-stone-border text-ash px-3 py-1.5 rounded-sm text-xs transition-all duration-180 hover:border-crimson hover:text-crimson-light" onClick={onDeactivate}>
                                    ✕ Mark Inactive
                                </button>
                            )}
                            <button className="flex items-center gap-1 bg-transparent border-none text-ash text-xs ml-auto hover:text-crimson-light" onClick={onRemove}><IconTrash size={13} /> Remove</button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
