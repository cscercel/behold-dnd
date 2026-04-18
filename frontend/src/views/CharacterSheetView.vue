<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCharactersStore } from '@/stores/characters'
import { useCharacter } from '@/composables/useCharacter'
import { useAuthStore } from '@/stores/auth'
import { inventoryAPI } from '@/api/inventory'
import { spellsAPI } from '@/api/spells'
import { featuresAPI } from '@/api/features'
import type { InventoryItem, Spell, Feature } from '@/types'

const route  = useRoute()
const router = useRouter()
const store  = useCharactersStore()
const auth   = useAuthStore()
const id     = route.params.id as string

onMounted(async () => {
  await store.fetchOne(id)
  if (store.current) await Promise.all([loadInventory(), loadSpells(), loadFeatures()])
})

const derived = computed(() => store.current ? useCharacter(store.current) : null)

const hpAction = ref<'damage' | 'heal' | 'temp' | null>(null)
const hpAmount = ref(0)
const actionLoading = ref(false)

async function submitHPAction() {
  if (!hpAmount.value || hpAmount.value <= 0) return
  actionLoading.value = true
  try {
    if (hpAction.value === 'damage') await store.applyDamage(id, hpAmount.value)
    if (hpAction.value === 'heal')   await store.heal(id, hpAmount.value)
    if (hpAction.value === 'temp')   await store.addTempHP(id, hpAmount.value)
    hpAction.value = null; hpAmount.value = 0
  } finally { actionLoading.value = false }
}

async function handleLongRest() {
  if (confirm('Take a long rest? Restores HP, hit dice, spell slots and clears conditions.'))
    await store.longRest(id)
}

async function handleShortRest() {
  const diceStr = prompt(`How many hit dice to spend? (${store.current?.hit_dice_remaining}d${store.current?.hit_dice_type} available)`)
  if (!diceStr) return
  const dice = parseInt(diceStr); if (isNaN(dice) || dice <= 0) return
  const hpStr = prompt('How much HP regained?')
  if (!hpStr) return
  const hp = parseInt(hpStr); if (isNaN(hp) || hp <= 0) return
  await store.shortRest(id, dice, hp)
}

function hpColor(pct: number) {
  return pct > 50 ? 'var(--hp-high)' : pct > 25 ? 'var(--hp-mid)' : 'var(--hp-low)'
}

const editingAbilities = ref(false)
const editingProfs = ref(false)
const abilityDraft = ref<Record<string,number>>({})
const profDraft = ref<Record<string,boolean|number>>({})

function startEditAbilities() {
  if (!store.current) return
  const c = store.current
  abilityDraft.value = { strength: c.strength, dexterity: c.dexterity, constitution: c.constitution, intelligence: c.intelligence, wisdom: c.wisdom, charisma: c.charisma }
  editingAbilities.value = true
}
async function saveAbilities() { await store.update(id, abilityDraft.value); editingAbilities.value = false }

function startEditProfs() {
  if (!store.current) return
  const c = store.current
  profDraft.value = {
    save_prof_strength: c.save_prof_strength, save_prof_dexterity: c.save_prof_dexterity,
    save_prof_constitution: c.save_prof_constitution, save_prof_intelligence: c.save_prof_intelligence,
    save_prof_wisdom: c.save_prof_wisdom, save_prof_charisma: c.save_prof_charisma,
    skill_acrobatics: c.skill_acrobatics, skill_animal_handling: c.skill_animal_handling,
    skill_arcana: c.skill_arcana, skill_athletics: c.skill_athletics, skill_deception: c.skill_deception,
    skill_history: c.skill_history, skill_insight: c.skill_insight, skill_intimidation: c.skill_intimidation,
    skill_investigation: c.skill_investigation, skill_medicine: c.skill_medicine, skill_nature: c.skill_nature,
    skill_perception: c.skill_perception, skill_performance: c.skill_performance, skill_persuasion: c.skill_persuasion,
    skill_religion: c.skill_religion, skill_sleight_of_hand: c.skill_sleight_of_hand,
    skill_stealth: c.skill_stealth, skill_survival: c.skill_survival,
  }
  editingProfs.value = true
}
async function saveProfs() { await store.update(id, profDraft.value as any); editingProfs.value = false }
function cycleProficiency(field: string) { profDraft.value[field] = ((profDraft.value[field] as number) + 1) % 3 }
function skillField(name: string) { return 'skill_' + name.toLowerCase().replace(/ /g,'_') }
function profLabel(level: number) { return level === 2 ? 'E' : level === 1 ? '●' : '○' }

type Tab = 'actions'|'spells'|'feats'|'inventory'|'proficiencies'|'personality'
const activeTab = ref<Tab>('actions')
const TABS: {key:Tab;label:string}[] = [
  {key:'actions',label:'Actions'},{key:'spells',label:'Spells'},{key:'feats',label:'Feats'},
  {key:'inventory',label:'Inventory'},{key:'proficiencies',label:'Proficiencies'},{key:'personality',label:'Personality'}
]

const inventory = ref<InventoryItem[]>([])
const inventoryLoading = ref(false)
const showItemForm = ref(false)
const editingItem = ref<InventoryItem|null>(null)
const itemDraft = ref({name:'',quantity:1,weight:0,value:0,description:'',is_equipped:false,requires_attunement:false})

async function loadInventory() {
  inventoryLoading.value = true
  try { inventory.value = await inventoryAPI.list(id) } finally { inventoryLoading.value = false }
}
function startAddItem() {
  editingItem.value = null
  itemDraft.value = {name:'',quantity:1,weight:0,value:0,description:'',is_equipped:false,requires_attunement:false}
  showItemForm.value = true
}
function startEditItem(item: InventoryItem) {
  editingItem.value = item
  itemDraft.value = {name:item.name,quantity:item.quantity,weight:item.weight,value:(item as any).value??0,description:item.description,is_equipped:item.is_equipped,requires_attunement:item.requires_attunement}
  showItemForm.value = true
}
async function saveItem() {
  if (editingItem.value) {
    const u = await inventoryAPI.update(id, editingItem.value.id, itemDraft.value)
    const i = inventory.value.findIndex(x=>x.id===editingItem.value!.id); if(i!==-1) inventory.value[i]=u
  } else { inventory.value.push(await inventoryAPI.create(id, itemDraft.value)) }
  showItemForm.value = false
}
async function deleteItem(itemId: string) {
  if (!confirm('Remove this item?')) return
  await inventoryAPI.delete(id, itemId); inventory.value = inventory.value.filter(i=>i.id!==itemId)
}
const totalWeight = computed(() => inventory.value.reduce((s,i)=>s+i.weight*i.quantity,0))

const spells = ref<Spell[]>([])
const spellsLoading = ref(false)
const showSpellForm = ref(false)
const editingSpell = ref<Spell|null>(null)
const spellDraft = ref({name:'',level:0,school:'',casting_time:'',range:'',components:'',duration:'',description:'',is_prepared:false})

