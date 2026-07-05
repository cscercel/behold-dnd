package service

import (
	"context"
	"fmt"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

type CharacterService struct {
	queries      *db.Queries
	spellService *SpellService
}

func NewCharacterService(queries *db.Queries, spellService *SpellService) *CharacterService {
	return &CharacterService{queries: queries, spellService: spellService}
}

func (s *CharacterService) ListCharacters(ctx context.Context) ([]db.Character, error) {
	characters, err := s.queries.ListCharacters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list characters: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) ListUserCharacters(ctx context.Context, ownerID uuid.UUID) ([]db.Character, error) {
	characters, err := s.queries.ListUserCharacters(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list characters: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) ListPlayerCharacters(ctx context.Context) ([]db.Character, error) {
	characters, err := s.queries.ListPlayerCharacters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list player characters: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) ListNPCs(ctx context.Context) ([]db.Character, error) {
	characters, err := s.queries.ListNPCs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list NPCs: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, id uuid.UUID) (db.GetCharacterRow, error) {
	character, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.GetCharacterRow{}, fmt.Errorf("character not found: %w", err)
	}
	return character, nil
}

func (s *CharacterService) CreateCharacter(ctx context.Context, params db.CreateCharacterParams) (db.Character, error) {
	character, err := s.queries.CreateCharacter(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to create character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterInfo(ctx context.Context, params db.UpdateCharacterInfoParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterInfo(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterAbilityScores(ctx context.Context, params db.UpdateCharacterAbilityScoresParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterAbilityScores(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterSkills(ctx context.Context, params db.UpdateCharacterSkillsParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterSkills(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterLevel(ctx context.Context, params db.UpdateCharacterLevelParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterLevel(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterTraining(ctx context.Context, params db.UpdateCharacterTrainingParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterTraining(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateCharacterCurrency(ctx context.Context, params db.UpdateCharacterCurrencyParams) (db.Character, error) {
	character, err := s.queries.UpdateCharacterCurrency(ctx, params)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}
	return character, nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteCharacter(ctx, id); err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}
	return nil
}

// Apply damage to temp HP before real HP and stop at 0
func (s *CharacterService) ApplyDamage(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	tempHP := int(char.TempHp)
	currentHP := int(char.CurrentHp)

	if tempHP > 0 {
		if amount <= tempHP {
			tempHP -= amount
			amount = 0
		} else {
			amount -= tempHP
			tempHP = 0
		}
	}

	currentHP = max(currentHP-amount, 0)

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID:        id,
		CurrentHp: int32(currentHP),
		TempHp:    int32(tempHP),
	})
}

// Heal up to Max HP
func (s *CharacterService) Heal(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	newHP := min(int(char.CurrentHp)+amount, int(char.MaxHp))

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID:        id,
		CurrentHp: int32(newHP),
		TempHp:    char.TempHp,
	})
}

// Temp HP does not stack
func (s *CharacterService) AddTempHP(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID:        id,
		CurrentHp: char.CurrentHp,
		TempHp:    int32(amount),
	})
}

// Record death saving throws
func (s *CharacterService) RecordDeathSave(ctx context.Context, id uuid.UUID, success bool) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	successes := int(char.DeathSaveSuccesses)
	failures := int(char.DeathSaveFailures)

	if success {
		successes = min(successes+1, 3)
	} else {
		failures = min(failures+1, 3)
	}

	return s.queries.UpdateDeathSaves(ctx, db.UpdateDeathSavesParams{
		ID:                 id,
		DeathSaveSuccesses: int32(successes),
		DeathSaveFailures:  int32(failures),
	})
}

// LongRest restores HP/hit dice/death saves/conditions and resets spell slots.
func (s *CharacterService) LongRest(ctx context.Context, id uuid.UUID) (db.Character, error) {
	character, err := s.queries.LongRest(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to apply long rest: %w", err)
	}

	if err := s.spellService.LongRestSlots(ctx, id); err != nil {
		return db.Character{}, fmt.Errorf("failed to reset spell slots: %w", err)
	}

	return character, nil
}

func (s *CharacterService) ShortRest(ctx context.Context, id uuid.UUID, hitDiceRemaining, currentHP int) (db.Character, error) {
	character, err := s.queries.ShortRest(ctx, db.ShortRestParams{
		ID:               id,
		HitDiceRemaining: int32(hitDiceRemaining),
		CurrentHp:        int32(currentHP),
	})
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to apply short rest: %w", err)
	}
	return character, nil
}

func (s *CharacterService) UpdateConditions(ctx context.Context, id uuid.UUID, conditions []string) (db.Character, error) {
	character, err := s.queries.UpdateConditions(ctx, db.UpdateConditionsParams{
		ID:         id,
		Conditions: conditions,
	})
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update conditions: %w", err)
	}
	return character, nil
}
