package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/cscercel/behold-dnd/internal/db"
)


type CharacterService struct {
	queries *db.Queries
}

func NewCharacterService(queries *db.Queries) *CharacterService {
	return &CharacterService{queries: queries}
}

// Apply damage to temp HP before real HP and stop at 0
func (s *CharacterService) ApplyDamage(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	tempHP := int(char.TempHp)
	currentHP := int(char.CurrentHp)

	// Hit temp HP first
	if tempHP > 0 {
		if amount <= tempHP {
			tempHP -= amount
			amount = 0
		} else {
			amount -= tempHP
			tempHP = 0
		}
	}

	// Apply rest to real HP up until it reaches 0
	currentHP = max(currentHP - amount, 0)

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID:	id,
		CurrentHp: int32(currentHP),
		TempHp: int32(tempHP),
	})
}

// Heal up to Max HP
func (s *CharacterService) Heal(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	newHP := min(int(char.CurrentHp) + amount, int(char.MaxHp))

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID:	id,
		CurrentHp: int32(newHP),
		TempHp: char.TempHp,
	})
}

// Temp HP does not stack, we would update and take max buff
func (s *CharacterService) AddTempHP(ctx context.Context, id uuid.UUID, amount int) (db.Character, error) {
	char, err := s.queries.GetCharacter(ctx, id)
	if err != nil {
		return db.Character{}, fmt.Errorf("character not found: %w", err)
	}

	newTempHP := amount

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID: id,
		CurrentHp: char.CurrentHp,
		TempHp: int32(newTempHP),
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
		successes = min(successes + 1, 3)
	} else {
		failures = min(failures + 1, 3)
	}

	return s.queries.UpdateDeathSaves(ctx, db.UpdateDeathSavesParams{
		ID:	id,
		DeathSaveSuccesses: int32(successes),
		DeathSaveFailures: int32(failures),
	})
}