async function loadSpells() {
  spellsLoading.value = true
  try { spells.value = await spellsAPI.list(id) } finally { spellsLoading.value = false }
}
function startAddSpell() {
  editingSpell.value = null
  spellDraft.value = {name:'',level:0,school:'',casting_time:'',range:'',components:'',duration:'',description:'',is_prepared:false}
  showSpellForm.value = true
}
function startEditSpell(spell: Spell) {
  editingSpell.value = spell
  spellDraft.value = {name:spell.name,level:spell.level,school:spell.school,casting_time:spell.casting_time,range:spell.range,components:spell.components,duration:spell.duration,description:spell.description,is_prepared:spell.is_prepared}
  showSpellForm.value = true
}
async function saveSpell() {
  if (editingSpell.value) {
    const u = await spellsAPI.update(id, editingSpell.value.id, spellDraft.value)
    const i = spells.value.findIndex(s=>s.id===editingSpell.value!.id); if(i!==-1) spells.value[i]=u
  } else { spells.value.push(await spellsAPI.create(id, spellDraft.value)) }
  showSpellForm.value = false
}
async function deleteSpell(spellId: string) {
  if (!confirm('Remove spell?')) return
  await spellsAPI.delete(id, spellId); spells.value = spells.value.filter(s=>s.id!==spellId)
}
async function togglePrepared(spell: Spell) {
  const u = await spellsAPI.togglePrepared(id, spell.id)
  const i = spells.value.findIndex(s=>s.id===spell.id); if(i!==-1) spells.value[i]=u
}
const spellsByLevel = computed(() => {
  const g: Record<number,Spell[]> = {}
  for (const s of spells.value) { if(!g[s.level]) g[s.level]=[]; g[s.level].push(s) }
  return Object.entries(g).sort(([a],[b])=>Number(a)-Number(b))
})

const features = ref<Feature[]>([])
const featuresLoading = ref(false)
const showFeatureForm = ref(false)
const editingFeature = ref<Feature | null>(null)
const featureDraft = ref({ name: '', source: '', description: '', action_type: 'none' as ActionType })

async function loadFeatures() {
  featuresLoading.value = true
  try { features.value = await featuresAPI.list(id) }
  finally { featuresLoading.value = false }
}

function startAddFeature(defaultType: ActionType = 'none') {
  editingFeature.value = null
  featureDraft.value = { name: '', source: '', description: '', action_type: defaultType }
  showFeatureForm.value = true
}

function startEditFeature(feature: Feature) {
  editingFeature.value = feature
  featureDraft.value = { name: feature.name, source: feature.source, description: feature.description, action_type: feature.action_type }
  showFeatureForm.value = true
}

async function saveFeature() {
  if (editingFeature.value) {
    const updated = await featuresAPI.update(id, editingFeature.value.id, featureDraft.value)
    const idx = features.value.findIndex(f => f.id === editingFeature.value!.id)
    if (idx !== -1) features.value[idx] = updated
  } else {
    features.value.push(await featuresAPI.create(id, featureDraft.value))
  }
  showFeatureForm.value = false
}

async function deleteFeature(featureId: string) {
  if (!confirm('Remove this feature?')) return
  await featuresAPI.delete(id, featureId)
  features.value = features.value.filter(f => f.id !== featureId)
}

// Computed feature groups for the tabs
const actions      = computed(() => features.value.filter(f => f.action_type === 'action'))
const bonusActions = computed(() => features.value.filter(f => f.action_type === 'bonus_action'))
const reactions    = computed(() => features.value.filter(f => f.action_type === 'reaction'))
const freeActions  = computed(() => features.value.filter(f => f.action_type === 'free'))
const feats        = computed(() => features.value.filter(f => f.action_type === 'none'))

const editingPersonality = ref(false)
const personalityDraft = ref({personality_traits:'',ideals:'',bonds:'',flaws:'',notes:''})
function startEditPersonality() {
  if (!store.current) return
  const c = store.current
  personalityDraft.value = {personality_traits:c.personality_traits,ideals:c.ideals,bonds:c.bonds,flaws:c.flaws,notes:c.notes}
  editingPersonality.value = true
}
async function savePersonality() { await store.update(id, personalityDraft.value); editingPersonality.value = false }

async function handleLevelChange() {
  const current = store.current!.level
  const input = prompt(`Current level: ${current}\nEnter +1 to level up or -1 to level down (min 1):`)
  if (!input) return
  const delta = parseInt(input)
  if (isNaN(delta)) return
  const newLevel = Math.max(1, current + delta)
  if (newLevel === current) return
  await store.update(id, { level: newLevel })
}

async function handleDeathSave(success: boolean) {
  await store.recordDeathSave(id, success)
}

function toggleComponent(c: string) {
  const parts = spellDraft.value.components
    .split(',').map(s => s.trim()).filter(Boolean)
  const idx = parts.indexOf(c)
  if (idx === -1) parts.push(c)
  else parts.splice(idx, 1)
  spellDraft.value.components = parts.join(', ')
}

const ABILITIES = [
  {key:'strength',label:'STR'},{key:'dexterity',label:'DEX'},{key:'constitution',label:'CON'},
  {key:'intelligence',label:'INT'},{key:'wisdom',label:'WIS'},{key:'charisma',label:'CHA'}
] as const

const SAVES = [
  {key:'strength',label:'Strength',prof:'save_prof_strength'},
  {key:'dexterity',label:'Dexterity',prof:'save_prof_dexterity'},
  {key:'constitution',label:'Constitution',prof:'save_prof_constitution'},
  {key:'intelligence',label:'Intelligence',prof:'save_prof_intelligence'},
  {key:'wisdom',label:'Wisdom',prof:'save_prof_wisdom'},
  {key:'charisma',label:'Charisma',prof:'save_prof_charisma'},
] as const
</script>

