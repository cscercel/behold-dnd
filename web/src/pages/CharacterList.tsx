import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { listCharacters, createCharacter, deleteCharacter } from '../lib/api';
import { IconPlus, IconSword, IconShield, IconTrash, IconChevronRight } from '../components/Icon';

const CLASS_COLORS: Record<string, string> = {
    Barbarian: '#c0392b', Bard: '#8e44ad', Cleric: '#f39c12', Druid: '#27ae60',
    Fighter: '#7f8c8d', Monk: '#16a085', Paladin: '#f1c40f', Ranger: '#2ecc71',
    Rogue: '#2c3e50', Sorcerer: '#e74c3c', Warlock: '#6c3483', Wizard: '#2980b9',
};

const fieldLabel = "text-[11px] font-medium text-ash-light uppercase tracking-wide";
const field = "flex flex-col gap-1.5";

export function CharacterList() {
    const [chars, setChars] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreate, setShowCreate] = useState(false);
    const [deleteTarget, setDeleteTarget] = useState<any>(null);
    const navigate = useNavigate();

    useEffect(() => { listCharacters().then(setChars).finally(() => setLoading(false)); }, []);

    const confirmDelete = async () => {
        if (!deleteTarget) return;
        await deleteCharacter(deleteTarget.id);
        setChars(c => c.filter(ch => ch.id !== deleteTarget.id));
        setDeleteTarget(null);
    };

    if (loading) return <div className="flex items-center justify-center h-[200px] text-ash">Loading characters…</div>;

    return (
        <div className="p-8 max-w-[1100px] mx-auto">
            <div className="flex items-start justify-between mb-8">
                <div>
                    <h1 className="font-display text-[28px] font-bold text-parchment">Your Characters</h1>
                    <p className="text-sm text-ash mt-1">{chars.length} adventurer{chars.length !== 1 ? 's' : ''} in your party</p>
                </div>
                <button className="flex items-center gap-2 bg-crimson border-none text-parchment px-[18px] py-2.5 rounded-md font-display text-[13px] font-semibold tracking-wide transition-colors duration-180 hover:bg-crimson-light" onClick={() => setShowCreate(true)}>
                    <IconPlus size={16} /> New Character
                </button>
            </div>

            {chars.length === 0 ? (
                <div className="text-center py-20 px-5 text-ash">
                    <IconSword size={48} className="mx-auto mb-4 opacity-30 block" />
                    <p>No characters yet. Create your first adventurer.</p>
                </div>
            ) : (
                <div className="grid grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-4">
                    {chars.map(char => (
                        <div key={char.id} className="flex items-stretch bg-stone border border-stone-border rounded-lg cursor-pointer overflow-hidden transition-[border-color,transform] duration-180 hover:border-gold-dim hover:-translate-y-0.5" onClick={() => navigate(`/characters/${char.id}`)}>
                            <div className="w-1 shrink-0" style={{ background: CLASS_COLORS[char.class] || '#8b1a1a' }} />
                            <div className="flex-1 p-4">
                                <div className="flex items-start justify-between mb-3">
                                    <div>
                                        <h2 className="font-display text-[17px] font-semibold text-parchment">{char.name}</h2>
                                        <p className="text-xs text-ash mt-0.5">Level {char.level} {char.race} {char.class}</p>
                                    </div>
                                    <button className="bg-transparent border-none text-stone-border p-1 rounded-sm transition-colors duration-180 shrink-0 hover:text-crimson-light" onClick={e => { e.stopPropagation(); setDeleteTarget(char); }}>
                                        <IconTrash size={14} />
                                    </button>
                                </div>
                                <div className="flex gap-3 items-center flex-wrap">
                                    <span className="flex items-center gap-1 text-xs text-ash-light"><IconShield size={12} /> {char.armor_class} AC</span>
                                    <span className="flex items-center gap-1 text-xs text-ash-light">❤ {char.current_hp}/{char.max_hp}</span>
                                    {char.is_npc && <span className="bg-gold/15 text-gold border border-gold/30 px-1.5 py-0.5 rounded text-[10px] font-semibold tracking-widest">NPC</span>}
                                </div>
                            </div>
                            <IconChevronRight size={16} className="text-stone-border self-center mr-3 shrink-0" />
                        </div>
                    ))}
                </div>
            )}

            {showCreate && (
                <CreateCharacterModal
                    onClose={() => setShowCreate(false)}
                    onCreated={c => { setChars(ch => [...ch, c]); setShowCreate(false); navigate(`/characters/${c.id}`); }}
                />
            )}

            {deleteTarget && (
                <div className="fixed inset-0 bg-black/75 z-[100] flex items-center justify-center p-5" onClick={() => setDeleteTarget(null)}>
                    <div className="bg-stone border border-stone-border rounded-lg p-7 w-full max-w-[400px]" onClick={e => e.stopPropagation()}>
                        <h2 className="font-display text-lg font-bold text-parchment mb-2">Delete Character?</h2>
                        <p className="text-sm text-ash-light leading-relaxed mb-6">
                            Are you sure you want to permanently delete <span className="text-parchment font-medium">{deleteTarget.name}</span>? This can't be undone.
                        </p>
                        <div className="flex justify-end gap-2.5">
                            <button className="bg-stone-mid border-none text-ash-light px-5 py-2.5 rounded-md text-[13px] hover:text-parchment" onClick={() => setDeleteTarget(null)}>Cancel</button>
                            <button className="bg-crimson border-none text-parchment px-5 py-2.5 rounded-md font-display text-[13px] font-semibold tracking-wide hover:bg-crimson-light" onClick={confirmDelete}>Delete</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function CreateCharacterModal({ onClose, onCreated }: { onClose: () => void; onCreated: (c: any) => void }) {
    const RACES = ['Human', 'Elf', 'Dwarf', 'Halfling', 'Gnome', 'Half-Elf', 'Half-Orc', 'Tiefling', 'Dragonborn', 'Aasimar', 'Tabaxi', 'Kenku', 'Firbolg'];
    const CLASSES = ['Barbarian', 'Bard', 'Cleric', 'Druid', 'Fighter', 'Monk', 'Paladin', 'Ranger', 'Rogue', 'Sorcerer', 'Warlock', 'Wizard', 'Artificer'];
    const ALIGNMENTS = ['Lawful Good', 'Neutral Good', 'Chaotic Good', 'Lawful Neutral', 'True Neutral', 'Chaotic Neutral', 'Lawful Evil', 'Neutral Evil', 'Chaotic Evil'];

    const [form, setForm] = useState({
        name: '', race: 'Human', class: 'Fighter', level: 1, background: '', alignment: 'True Neutral',
        strength: 10, dexterity: 10, constitution: 10, intelligence: 10, wisdom: 10, charisma: 10,
        max_hp: 10, current_hp: 10, armor_class: 10, speed: 30, hit_dice_type: 8, hit_dice_remaining: 1,
        spellcasting_ability: 'none', is_npc: false,
    });
    const [error, setError] = useState('');

    const set = (k: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
        setForm(f => ({ ...f, [k]: e.target.type === 'number' ? Number(e.target.value) : e.target.type === 'checkbox' ? (e.target as HTMLInputElement).checked : e.target.value }));

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        try { onCreated(await createCharacter(form)); }
        catch (err: any) { setError(err?.data?.error || 'Failed to create character'); }
    };

    return (
        <div className="fixed inset-0 bg-black/75 z-[100] flex items-center justify-center p-5" onClick={onClose}>
            <div className="bg-stone border border-stone-border rounded-lg p-8 w-full max-w-[560px] max-h-[90vh] overflow-y-auto" onClick={e => e.stopPropagation()}>
                <h2 className="font-display text-xl font-bold text-parchment mb-6">Create Character</h2>
                <form onSubmit={submit} className="flex flex-col gap-4">
                    <div className="grid grid-cols-2 gap-3">
                        <div className={field}><label className={fieldLabel}>Name *</label><input value={form.name} onChange={set('name')} required placeholder="Character name" /></div>
                        <div className={field}><label className={fieldLabel}>Level</label><input type="number" min={1} max={20} value={form.level} onChange={set('level')} /></div>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                        <div className={field}><label className={fieldLabel}>Race</label><select value={form.race} onChange={set('race')}>{RACES.map(r => <option key={r}>{r}</option>)}</select></div>
                        <div className={field}><label className={fieldLabel}>Class</label><select value={form.class} onChange={set('class')}>{CLASSES.map(c => <option key={c}>{c}</option>)}</select></div>
                    </div>
                    <div className="grid grid-cols-2 gap-3">
                        <div className={field}><label className={fieldLabel}>Background</label><input value={form.background} onChange={set('background')} placeholder="Soldier, Sage…" /></div>
                        <div className={field}><label className={fieldLabel}>Alignment</label><select value={form.alignment} onChange={set('alignment')}>{ALIGNMENTS.map(a => <option key={a}>{a}</option>)}</select></div>
                    </div>
                    <div className="text-[11px] font-semibold text-gold-dim uppercase tracking-widest border-b border-stone-border pb-1.5">Ability Scores</div>
                    <div className="grid grid-cols-6 gap-2">
                        {(['strength', 'dexterity', 'constitution', 'intelligence', 'wisdom', 'charisma'] as const).map(ab => (
                            <div key={ab} className="flex flex-col gap-1">
                                <label className="text-[10px] font-bold text-ash text-center tracking-wider">{ab.slice(0, 3).toUpperCase()}</label>
                                <input className="text-center" type="number" min={1} max={30} value={form[ab]} onChange={set(ab)} />
                            </div>
                        ))}
                    </div>
                    <div className="text-[11px] font-semibold text-gold-dim uppercase tracking-widest border-b border-stone-border pb-1.5">Combat Stats</div>
                    <div className="grid grid-cols-4 gap-3">
                        <div className={field}><label className={fieldLabel}>Max HP</label><input type="number" min={1} value={form.max_hp} onChange={set('max_hp')} /></div>
                        <div className={field}><label className={fieldLabel}>AC</label><input type="number" min={1} value={form.armor_class} onChange={set('armor_class')} /></div>
                        <div className={field}><label className={fieldLabel}>Speed</label><input type="number" min={0} value={form.speed} onChange={set('speed')} /></div>
                        <div className={field}><label className={fieldLabel}>Hit Die</label><input type="number" min={4} value={form.hit_dice_type} onChange={set('hit_dice_type')} /></div>
                    </div>
                    {error && <p className="text-[13px] text-crimson-light">{error}</p>}
                    <div className="flex justify-end gap-2.5 mt-2">
                        <button type="button" className="bg-stone-mid border-none text-ash-light px-5 py-2.5 rounded-md text-[13px] hover:text-parchment" onClick={onClose}>Cancel</button>
                        <button type="submit" className="bg-crimson border-none text-parchment px-5 py-2.5 rounded-md font-display text-[13px] font-semibold tracking-wide hover:bg-crimson-light">Create Character</button>
                    </div>
                </form>
            </div>
        </div>
    );
}
