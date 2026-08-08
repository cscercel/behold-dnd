import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { listCharacters, createCharacter, deleteCharacter } from '../api';
import { IconPlus, IconSword, IconShield, IconTrash, IconChevronRight } from '../components/Icon';
import styles from './CharacterList.module.css';

const CLASS_COLORS: Record<string, string> = {
  Barbarian: '#c0392b', Bard: '#8e44ad', Cleric: '#f39c12', Druid: '#27ae60',
  Fighter: '#7f8c8d', Monk: '#16a085', Paladin: '#f1c40f', Ranger: '#2ecc71',
  Rogue: '#2c3e50', Sorcerer: '#e74c3c', Warlock: '#6c3483', Wizard: '#2980b9',
};

export function CharacterList() {
  const [chars, setChars] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const navigate = useNavigate();

  useEffect(() => { listCharacters().then(setChars).finally(() => setLoading(false)); }, []);

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm('Delete this character permanently?')) return;
    await deleteCharacter(id);
    setChars(c => c.filter(ch => ch.id !== id));
  };

  if (loading) return <div className={styles.loading}>Loading characters…</div>;

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <div>
          <h1 className={styles.title}>Your Characters</h1>
          <p className={styles.subtitle}>{chars.length} adventurer{chars.length !== 1 ? 's' : ''} in your party</p>
        </div>
        <button className={styles.createBtn} onClick={() => setShowCreate(true)}>
          <IconPlus size={16} /> New Character
        </button>
      </div>

      {chars.length === 0 ? (
        <div className={styles.empty}>
          <IconSword size={48} className={styles.emptyIcon} />
          <p>No characters yet. Create your first adventurer.</p>
        </div>
      ) : (
        <div className={styles.grid}>
          {chars.map(char => (
            <div key={char.id} className={styles.card} onClick={() => navigate(`/characters/${char.id}`)}>
              <div className={styles.cardAccent} style={{ background: CLASS_COLORS[char.class] || '#8b1a1a' }} />
              <div className={styles.cardBody}>
                <div className={styles.cardTop}>
                  <div>
                    <h2 className={styles.charName}>{char.name}</h2>
                    <p className={styles.charMeta}>Level {char.level} {char.race} {char.class}</p>
                  </div>
                  <button className={styles.deleteBtn} onClick={e => handleDelete(char.id, e)}>
                    <IconTrash size={14} />
                  </button>
                </div>
                <div className={styles.cardStats}>
                  <span className={styles.stat}><IconShield size={12} /> {char.armor_class} AC</span>
                  <span className={styles.stat}>❤ {char.current_hp}/{char.max_hp}</span>
                  {char.is_npc && <span className={styles.npcTag}>NPC</span>}
                </div>
              </div>
              <IconChevronRight size={16} className={styles.chevron} />
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
    </div>
  );
}

function CreateCharacterModal({ onClose, onCreated }: { onClose: () => void; onCreated: (c: any) => void }) {
  const RACES = ['Human','Elf','Dwarf','Halfling','Gnome','Half-Elf','Half-Orc','Tiefling','Dragonborn','Aasimar','Tabaxi','Kenku','Firbolg'];
  const CLASSES = ['Barbarian','Bard','Cleric','Druid','Fighter','Monk','Paladin','Ranger','Rogue','Sorcerer','Warlock','Wizard','Artificer'];
  const ALIGNMENTS = ['Lawful Good','Neutral Good','Chaotic Good','Lawful Neutral','True Neutral','Chaotic Neutral','Lawful Evil','Neutral Evil','Chaotic Evil'];

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
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={e => e.stopPropagation()}>
        <h2 className={styles.modalTitle}>Create Character</h2>
        <form onSubmit={submit} className={styles.modalForm}>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Name *</label><input value={form.name} onChange={set('name')} required placeholder="Character name" /></div>
            <div className={styles.field}><label>Level</label><input type="number" min={1} max={20} value={form.level} onChange={set('level')} /></div>
          </div>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Race</label><select value={form.race} onChange={set('race')}>{RACES.map(r => <option key={r}>{r}</option>)}</select></div>
            <div className={styles.field}><label>Class</label><select value={form.class} onChange={set('class')}>{CLASSES.map(c => <option key={c}>{c}</option>)}</select></div>
          </div>
          <div className={styles.formRow}>
            <div className={styles.field}><label>Background</label><input value={form.background} onChange={set('background')} placeholder="Soldier, Sage…" /></div>
            <div className={styles.field}><label>Alignment</label><select value={form.alignment} onChange={set('alignment')}>{ALIGNMENTS.map(a => <option key={a}>{a}</option>)}</select></div>
          </div>
          <div className={styles.sectionLabel}>Ability Scores</div>
          <div className={styles.abilityGrid}>
            {(['strength','dexterity','constitution','intelligence','wisdom','charisma'] as const).map(ab => (
              <div key={ab} className={styles.abilityField}>
                <label>{ab.slice(0,3).toUpperCase()}</label>
                <input type="number" min={1} max={30} value={form[ab]} onChange={set(ab)} />
              </div>
            ))}
          </div>
          <div className={styles.sectionLabel}>Combat Stats</div>
          <div className={styles.formRow3}>
            <div className={styles.field}><label>Max HP</label><input type="number" min={1} value={form.max_hp} onChange={set('max_hp')} /></div>
            <div className={styles.field}><label>AC</label><input type="number" min={1} value={form.armor_class} onChange={set('armor_class')} /></div>
            <div className={styles.field}><label>Speed</label><input type="number" min={0} value={form.speed} onChange={set('speed')} /></div>
            <div className={styles.field}><label>Hit Die</label><input type="number" min={4} value={form.hit_dice_type} onChange={set('hit_dice_type')} /></div>
          </div>
          {error && <p className={styles.error}>{error}</p>}
          <div className={styles.modalActions}>
            <button type="button" className={styles.cancelBtn} onClick={onClose}>Cancel</button>
            <button type="submit" className={styles.submitBtn}>Create Character</button>
          </div>
        </form>
      </div>
    </div>
  );
}