<template>
  <div class="sheet-page">
    <header class="page-header">
      <div class="header-left">
        <button class="btn btn-ghost back-btn" @click="router.back()">← Back</button>
        <span class="header-icon">⚔</span>
        <span class="header-title">Behold DnD</span>
      </div>
      <button class="btn btn-ghost" @click="auth.logout()">Logout</button>
    </header>

    <div v-if="store.loading" class="loading-state">
      <div class="skeleton" style="height:40px;width:30%;margin-bottom:8px" />
      <div class="skeleton" style="height:20px;width:20%" />
    </div>
    <p v-else-if="store.error" class="page-error">{{ store.error }}</p>

    <template v-else-if="store.current && derived">
      <!-- Identity bar -->
      <div class="identity-bar">
        <div class="identity-portrait">{{ store.current.name[0] }}</div>
        <div class="identity-info">
          <h1 class="identity-name">{{ store.current.name }}</h1>
          <p class="identity-sub">
              {{ store.current.race }} ·
              {{ store.current.class }} <span class="identity-level-inline">{{ store.current.level }}</span>
              <span v-if="store.current.background"> · {{ store.current.background }}</span>
              <span v-if="store.current.alignment"> · {{ store.current.alignment }}</span>
          </p>
          <p class="identity-xp">{{ store.current.xp }} XP</p>
        </div>
        <button class="btn btn-ghost level-btn" @click="handleLevelChange">⬆ Level</button>
      </div>

      <!-- HP + rest bar -->
      <div class="hp-bar-zone">
        <div class="hp-zone-left">
          <div class="hp-numbers">
            <span class="hp-current">{{ store.current.current_hp }}</span>
            <span class="hp-sep">/</span>
            <span class="hp-max">{{ store.current.max_hp }}</span>
            <span v-if="store.current.temp_hp > 0" class="hp-temp">+{{ store.current.temp_hp }} temp</span>
          </div>
          <div class="hp-progress"><div class="hp-progress-fill" :style="{width:derived.hpPercentage.value+'%',background:hpColor(derived.hpPercentage.value)}" /></div>
          <div class="hp-hit-dice">Hit Dice: {{ store.current.hit_dice_remaining }}d{{ store.current.hit_dice_type }}</div>
        </div>
        <div class="hp-zone-mid">
          <button class="hp-action-btn" :class="{active:hpAction==='damage'}" @click="hpAction=hpAction==='damage'?null:'damage'">Damage</button>
          <button class="hp-action-btn hp-heal" :class="{active:hpAction==='heal'}" @click="hpAction=hpAction==='heal'?null:'heal'">Heal</button>
          <button class="hp-action-btn hp-temp" :class="{active:hpAction==='temp'}" @click="hpAction=hpAction==='temp'?null:'temp'">Temp HP</button>
          <template v-if="hpAction">
            <input v-model.number="hpAmount" class="input hp-input" type="number" min="1" placeholder="Amount" @keyup.enter="submitHPAction" />
            <button class="btn btn-primary btn-sm" :disabled="actionLoading||hpAmount<=0" @click="submitHPAction">Apply</button>
            <button class="btn btn-ghost btn-sm" @click="hpAction=null;hpAmount=0">✕</button>
          </template>
        </div>
        <div v-if="store.current.current_hp===0" class="death-saves-inline">
            <div class="death-row">
                <span class="death-lbl">Successes</span>
                <span
                        v-for="i in 3" :key="'s'+i"
                        class="death-dot success"
                        :class="{filled:i<=store.current.death_save_successes}"
                        @click="handleDeathSave(true)"
                        title="Record success"
                        />
            </div>
            <div class="death-row">
                <span class="death-lbl">Failures</span>
                <span
                        v-for="i in 3" :key="'f'+i"
                        class="death-dot failure"
                        :class="{filled:i<=store.current.death_save_failures}"
                        @click="handleDeathSave(false)"
                        title="Record failure"
                        />
            </div>
            <button class="btn btn-ghost" style="font-size:9px;padding:3px 8px;margin-top:4px" @click="store.update(id,{death_save_successes:0,death_save_failures:0})">Reset</button>
        </div>
        <div class="rest-zone">
          <button class="btn btn-ghost rest-btn" @click="handleShortRest">☀ Short Rest</button>
          <button class="btn btn-ghost rest-btn" @click="handleLongRest">🌙 Long Rest</button>
        </div>
      </div>

      <!-- Quickstats bar -->
      <div class="quickstats-bar">
        <div class="quickstat"><span class="qs-val">{{ store.current.armor_class }}</span><span class="qs-label">Armor Class</span></div>
        <div class="qs-div" />
        <div class="quickstat"><span class="qs-val">{{ derived.initiative.value }}</span><span class="qs-label">Initiative</span></div>
        <div class="qs-div" />
        <div class="quickstat"><span class="qs-val">{{ store.current.speed }} ft</span><span class="qs-label">Speed</span></div>
        <div class="qs-div" />
        <div class="quickstat"><span class="qs-val">{{ derived.signedModifier(derived.proficiencyBonus.value) }}</span><span class="qs-label">Proficiency</span></div>
        <div class="qs-div" />
        <div class="qs-div" />
        <div class="quickstat"><span class="qs-val">{{ store.current.inspiration ? '✦' : '—' }}</span><span class="qs-label">Inspiration</span></div>
        <template v-if="store.current.conditions?.length"><div class="qs-div" /><div class="qs-conditions"><span v-for="c in store.current.conditions" :key="c" class="condition-tag">{{ c }}</span></div></template>
      </div>

      <!-- Three columns -->
      <div class="sheet-body">
        <!-- LEFT -->
        <div class="col-left">
          <div class="panel">
            <div class="panel-header"><span class="panel-title">Ability Scores</span><button class="edit-btn" @click="startEditAbilities">✎ Edit</button></div>
            <div class="ability-grid">
              <div v-for="ab in ABILITIES" :key="ab.key" class="ability-block">
                <span class="ab-label">{{ ab.label }}</span>
                <template v-if="editingAbilities"><input v-model.number="abilityDraft[ab.key]" class="ab-edit-input" type="number" min="1" max="30" /></template>
                <template v-else><span class="ab-score">{{ store.current[ab.key as keyof typeof store.current] }}</span><span class="ab-mod">{{ derived.signedModifier(derived.modifiers.value[ab.key]) }}</span></template>
              </div>
            </div>
            <div v-if="editingAbilities" class="edit-actions"><button class="btn btn-primary btn-sm" @click="saveAbilities">Save</button><button class="btn btn-ghost btn-sm" @click="editingAbilities=false">Cancel</button></div>
          </div>

          <div class="panel">
            <div class="panel-header"><span class="panel-title">Saving Throws</span><button class="edit-btn" @click="startEditProfs">✎ Edit</button></div>
            <div class="save-list">
              <div v-for="save in SAVES" :key="save.key" class="save-row">
                <template v-if="editingProfs"><input type="checkbox" :checked="!!profDraft[save.prof]" @change="profDraft[save.prof]=!profDraft[save.prof]" class="prof-checkbox" /></template>
                <span v-else class="prof-dot" :class="{active:store.current[save.prof as keyof typeof store.current]}" />
                <span class="save-val">{{ derived.signedModifier(derived.savingThrows.value[save.key]) }}</span>
                <span class="save-label">{{ save.label }}</span>
              </div>
            </div>
            <div v-if="editingProfs" class="edit-actions"><button class="btn btn-primary btn-sm" @click="saveProfs">Save</button><button class="btn btn-ghost btn-sm" @click="editingProfs=false">Cancel</button></div>
          </div>

          <div class="panel">
            <div class="panel-header"><span class="panel-title">Senses</span></div>
            <div class="senses-list">
              <div class="sense-row"><span class="sense-label">Passive Perception</span><span class="sense-val">{{ derived.passivePerception.value }}</span></div>
              <div class="sense-row"><span class="sense-label">Passive Insight</span><span class="sense-val">{{ 10+(derived.skills.value.find(s=>s.name==='Insight')?.bonus??0) }}</span></div>
              <div class="sense-row"><span class="sense-label">Passive Investigation</span><span class="sense-val">{{ 10+(derived.skills.value.find(s=>s.name==='Investigation')?.bonus??0) }}</span></div>
            </div>
          </div>
        </div>

        <!-- MIDDLE: Skills -->
        <div class="col-mid">
          <div class="panel panel-skills">
            <div class="panel-header"><span class="panel-title">Skills</span><button class="edit-btn" @click="startEditProfs">✎ Edit</button></div>
            <div class="skill-list">
              <div v-for="skill in derived.skills.value" :key="skill.name" class="skill-row">
                <template v-if="editingProfs">
                  <span class="skill-prof-toggle" :class="`level-${profDraft[skillField(skill.name)]??0}`" @click="cycleProficiency(skillField(skill.name))">{{ profLabel(profDraft[skillField(skill.name)] as number??0) }}</span>
                </template>
                <span v-else class="prof-dot-sm" :class="{active:skill.profLevel>=1,expert:skill.profLevel===2}" />
                <span class="skill-val">{{ skill.display }}</span>
                <span class="skill-name-ability">{{ skill.name }} <span class="skill-ability-abbr">({{ skill.ability.slice(0,3).toUpperCase() }})</span></span>
              </div>
            </div>
            <div v-if="editingProfs" class="edit-actions"><button class="btn btn-primary btn-sm" @click="saveProfs">Save</button><button class="btn btn-ghost btn-sm" @click="editingProfs=false">Cancel</button></div>
          </div>
        </div>

        <!-- RIGHT: Tabs -->
        <div class="col-right">
          <div class="panel panel-tabs">
            <div class="tab-nav">
              <button v-for="tab in TABS" :key="tab.key" class="tab-btn" :class="{active:activeTab===tab.key}" @click="activeTab=tab.key">{{ tab.label }}</button>
            </div>
            <div class="tab-content">

                <div v-if="activeTab==='actions'" class="tab-panel">
                    <div class="tab-toolbar">
                        <span class="tab-count">{{ actions.length + bonusActions.length + reactions.length }} entries</span>
                        <button class="btn btn-ghost btn-sm" @click="startAddFeature('action')">+ Add</button>
                    </div>

                    <template v-if="actions.length">
                        <div class="feature-section-header">Actions</div>
                        <div v-for="f in actions" :key="f.id" class="feature-row">
                            <div class="feature-info">
                                <span class="feature-name">{{ f.name }}</span>
                                <span v-if="f.source" class="feature-source">{{ f.source }}</span>
                            </div>
                            <p v-if="f.description" class="feature-desc">{{ f.description }}</p>
                            <div class="row-actions">
                                <button class="icon-btn" @click="startEditFeature(f)">✎</button>
                                <button class="icon-btn danger" @click="deleteFeature(f.id)">✕</button>
                            </div>
                        </div>
                    </template>

                    <template v-if="bonusActions.length">
                        <div class="feature-section-header">Bonus Actions</div>
                        <div v-for="f in bonusActions" :key="f.id" class="feature-row">
                            <div class="feature-info">
                                <span class="feature-name">{{ f.name }}</span>
                                <span v-if="f.source" class="feature-source">{{ f.source }}</span>
                            </div>
                            <p v-if="f.description" class="feature-desc">{{ f.description }}</p>
                            <div class="row-actions">
                                <button class="icon-btn" @click="startEditFeature(f)">✎</button>
                                <button class="icon-btn danger" @click="deleteFeature(f.id)">✕</button>
                            </div>
                        </div>
                    </template>

                    <template v-if="reactions.length">
                        <div class="feature-section-header">Reactions</div>
                        <div v-for="f in reactions" :key="f.id" class="feature-row">
                            <div class="feature-info">
                                <span class="feature-name">{{ f.name }}</span>
                                <span v-if="f.source" class="feature-source">{{ f.source }}</span>
                            </div>
                            <p v-if="f.description" class="feature-desc">{{ f.description }}</p>
                            <div class="row-actions">
                                <button class="icon-btn" @click="startEditFeature(f)">✎</button>
                                <button class="icon-btn danger" @click="deleteFeature(f.id)">✕</button>
                            </div>
                        </div>
                    </template>

                    <template v-if="freeActions.length">
                        <div class="feature-section-header">Free Actions</div>
                        <div v-for="f in freeActions" :key="f.id" class="feature-row">
                            <div class="feature-info">
                                <span class="feature-name">{{ f.name }}</span>
                                <span v-if="f.source" class="feature-source">{{ f.source }}</span>
                            </div>
                            <p v-if="f.description" class="feature-desc">{{ f.description }}</p>
                            <div class="row-actions">
                                <button class="icon-btn" @click="startEditFeature(f)">✎</button>
                                <button class="icon-btn danger" @click="deleteFeature(f.id)">✕</button>
                            </div>
                        </div>
                    </template>

                    <p v-if="!actions.length && !bonusActions.length && !reactions.length && !freeActions.length" class="empty-tab">
                    No actions added yet. Add class features, spells or abilities.
                    </p>
                </div>

              <div v-if="activeTab==='spells'" class="tab-panel">
              <div class="spell-ability-edit">
                  <span class="spell-ability-label">Spellcasting Ability</span>
                  <select
                          class="input spell-ability-select"
                          :value="store.current.spellcasting_ability"
                          @change="store.update(id, { spellcasting_ability: ($event.target as HTMLSelectElement).value })"
                          >
                          <option value="">None</option>
                          <option value="strength">Strength</option>
                          <option value="dexterity">Dexterity</option>
                          <option value="constitution">Constitution</option>
                          <option value="intelligence">Intelligence</option>
                          <option value="wisdom">Wisdom</option>
                          <option value="charisma">Charisma</option>
                  </select>
              </div>

              <!-- Spellcasting stats -->
              <div v-if="store.current.spellcasting_ability" class="spell-stats-row">
                  <div class="spell-stat">
                      <span class="spell-stat-val">{{ store.current.spellcasting_ability.toUpperCase().slice(0,3) }}</span>
                      <span class="spell-stat-label">Ability</span>
                  </div>
                  <div class="spell-stat">
                      <span class="spell-stat-val">{{ derived.spellSaveDC.value ?? '—' }}</span>
                      <span class="spell-stat-label">Save DC</span>
                  </div>
                  <div class="spell-stat">
                      <span class="spell-stat-val">{{ derived.spellAttackBonus.value !== null ? derived.signedModifier(derived.spellAttackBonus.value) : '—' }}</span>
                      <span class="spell-stat-label">Attack Bonus</span>
                  </div>
              </div>
              <p v-else class="spell-no-ability">
              No spellcasting ability set.
              <button class="btn-link" @click="activeTab = 'personality'">Set it in character details.</button>
              </p>

                <div class="tab-toolbar"><span class="tab-count">{{ spells.length }} spell{{ spells.length!==1?'s':'' }}</span><button class="btn btn-ghost btn-sm" @click="startAddSpell">+ Add Spell</button></div>
                <div v-if="spellsLoading" class="skeleton" style="height:40px;margin-top:8px" />
                <template v-else-if="spells.length">
                  <div v-for="[level,group] in spellsByLevel" :key="level" class="spell-group">
                    <div class="spell-level-header">{{ level==='0'?'Cantrips':`Level ${level}` }}</div>
                    <div v-for="spell in group" :key="spell.id" class="spell-row">
                      <span class="spell-prepared" :class="{prepared:spell.is_prepared}" @click="togglePrepared(spell)">{{ spell.is_prepared?'◆':'◇' }}</span>
                      <span class="spell-name">{{ spell.name }}</span>
                      <span class="spell-meta">{{ spell.school }}</span>
                      <div class="row-actions"><button class="icon-btn" @click="startEditSpell(spell)">✎</button><button class="icon-btn danger" @click="deleteSpell(spell.id)">✕</button></div>
                    </div>
                  </div>
                </template>
                <p v-else class="empty-tab">No spells added yet.</p>
              </div>

              <div v-if="activeTab==='feats'" class="tab-panel">
                  <div class="tab-toolbar">
                      <span class="tab-count">{{ feats.length }} feature{{ feats.length !== 1 ? 's' : '' }}</span>
                      <button class="btn btn-ghost btn-sm" @click="startAddFeature('none')">+ Add Feature</button>
                  </div>
                  <template v-if="feats.length">
                      <div v-for="f in feats" :key="f.id" class="feature-row">
                          <div class="feature-info">
                              <span class="feature-name">{{ f.name }}</span>
                              <span v-if="f.source" class="feature-source">{{ f.source }}</span>
                          </div>
                          <p v-if="f.description" class="feature-desc">{{ f.description }}</p>
                          <div class="row-actions">
                              <button class="icon-btn" @click="startEditFeature(f)">✎</button>
                              <button class="icon-btn danger" @click="deleteFeature(f.id)">✕</button>
                          </div>
                      </div>
                  </template>
                  <p v-else class="empty-tab">No features or feats added yet.</p>
              </div>

              <div v-if="activeTab==='inventory'" class="tab-panel">
                <div class="tab-toolbar"><span class="tab-count">{{ inventory.length }} items · {{ totalWeight }} lb</span><button class="btn btn-ghost btn-sm" @click="startAddItem">+ Add Item</button></div>
                <div v-if="inventoryLoading" class="skeleton" style="height:60px;margin-top:8px" />
                <template v-else-if="inventory.length">
                  <div class="inventory-header-row"><span>Item</span><span>Qty</span><span>Wt</span><span>Val</span><span /></div>
                  <div v-for="item in inventory" :key="item.id" class="inventory-row" :class="{equipped:item.is_equipped,attuned:item.is_attuned}">
                    <div class="inv-name">{{ item.name }}<span v-if="item.is_equipped" class="inv-tag">Equipped</span><span v-if="item.is_attuned" class="inv-tag attuned-tag">Attuned</span></div>
                    <span class="inv-cell">{{ item.quantity }}</span>
                    <span class="inv-cell">{{ item.weight }}</span>
                    <span class="inv-cell">{{ (item as any).value??0 }}g</span>
                    <div class="row-actions"><button class="icon-btn" @click="startEditItem(item)">✎</button><button class="icon-btn danger" @click="deleteItem(item.id)">✕</button></div>
                  </div>
                  <div class="currency-row">
                    <span class="currency-item"><span class="copper">{{ store.current.copper }}</span> CP</span>
                    <span class="currency-item"><span class="silver">{{ store.current.silver }}</span> SP</span>
                    <span class="currency-item"><span class="electrum">{{ store.current.electrum }}</span> EP</span>
                    <span class="currency-item"><span class="gold">{{ store.current.gold }}</span> GP</span>
                    <span class="currency-item"><span class="platinum">{{ store.current.platinum }}</span> PP</span>
                  </div>
                </template>
                <p v-else class="empty-tab">No items in inventory.</p>
              </div>

              <div v-if="activeTab==='proficiencies'" class="tab-panel">
                <div v-if="store.current.training_armor?.length" class="prof-group"><span class="prof-group-label">Armor</span><p class="prof-group-vals">{{ store.current.training_armor.join(', ') }}</p></div>
                <div v-if="store.current.training_weapons?.length" class="prof-group"><span class="prof-group-label">Weapons</span><p class="prof-group-vals">{{ store.current.training_weapons.join(', ') }}</p></div>
                <div v-if="store.current.training_tools?.length" class="prof-group"><span class="prof-group-label">Tools</span><p class="prof-group-vals">{{ store.current.training_tools.join(', ') }}</p></div>
                <div v-if="store.current.training_languages?.length" class="prof-group"><span class="prof-group-label">Languages</span><p class="prof-group-vals">{{ store.current.training_languages.join(', ') }}</p></div>
                <div v-if="store.current.resistances?.length" class="prof-group"><span class="prof-group-label">Resistances</span><div class="defense-tags"><span v-for="r in store.current.resistances" :key="r" class="defense-tag resist">{{ r }}</span></div></div>
                <div v-if="store.current.immunities?.length" class="prof-group"><span class="prof-group-label">Immunities</span><div class="defense-tags"><span v-for="im in store.current.immunities" :key="im" class="defense-tag immune">{{ im }}</span></div></div>
                <div v-if="store.current.vulnerabilities?.length" class="prof-group"><span class="prof-group-label">Vulnerabilities</span><div class="defense-tags"><span v-for="v in store.current.vulnerabilities" :key="v" class="defense-tag vuln">{{ v }}</span></div></div>
                <p v-if="!store.current.training_armor?.length&&!store.current.training_weapons?.length&&!store.current.training_languages?.length" class="empty-tab">No proficiencies recorded.</p>
              </div>

              <div v-if="activeTab==='personality'" class="tab-panel">
                <div class="tab-toolbar"><span class="tab-count">Character Background</span><button class="btn btn-ghost btn-sm" @click="startEditPersonality">✎ Edit</button></div>
                <template v-if="editingPersonality">
                  <div class="personality-form">
                    <label class="field-label">Personality Traits</label><textarea v-model="personalityDraft.personality_traits" class="input textarea" rows="2" />
                    <label class="field-label">Ideals</label><textarea v-model="personalityDraft.ideals" class="input textarea" rows="2" />
                    <label class="field-label">Bonds</label><textarea v-model="personalityDraft.bonds" class="input textarea" rows="2" />
                    <label class="field-label">Flaws</label><textarea v-model="personalityDraft.flaws" class="input textarea" rows="2" />
                    <label class="field-label">Notes</label><textarea v-model="personalityDraft.notes" class="input textarea" rows="3" />
                    <div class="edit-actions"><button class="btn btn-primary btn-sm" @click="savePersonality">Save</button><button class="btn btn-ghost btn-sm" @click="editingPersonality=false">Cancel</button></div>
                  </div>
                </template>
                <template v-else>
                  <div v-if="store.current.personality_traits" class="trait-block"><span class="trait-label">Traits</span><p class="trait-text">{{ store.current.personality_traits }}</p></div>
                  <div v-if="store.current.ideals" class="trait-block"><span class="trait-label">Ideals</span><p class="trait-text">{{ store.current.ideals }}</p></div>
                  <div v-if="store.current.bonds" class="trait-block"><span class="trait-label">Bonds</span><p class="trait-text">{{ store.current.bonds }}</p></div>
                  <div v-if="store.current.flaws" class="trait-block"><span class="trait-label">Flaws</span><p class="trait-text">{{ store.current.flaws }}</p></div>
                  <div v-if="store.current.notes" class="trait-block"><span class="trait-label">Notes</span><p class="trait-text">{{ store.current.notes }}</p></div>
                  <p v-if="!store.current.personality_traits&&!store.current.ideals&&!store.current.bonds&&!store.current.flaws" class="empty-tab">No personality traits recorded.</p>
                </template>
              </div>

            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Item modal -->
    <div v-if="showItemForm" class="modal-overlay" @click.self="showItemForm=false">
      <div class="modal">
        <h3 class="modal-title">{{ editingItem?'Edit Item':'Add Item' }}</h3>
        <div class="modal-form">
          <label class="field-label">Name *</label><input v-model="itemDraft.name" class="input" type="text" />
          <div class="modal-row">
            <div><label class="field-label">Quantity</label><input v-model.number="itemDraft.quantity" class="input" type="number" min="1" /></div>
            <div><label class="field-label">Weight (lb)</label><input v-model.number="itemDraft.weight" class="input" type="number" min="0" /></div>
            <div><label class="field-label">Value (gp)</label><input v-model.number="itemDraft.value" class="input" type="number" min="0" /></div>
          </div>
          <label class="field-label">Description</label><textarea v-model="itemDraft.description" class="input textarea" rows="2" />
          <div class="modal-checks">
            <label class="check-label"><input v-model="itemDraft.is_equipped" type="checkbox" /> Equipped</label>
            <label class="check-label"><input v-model="itemDraft.requires_attunement" type="checkbox" /> Requires Attunement</label>
          </div>
        </div>
        <div class="modal-footer"><button class="btn btn-ghost" @click="showItemForm=false">Cancel</button><button class="btn btn-primary" :disabled="!itemDraft.name" @click="saveItem">Save</button></div>
      </div>
    </div>

    <!-- Feature modal -->
    <div v-if="showFeatureForm" class="modal-overlay" @click.self="showFeatureForm=false">
        <div class="modal">
            <h3 class="modal-title">{{ editingFeature ? 'Edit Feature' : 'Add Feature' }}</h3>
            <div class="modal-form">
                <label class="field-label">Name *</label>
                <input v-model="featureDraft.name" class="input" type="text" />

                <div class="modal-row" style="grid-template-columns:1fr 1fr">
                    <div>
                        <label class="field-label">Type</label>
                        <select v-model="featureDraft.action_type" class="input">
                            <option value="none">Feature / Feat</option>
                            <option value="action">Action</option>
                            <option value="bonus_action">Bonus Action</option>
                            <option value="reaction">Reaction</option>
                            <option value="free">Free Action</option>
                        </select>
                    </div>
                    <div>
                        <label class="field-label">Source</label>
                        <input v-model="featureDraft.source" class="input" type="text" placeholder="e.g. Fighter 1" />
                    </div>
                </div>

                <label class="field-label">Description</label>
                <textarea v-model="featureDraft.description" class="input textarea" rows="4" />
            </div>
            <div class="modal-footer">
                <button class="btn btn-ghost" @click="showFeatureForm=false">Cancel</button>
                <button class="btn btn-primary" :disabled="!featureDraft.name" @click="saveFeature">Save</button>
            </div>
        </div>
    </div>

    <!-- Spell modal -->
    <div v-if="showSpellForm" class="modal-overlay" @click.self="showSpellForm=false">
      <div class="modal">
        <h3 class="modal-title">{{ editingSpell?'Edit Spell':'Add Spell' }}</h3>
        <div class="modal-form">
          <label class="field-label">Name *</label><input v-model="spellDraft.name" class="input" type="text" />
          <div class="modal-row">
              <div>
                  <label class="field-label">Level</label>
                  <select v-model.number="spellDraft.level" class="input">
                      <option value="0">Cantrip</option>
                      <option v-for="l in 9" :key="l" :value="l">Level {{ l }}</option>
                  </select>
              </div>
              <div>
                  <label class="field-label">School</label>
                  <select v-model="spellDraft.school" class="input">
                      <option value="">— Select —</option>
                      <option>Abjuration</option>
                      <option>Conjuration</option>
                      <option>Divination</option>
                      <option>Enchantment</option>
                      <option>Evocation</option>
                      <option>Illusion</option>
                      <option>Necromancy</option>
                      <option>Transmutation</option>
                  </select>
              </div>
          </div>
          <div class="modal-row">
              <div>
                  <label class="field-label">Casting Time</label>
                  <select v-model="spellDraft.casting_time" class="input">
                      <option value="">— Select —</option>
                      <option>1 action</option>
                      <option>1 bonus action</option>
                      <option>1 reaction</option>
                      <option>1 minute</option>
                      <option>10 minutes</option>
                      <option>1 hour</option>
                      <option>8 hours</option>
                      <option>24 hours</option>
                  </select>
              </div>
            <div><label class="field-label">Range</label><input v-model="spellDraft.range" class="input" type="text" /></div>
          </div>
          <div class="modal-row">
              <div>
                  <label class="field-label">Components</label>
                  <div class="component-checks">
                      <label class="check-label"><input type="checkbox" :checked="spellDraft.components.includes('V')" @change="toggleComponent('V')" /> V</label>
                      <label class="check-label"><input type="checkbox" :checked="spellDraft.components.includes('S')" @change="toggleComponent('S')" /> S</label>
                      <label class="check-label"><input type="checkbox" :checked="spellDraft.components.includes('M')" @change="toggleComponent('M')" /> M</label>
                  </div>
              </div>
            <div><label class="field-label">Duration</label><input v-model="spellDraft.duration" class="input" type="text" /></div>
          </div>
          <label class="field-label">Description</label><textarea v-model="spellDraft.description" class="input textarea" rows="3" />
          <label class="check-label" style="margin-top:8px"><input v-model="spellDraft.is_prepared" type="checkbox" /> Prepared</label>
        </div>
        <div class="modal-footer"><button class="btn btn-ghost" @click="showSpellForm=false">Cancel</button><button class="btn btn-primary" :disabled="!spellDraft.name" @click="saveSpell">Save</button></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sheet-page{min-height:100vh;background:var(--bg-dark)}
