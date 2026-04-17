import { computed, initCustomFormatter } from 'vue'
import type { Character } from '@/types'


const SKILLS = [
  { name: 'Acrobatics',      ability: 'dexterity',      field: 'skill_acrobatics'      },
  { name: 'Animal Handling', ability: 'wisdom',         field: 'skill_animal_handling' },
  { name: 'Arcana',          ability: 'intelligence',   field: 'skill_arcana'          },
  { name: 'Athletics',       ability: 'strength',       field: 'skill_athletics'       },
  { name: 'Deception',       ability: 'charisma',       field: 'skill_deception'       },
  { name: 'History',         ability: 'intelligence',   field: 'skill_history'         },
  { name: 'Insight',         ability: 'wisdom',         field: 'skill_insight'         },
  { name: 'Intimidation',    ability: 'charisma',       field: 'skill_intimidation'    },
  { name: 'Investigation',   ability: 'intelligence',   field: 'skill_investigation'   },
  { name: 'Medicine',        ability: 'wisdom',         field: 'skill_medicine'        },
  { name: 'Nature',          ability: 'intelligence',   field: 'skill_nature'          },
  { name: 'Perception',      ability: 'wisdom',         field: 'skill_perception'      },
  { name: 'Performance',     ability: 'charisma',       field: 'skill_performance'     },
  { name: 'Persuasion',      ability: 'charisma',       field: 'skill_persuasion'      },
  { name: 'Religion',        ability: 'intelligence',   field: 'skill_religion'        },
  { name: 'Sleight of Hand', ability: 'dexterity',      field: 'skill_sleight_of_hand' },
  { name: 'Stealth',         ability: 'dexterity',      field: 'skill_stealth'         },
  { name: 'Survival',        ability: 'wisdom',         field: 'skill_survival'        },
] as const

type AbilityKey = 'strength' | 'dexterity' | 'constitution' | 'intelligence' | 'wisdom' | 'charisma'

export function useCharacter(character: Character) {
    const abilityModifier = (score: number): number =>
        Math.floor((score - 10) / 2)

    const proficiencyBonus = computed((): number =>
        Math.ceil(1 + character.level / 4)
    )

    const signedModifier = (value: number): string =>
        value >= 0 ? `+${value}` : `${value}`

    const modifiers = computed(() => ({
        strength:       abilityModifier(character.strength),
        dexterity:      abilityModifier(character.dexterity),
        constitution:   abilityModifier(character.constitution),
        intelligence:   abilityModifier(character.intelligence),
        wisdom:         abilityModifier(character.wisdom),
        charisma:       abilityModifier(character.charisma),
    }))

    const savingThrows = computed(() => ({
        strength: modifiers.value.strength + (character.save_prof_strength ? proficiencyBonus.value : 0),
        dexterity: modifiers.value.dexterity + (character.save_prof_dexterity ? proficiencyBonus.value : 0),
        constitution: modifiers.value.constitution + (character.save_prof_constitution ? proficiencyBonus.value : 0),
        intelligence: modifiers.value.intelligence + (character.save_prof_intelligence ? proficiencyBonus.value : 0),
        wisdom: modifiers.value.wisdom + (character.save_prof_wisdom ? proficiencyBonus.value : 0),
        charisma: modifiers.value.charisma + (character.save_prof_charisma ? proficiencyBonus.value : 0),
    }))

    const skills = computed(() =>
        SKILLS.map((skill) => {
            const abilityMod = modifiers.value[skill.ability as AbilityKey]
            const profLevel = character[skill.field as keyof Character] as number
            const bonus = abilityMod + profLevel * proficiencyBonus.value
            return {
                name: skill.name,
                ability: skill.ability,
                profLevel,
                bonus,
                display: signedModifier(bonus),
            }
        })
    )

    const passivePerception = computed(() => {
        const perception = skills.value.find((s) => s.name === 'Perception')
        return 10 + (perception?.bonus ?? 0)
    })

    const initiative = computed(() =>
        signedModifier(modifiers.value.dexterity)
    )

    const hpPercentage = computed(() =>
        character.max_hp > 0
        ? Math.round((character.current_hp / character.max_hp) * 100)
        : 0
    )

    return {
        proficiencyBonus,
        modifiers,
        savingThrows,
        skills,
        passivePerception,
        initiative,
        hpPercentage,
        signedModifier,
        abilityModifier,
    }
}
