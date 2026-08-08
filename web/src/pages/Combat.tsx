import { useState, useEffect } from 'react';
import * as api from '../api';
import { IconPlus, IconSword, IconSkip, IconStop, IconPlay, IconTrash, IconChevronDown, IconChevronUp } from '../components/Icon';
import styles from './Combat.module.css';

export function Combat() {
  const [encounters, setEncounters] = useState<any[]>([]);
  const [active, setActive] = useState<any>(null);
  const [participants, setParticipants] = useState<any[]>([]);
  const [characters, setCharacters] = useState<any[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [newName, setNewName] = useState('');
  const [loading, setLoading] = useState(true);
  const [showAddPart, setShowAddPart] = useState(false);
  const [partForm, setPartForm] = useState({ character_id: '', name: '', initiative: 0, max_hp: 10, current_hp: 10, armor_class: 10, speed: 30 });

  const reload = async () => {
    const [enc, chars] = await Promise.all([api.listEncounters(), api.listCharacters()]);
    setEncounters(enc || []); setCharacters(chars || []);
    try {
      const act = await api.getActiveEncounter();
      setActive(act);
      if (act) {
        const parts = await api.listParticipants(act.id);
        setParticipants((parts || []).sort((a: any, b: any) => b.initiative - a.initiative));
      }
    } catch { setActive(null); setParticipants([]); }
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
    await api.addParticipant(encId, partForm);
    setShowAddPart(false);
    setPartForm({ character_id: '', name: '', initiative: 0, max_hp: 10, current_hp: 10, armor_class: 10, speed: 30 });
    reload();
  };

  const selectChar = (charId: string) => {
    const c = characters.find((ch: any) => ch.id === charId);
    if (c) setPartForm(f => ({ ...f, character_id: charId, name: c.name, max_hp: c.max_hp, current_hp: c.current_hp, armor_class: c.armor_class, speed: c.speed }));
    else   setPartForm(f => ({ ...f, character_id: charId }));
  };

  const encId = active?.id || selectedId;
  const currentEnc = active || encounters.find((e: any) => e.id === selectedId);

  if (loading) return <div className={styles.loading}>Loading…</div>;

  return (
    <div className={styles.page}>
      <div className={styles.sidebar}>
        <h2 className={styles.sideTitle}>Encounters</h2>
        <div className={styles.newEncounter}>
          <input value={newName} onChange={e => setNewName(e.target.value)}
            placeholder="Encounter name…" onKeyDown={e => e.key === 'Enter' && createEnc()} />
          <button className={styles.createBtn} onClick={createEnc}><IconPlus size={14} /></button>
        </div>
        <div className={styles.encList}>
          {encounters.map((enc: any) => (
            <div key={enc.id}
              className={`${styles.encItem} ${enc.id === encId ? styles.encActive : ''} ${enc.is_active ? styles.encLive : ''}`}
              onClick={() => loadEnc(enc.id)}>
              <div>
                <div className={styles.encName}>{enc.name}</div>
                <div className={styles.encMeta}>{enc.is_active ? `Round ${enc.round} · LIVE` : 'Inactive'}</div>
              </div>
              {enc.is_active && <div className={styles.liveDot} />}
            </div>
          ))}
          {encounters.length === 0 && <div className={styles.emptyEnc}>No encounters yet.</div>}
        </div>
      </div>

      <div className={styles.main}>
        {currentEnc ? (
          <>
            <div className={styles.encHeader}>
              <div>
                <h1 className={styles.encTitle}>{currentEnc.name}</h1>
                {currentEnc.is_active && <span className={styles.roundBadge}>Round {currentEnc.round}</span>}
              </div>
              <div className={styles.encActions}>
                {!currentEnc.is_active
                  ? <button className={styles.actionBtn} onClick={() => api.startEncounter(currentEnc.id).then(reload)}><IconPlay size={15} /> Start</button>
                  : <>
                      <button className={styles.actionBtn} onClick={() => api.nextRound(currentEnc.id).then(reload)}><IconSkip size={15} /> Next Round</button>
                      <button className={styles.actionBtnDanger} onClick={() => api.endEncounter(currentEnc.id).then(reload)}><IconStop size={15} /> End</button>
                    </>
                }
                <button className={styles.addPartBtn} onClick={() => setShowAddPart(v => !v)}><IconPlus size={15} /> Add Combatant</button>
              </div>
            </div>

            {showAddPart && (
              <form onSubmit={addPart} className={styles.addForm}>
                <div className={styles.formRow}>
                  <div className={styles.field}>
                    <label>From Characters</label>
                    <select value={partForm.character_id} onChange={e => selectChar(e.target.value)}>
                      <option value="">— Custom NPC —</option>
                      {characters.map((c: any) => <option key={c.id} value={c.id}>{c.name}</option>)}
                    </select>
                  </div>
                  <div className={styles.field}><label>Name *</label><input value={partForm.name} onChange={e => setPartForm(f => ({ ...f, name: e.target.value }))} required /></div>
                  <div className={styles.field}><label>Initiative</label><input type="number" value={partForm.initiative} onChange={e => setPartForm(f => ({ ...f, initiative: +e.target.value }))} /></div>
                  <div className={styles.field}><label>Max HP</label><input type="number" min={1} value={partForm.max_hp} onChange={e => setPartForm(f => ({ ...f, max_hp: +e.target.value, current_hp: +e.target.value }))} /></div>
                  <div className={styles.field}><label>AC</label><input type="number" min={1} value={partForm.armor_class} onChange={e => setPartForm(f => ({ ...f, armor_class: +e.target.value }))} /></div>
                </div>
                <div className={styles.modalActions}>
                  <button type="button" className={styles.cancelBtn} onClick={() => setShowAddPart(false)}>Cancel</button>
                  <button type="submit" className={styles.submitBtn}>Add to Battle</button>
                </div>
              </form>
            )}

            <div className={styles.tracker}>
              {participants.length === 0 && (
                <div className={styles.emptyTracker}><IconSword size={40} className={styles.emptyIcon} /><p>No combatants yet.</p></div>
              )}
              {participants.map((p: any) => (
                <ParticipantRow key={p.id} p={p} encId={encId}
                  onDmg={n => api.participantDamage(encId, p.id, n).then(reload)}
                  onHeal={n => api.participantHeal(encId, p.id, n).then(reload)}
                  onRemove={() => api.removeParticipant(encId, p.id).then(reload)}
                />
              ))}
            </div>
          </>
        ) : (
          <div className={styles.noEnc}>
            <IconSword size={64} className={styles.noEncIcon} />
            <p>Select or create an encounter to begin</p>
          </div>
        )}
      </div>
    </div>
  );
}

function ParticipantRow({ p, onDmg, onHeal, onRemove }: {
  p: any; encId: string;
  onDmg: (n: number) => void;
  onHeal: (n: number) => void;
  onRemove: () => void;
}) {
  const [expanded, setExpanded] = useState(false);
  const [amount, setAmount] = useState('');
  const hpPct   = Math.max(0, Math.min(100, (p.current_hp / p.max_hp) * 100));
  const hpColor = hpPct > 50 ? '#27ae60' : hpPct > 25 ? '#f39c12' : '#c0392b';

  return (
    <div className={`${styles.participantCard} ${!p.is_active ? styles.partInactive : ''}`}>
      <div className={styles.partInit}>{p.initiative}</div>
      <div className={styles.partMain}>
        <div className={styles.partTop} onClick={() => setExpanded(v => !v)}>
          <div>
            <div className={styles.partName}>{p.name}</div>
            {p.conditions?.length > 0 && (
              <div className={styles.partConditions}>
                {p.conditions.map((c: string) => <span key={c} className={styles.cond}>{c}</span>)}
              </div>
            )}
          </div>
          <div className={styles.partHpArea}>
            <div className={styles.partHpBar}>
              <div style={{ width: `${hpPct}%`, background: hpColor, height: '100%', borderRadius: 3 }} />
            </div>
            <span className={styles.partHpNum} style={{ color: hpColor }}>{p.current_hp}</span>
            <span className={styles.partHpMax}>/{p.max_hp}</span>
          </div>
          <div className={styles.partStats}>
            <span>AC {p.armor_class}</span>
            {p.concentration && <span className={styles.concentrating}>⊛</span>}
          </div>
          {expanded ? <IconChevronUp size={14} className={styles.partChevron} /> : <IconChevronDown size={14} className={styles.partChevron} />}
        </div>
        {expanded && (
          <div className={styles.partExpanded}>
            <input className={styles.partAmountInput} type="number" min={1} value={amount}
              onChange={e => setAmount(e.target.value)} placeholder="Amount" />
            <button className={styles.dmgBtn} onClick={() => { if (amount) { onDmg(+amount); setAmount(''); } }}>⚔ Damage</button>
            <button className={styles.healBtn} onClick={() => { if (amount) { onHeal(+amount); setAmount(''); } }}>❤ Heal</button>
            <button className={styles.removeBtn} onClick={onRemove}><IconTrash size={13} /></button>
          </div>
        )}
      </div>
    </div>
  );
}