.page-header{display:flex;align-items:center;justify-content:space-between;padding:10px 20px;background:var(--bg-surface);border-bottom:1px solid var(--border);position:sticky;top:0;z-index:20}
.header-left{display:flex;align-items:center;gap:12px}
.header-icon{font-size:18px;color:var(--accent-gold)}
.header-title{font-family:var(--font-deco);font-size:15px;letter-spacing:.08em}
.back-btn{padding:5px 10px;font-size:11px}
.identity-bar{display:flex;align-items:center;gap:14px;padding:12px 20px;background:var(--bg-surface);border-bottom:1px solid var(--border)}
.identity-portrait{width:44px;height:44px;border-radius:50%;background:var(--bg-card);border:2px solid var(--accent-gold-dim);display:flex;align-items:center;justify-content:center;font-family:var(--font-display);font-size:18px;color:var(--accent-gold);flex-shrink:0}
.identity-info{flex:1}
.identity-name{font-family:var(--font-display);font-size:18px;font-weight:600;letter-spacing:.04em}
.identity-sub{font-family:var(--font-body);font-size:13px;color:var(--text-secondary);margin-top:1px}
.identity-stat{text-align:center;padding:0 10px;flex-shrink:0}
.identity-level-inline {
  display: inline-flex; align-items: center; justify-content: center;
  width: 20px; height: 20px; border-radius: 50%;
  background: var(--accent-gold-dim); border: 1px solid var(--accent-gold);
  font-family: var(--font-display); font-size: 11px; font-weight: 700;
  color: var(--accent-gold); margin-left: 2px; vertical-align: middle;
}
.identity-xp { font-family: var(--font-display); font-size: 10px; letter-spacing: .08em; text-transform: uppercase; color: var(--text-muted); margin-top: 2px; }
.level-btn { padding: 5px 10px; font-size: 10px; flex-shrink: 0; }
.istat-val{display:block;font-family:var(--font-display);font-size:20px;font-weight:700;color:var(--accent-gold);line-height:1}
.istat-label{display:block;font-family:var(--font-display);font-size:9px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);margin-top:2px}
.hp-bar-zone{display:flex;align-items:center;gap:14px;flex-wrap:wrap;padding:10px 20px;background:var(--bg-card);border-bottom:1px solid var(--border)}
.hp-zone-left{display:flex;flex-direction:column;gap:4px;min-width:160px}
.hp-numbers{display:flex;align-items:baseline;gap:3px}
.hp-current{font-family:var(--font-display);font-size:28px;font-weight:700;line-height:1}
.hp-sep{font-family:var(--font-display);font-size:16px;color:var(--text-muted)}
.hp-max{font-family:var(--font-display);font-size:16px;color:var(--text-secondary)}
.hp-temp{font-family:var(--font-display);font-size:14px;color:var(--accent-gold);margin-left:6px}
.hp-progress{height:5px;background:var(--bg-surface);border-radius:3px;overflow:hidden}
.hp-progress-fill{height:100%;border-radius:3px;transition:width .4s ease,background .4s ease}
.hp-hit-dice{font-family:var(--font-display);font-size:11px;letter-spacing:.08em;text-transform:uppercase;color:var(--text-muted)}
.hp-zone-mid{display:flex;align-items:center;gap:6px;flex-wrap:wrap;flex:1}
.hp-action-btn{padding:4px 10px;font-family:var(--font-display);font-size:10px;font-weight:600;letter-spacing:.08em;text-transform:uppercase;border:1px solid var(--border-light);border-radius:var(--radius-sm);background:transparent;color:var(--text-muted);cursor:pointer;transition:all .15s}
.hp-action-btn:hover,.hp-action-btn.active{border-color:var(--accent-red-bright);color:var(--accent-red-bright)}
.hp-heal:hover,.hp-heal.active{border-color:var(--hp-high);color:#4ade80}
.hp-temp:hover,.hp-temp.active{border-color:var(--accent-gold-dim);color:var(--accent-gold)}
.hp-input{width:75px;padding:5px 8px;font-size:13px}
.death-saves-inline{display:flex;flex-direction:column;gap:3px}
.death-row{display:flex;align-items:center;gap:5px}
.death-lbl{font-family:var(--font-display);font-size:9px;letter-spacing:.08em;text-transform:uppercase;color:var(--text-muted);width:60px}
.death-dot{width:11px;height:11px;border-radius:50%;border:1px solid var(--border-light);background:transparent}
.death-dot.success.filled{background:var(--hp-high);border-color:var(--hp-high)}
.death-dot.failure.filled{background:var(--accent-red-bright);border-color:var(--accent-red-bright)}
.death-dot { cursor: pointer; transition: transform .15s, background .15s; }
.death-dot:hover { transform: scale(1.3); }
.rest-zone{display:flex;gap:6px;margin-left:auto;flex-shrink:0}
.rest-btn{padding:6px 12px;font-size:10px}
.quickstats-bar{display:flex;align-items:center;flex-wrap:wrap;padding:8px 20px;background:var(--bg-surface);border-bottom:1px solid var(--border);gap:0;overflow-x:auto;justify-content:center}
.quickstat{text-align:center;padding:0 14px}
.qs-val{display:block;font-family:var(--font-display);font-size:16px;font-weight:600;color:var(--text-primary);line-height:1}
.qs-label{display:block;font-family:var(--font-display);font-size:11px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);margin-top:2px}
.qs-div{width:1px;height:22px;background:var(--border);flex-shrink:0}
.qs-conditions{display:flex;align-items:center;gap:4px;flex-wrap:wrap;padding:0 10px}
.condition-tag{padding:2px 6px;background:rgba(139,26,26,.2);border:1px solid var(--accent-red);border-radius:2px;font-family:var(--font-display);font-size:9px;letter-spacing:.08em;text-transform:uppercase;color:#f5a0a0}
.sheet-body{display:grid;grid-template-columns:294px 294px 1fr;gap:0;min-height:calc(100vh - 240px);align-items:start}
.col-left,.col-mid,.col-right{padding:10px;border-right:1px solid var(--border);display:flex;flex-direction:column;gap:8px}
.col-right{border-right:none}
.panel{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:10px}
.panel-header{display:flex;align-items:center;justify-content:space-between;margin-bottom:8px}
.panel-title{font-family:var(--font-display);font-size:13px;font-weight:600;letter-spacing:.14em;text-transform:uppercase;color:var(--text-muted)}
.edit-btn{font-family:var(--font-display);font-size:11px;letter-spacing:.06em;color:var(--accent-gold-dim);background:none;border:none;cursor:pointer;padding:2px 4px;transition:color .15s}
.edit-btn:hover{color:var(--accent-gold)}
.ability-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:5px}
.ability-block{background:var(--bg-surface);border:1px solid var(--border-light);border-radius:var(--radius-sm);text-align:center;padding:7px 4px;display:flex;flex-direction:column;align-items:center;gap:1px}
.ab-label{font-family:var(--font-display);font-size:12px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted)}
.ab-score{font-family:var(--font-display);font-size:17px;font-weight:600;line-height:1}
.ab-mod{font-family:var(--font-display);font-size:13px;color:var(--accent-gold)}
.ab-edit-input{width:100%;padding:2px 4px;background:var(--bg-input);border:1px solid var(--border-light);border-radius:2px;color:var(--text-primary);font-family:var(--font-display);font-size:13px;text-align:center;outline:none}
.save-list{display:flex;flex-direction:column;gap:2px}
.save-row{display:flex;align-items:center;gap:5px;padding:2px 0}
.prof-dot{width:9px;height:9px;border-radius:50%;border:1px solid var(--border-light);background:transparent;flex-shrink:0}
.prof-dot.active{background:var(--accent-gold);border-color:var(--accent-gold)}
.prof-checkbox{width:12px;height:12px;cursor:pointer;accent-color:var(--accent-gold);flex-shrink:0}
.save-val{font-family:var(--font-display);font-size:14px;font-weight:600;width:22px;text-align:right}
.save-label{font-family:var(--font-body);font-size:13px;color:var(--text-secondary)}
.senses-list{display:flex;flex-direction:column;gap:3px}
.sense-row{display:flex;justify-content:space-between;align-items:center;padding:2px 0}
.sense-label{font-family:var(--font-body);font-size:13px;color:var(--text-secondary)}
.sense-val{font-family:var(--font-display);font-size:14px;font-weight:600;color:var(--accent-gold)}
.panel-skills{height:fit-content}
.skill-list{display:flex;flex-direction:column;gap:1.3px}
.skill-row{display:flex;align-items:center;gap:4px;padding:2px 2px;border-radius:2px;transition:background .1s}
.skill-row:hover{background:var(--bg-surface)}
.prof-dot-sm{width:10px;height:10px;border-radius:50%;flex-shrink:0;border:1px solid var(--border-light);background:transparent}
.prof-dot-sm.active{background:var(--accent-gold);border-color:var(--accent-gold)}
.prof-dot-sm.expert{background:var(--accent-gold);border-color:var(--accent-gold);box-shadow:0 0 0 1px var(--accent-gold-dim)}
.skill-prof-toggle{width:13px;height:13px;border-radius:50%;flex-shrink:0;display:flex;align-items:center;justify-content:center;font-size:8px;cursor:pointer;border:1px solid var(--border-light);color:var(--text-muted);transition:all .1s}
.skill-prof-toggle.level-1{background:var(--accent-gold-dim);border-color:var(--accent-gold);color:var(--accent-gold)}
.skill-prof-toggle.level-2{background:var(--accent-gold);border-color:var(--accent-gold);color:#000}
.skill-val{font-family:var(--font-display);font-size:15px;font-weight:600;width:20px;text-align:right;flex-shrink:0}
.skill-name-ability{font-family:var(--font-body);font-size:13px;color:var(--text-secondary);line-height:1.2}
.skill-ability-abbr{font-family:var(--font-display);font-size:11px;color:var(--text-muted)}
.edit-actions{display:flex;gap:5px;margin-top:8px}
.btn-sm{padding:5px 10px;font-size:11px}
.panel-tabs{padding:0;overflow:hidden}
.tab-nav{display:flex;overflow-x:auto;background:var(--bg-surface);border-bottom:1px solid var(--border)}
.tab-btn{padding:9px 12px;flex-shrink:0;font-family:var(--font-display);font-size:13px;font-weight:600;letter-spacing:.08em;text-transform:uppercase;color:var(--text-muted);background:none;border:none;border-bottom:2px solid transparent;cursor:pointer;transition:color .15s,border-color .15s}
.tab-btn:hover{color:var(--text-secondary)}
.tab-btn.active{color:var(--accent-gold);border-bottom-color:var(--accent-gold)}
.tab-content{padding:12px}
.tab-panel{display:flex;flex-direction:column;gap:6px}
.tab-toolbar{display:flex;align-items:center;justify-content:space-between;margin-bottom:4px}
.tab-count{font-family:var(--font-display);font-size:10px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted)}
.empty-tab{font-family:var(--font-body);font-style:italic;color:var(--text-muted);font-size:14px;padding:6px 0}
.tab-empty-hint{font-family:var(--font-body);font-size:13px;color:var(--text-muted);font-style:italic}
.inventory-header-row{display:grid;grid-template-columns:1fr 36px 36px 48px 40px;gap:6px;padding:3px 4px;font-family:var(--font-display);font-size:11px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);border-bottom:1px solid var(--border)}
.inventory-row{display:grid;grid-template-columns:1fr 36px 36px 48px 40px;gap:6px;padding:5px 4px;align-items:center;border-bottom:1px solid var(--border);transition:background .1s}
.inventory-row:hover{background:var(--bg-surface)}
.inventory-row.equipped{border-left:2px solid var(--accent-gold-dim)}
.inventory-row.attuned{border-left:2px solid #7ab3e0}
.inv-name{font-family:var(--font-body);font-size:15px;display:flex;align-items:center;gap:4px;flex-wrap:wrap}
.inv-tag{padding:1px 4px;border-radius:2px;font-family:var(--font-display);font-size:12px;letter-spacing:.06em;text-transform:uppercase;background:var(--accent-gold-dim);color:var(--accent-gold);border:1px solid var(--accent-gold-dim)}
.attuned-tag{background:rgba(26,58,92,.4);border-color:#1a3a5c;color:#7ab3e0}
.inv-cell{font-family:var(--font-display);font-size:11px;color:var(--text-secondary);text-align:center}
.currency-row{display:flex;gap:10px;padding:6px 2px;margin-top:2px;border-top:1px solid var(--border);font-family:var(--font-display);font-size:13px;color:var(--text-muted)}
.currency-item{display:flex;align-items:center;gap:3px}
.copper{color:#b87333;font-weight:600}
.silver{color:#aaa9ad;font-weight:600}
.electrum{color:#7ecece;font-weight:600}
.gold{color:var(--accent-gold);font-weight:600}
.platinum{color:#e5e4e2;font-weight:600}
.feature-section-header {
  font-family: var(--font-display); font-size: 9px; letter-spacing: .12em;
  text-transform: uppercase; color: var(--text-muted);
  border-bottom: 1px solid var(--border); padding: 3px 0; margin: 6px 0 3px;
}
.feature-row {
  padding: 6px 4px; border-bottom: 1px solid var(--border);
  display: flex; flex-direction: column; gap: 3px;
  transition: background .1s;
}
.feature-row:hover { background: var(--bg-surface); }
.feature-info { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.feature-name { font-family: var(--font-body); font-size: 14px; flex: 1; }
.feature-source { font-family: var(--font-display); font-size: 9px; letter-spacing: .06em; color: var(--text-muted); text-transform: uppercase; }
.feature-desc { font-family: var(--font-body); font-size: 12px; color: var(--text-secondary); font-style: italic; line-height: 1.4; padding: 0 2px; }

.spell-stats-row {
  display: flex; gap: 0; background: var(--bg-surface);
  border: 1px solid var(--border); border-radius: var(--radius-sm);
  margin-bottom: 6px; overflow: hidden;
}
.spell-stat { flex: 1; text-align: center; padding: 8px 4px; border-right: 1px solid var(--border); }
.spell-stat:last-child { border-right: none; }
.spell-stat-val { display: block; font-family: var(--font-display); font-size: 16px; font-weight: 600; color: var(--accent-gold); }
.spell-stat-label { display: block; font-family: var(--font-display); font-size: 9px; letter-spacing: .1em; text-transform: uppercase; color: var(--text-muted); margin-top: 2px; }
.spell-no-ability { font-family: var(--font-body); font-size: 13px; color: var(--text-muted); font-style: italic; }
.spell-ability-edit {
  display: flex; align-items: center; justify-content: space-between;
  gap: 12px; margin-top: 10px; padding-top: 10px;
  border-top: 1px solid var(--border);
}
.spell-ability-label {
  font-family: var(--font-display); font-size: 10px; letter-spacing: .1em;
  text-transform: uppercase; color: var(--text-muted); flex-shrink: 0;
}
.spell-ability-select { width: 140px; padding: 5px 8px; font-size: 13px; }
.btn-link { background: none; border: none; color: var(--accent-gold-dim); font-family: var(--font-body); font-size: 13px; cursor: pointer; text-decoration: underline; padding: 0; }
.btn-link:hover { color: var(--accent-gold); }
.spell-group{margin-bottom:6px}
.spell-level-header{font-family:var(--font-display);font-size:13px;letter-spacing:.12em;text-transform:uppercase;color:var(--text-muted);border-bottom:1px solid var(--border);padding:3px 0;margin-bottom:3px}
.spell-row{display:flex;align-items:center;gap:6px;padding:4px 3px;border-radius:2px;transition:background .1s}
.spell-row:hover{background:var(--bg-surface)}
.spell-prepared{font-size:12px;cursor:pointer;color:var(--text-muted);transition:color .15s;flex-shrink:0}
.spell-prepared.prepared{color:var(--accent-gold)}
.spell-name{font-family:var(--font-body);font-size:15px;flex:1}
.spell-meta{font-family:var(--font-display);font-size:9px;color:var(--text-muted);letter-spacing:.06em}
.component-checks { display: flex; gap: 12px; margin-top: 4px; }
select.input { cursor: pointer; }
select.input option { background: var(--bg-card); color: var(--text-primary); }
.row-actions{display:flex;gap:3px;flex-shrink:0}
.icon-btn{width:20px;height:20px;display:flex;align-items:center;justify-content:center;background:none;border:1px solid var(--border);border-radius:2px;color:var(--text-muted);font-size:10px;cursor:pointer;transition:all .15s}
.icon-btn:hover{border-color:var(--accent-gold-dim);color:var(--accent-gold)}
.icon-btn.danger:hover{border-color:var(--accent-red);color:var(--accent-red-bright)}
.prof-group{margin-bottom:8px}
.prof-group-label{display:block;font-family:var(--font-display);font-size:10px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);margin-bottom:2px}
.prof-group-vals{font-family:var(--font-body);font-size:13px;color:var(--text-secondary)}
.defense-tags{display:flex;flex-wrap:wrap;gap:3px;margin-top:3px}
.defense-tag{padding:2px 6px;border-radius:2px;font-family:var(--font-display);font-size:9px;letter-spacing:.06em;text-transform:uppercase}
.resist{background:rgba(26,58,92,.3);border:1px solid #1a3a5c;color:#7ab3e0}
.immune{background:rgba(45,122,79,.2);border:1px solid var(--hp-high);color:#4ade80}
.vuln{background:rgba(139,26,26,.2);border:1px solid var(--accent-red);color:#f5a0a0}
.trait-block{margin-bottom:8px}
.trait-label{display:block;font-family:var(--font-display);font-size:10px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);margin-bottom:2px}
.trait-text{font-family:var(--font-body);font-style:italic;font-size:13px;color:var(--text-secondary);line-height:1.5}
.personality-form{display:flex;flex-direction:column;gap:7px}
.field-label{font-family:var(--font-display);font-size:10px;letter-spacing:.1em;text-transform:uppercase;color:var(--text-muted);display:block;margin-bottom:2px}
.textarea{resize:vertical;font-family:var(--font-body);line-height:1.5}
.modal-overlay{position:fixed;inset:0;background:rgba(0,0,0,.78);display:flex;align-items:center;justify-content:center;z-index:100;padding:20px}
.modal{background:var(--bg-card);border:1px solid var(--border-light);border-radius:var(--radius-md);padding:22px;width:100%;max-width:460px;max-height:90vh;overflow-y:auto}
.modal-title{font-family:var(--font-display);font-size:15px;font-weight:600;letter-spacing:.05em;margin-bottom:14px}
.modal-form{display:flex;flex-direction:column;gap:9px}
.modal-row{display:grid;grid-template-columns:1fr 1fr 1fr;gap:9px}
.modal-checks{display:flex;gap:14px}
.check-label{display:flex;align-items:center;gap:5px;cursor:pointer;font-family:var(--font-body);font-size:13px;color:var(--text-secondary)}
.check-label input{accent-color:var(--accent-gold)}
.modal-footer{display:flex;justify-content:flex-end;gap:7px;margin-top:14px}
.loading-state{padding:48px 24px}
.page-error{color:var(--accent-red-bright);text-align:center;padding:32px}
</style>
