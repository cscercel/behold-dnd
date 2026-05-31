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

func (s *CharacterService) CreateCharacter(
	ctx context.Context, characterData db.CreateCharacterParams,
) (db.Character, error) {
	character, err := s.queries.CreateCharacter(ctx, characterData)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to create character: %w", err)
	}

	return character, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, characterID uuid.UUID) (db.GetCharacterRow, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.GetCharacterRow{}, fmt.Errorf("failed to get character: %w", err)
	}

	return character, nil
}

func (s *CharacterService) ListCharacters(ctx context.Context) ([]db.Character, error) {
	characters, err := s.queries.ListCharacters(ctx)
	if err != nil {
		return []db.Character{}, fmt.Errorf("failed to list characters: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) ListUserCharacters(ctx context.Context, userID uuid.UUID) ([]db.Character, error) {
	characters, err := s.queries.ListUserCharacters(ctx, userID)
	if err != nil {
		return []db.Character{}, fmt.Errorf("failed to list user characters: %w", err)
	}
	return characters, nil
}

func (s *CharacterService) ListPlayerCharacters(ctx context.Context) ([]db.Character, error) {
	characters, err := s.queries.ListPlayerCharacters(ctx)
	if err != nil {
		return []db.Character{}, fmt.Errorf("failed to list player characters: %w", err)
	}

	return characters, nil
}

func (s *CharacterService) ListNPCs(ctx context.Context) ([]db.Character, error) {
	npcs, err := s.queries.ListNPCs(ctx)
	if err != nil {
		return []db.Character{}, fmt.Errorf("failed to list npcs: %w", err)
	}

	return npcs, nil
}

func (s *CharacterService) UpdateCharacter(
	ctx context.Context, characterData db.UpdateCharacterParams,
) (db.Character, error) {
	character, err := s.queries.UpdateCharacter(ctx, characterData)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to update character: %w", err)
	}

	return character, nil
}

func (s *CharacterService) HealCharacter(
	ctx context.Context, characterID uuid.UUID, amount int32,
) (db.Character, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to get character: %w", err)
	}

	newHP := min((character.CurrentHp) + amount, character.MaxHp)
	
	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID: character.ID,
		CurrentHp: newHP,
		TempHp: character.TempHp,
	})
}

func (s *CharacterService) DamageCharacter(
	ctx context.Context, characterID uuid.UUID, amount int32,
) (db.Character, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to get character: %w", err)
	}

	tempHP := character.TempHp
	currentHP := character.CurrentHp

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
		ID:	character.ID,
		CurrentHp: int32(currentHP),
		TempHp: int32(tempHP),
	})
}

func (s *CharacterService) UpdateCharacterTempHP(
	ctx context.Context, characterID uuid.UUID, amount int32,
) (db.Character, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to get character: %w", err)
	}

	return s.queries.UpdateCharacterHP(ctx, db.UpdateCharacterHPParams{
		ID: character.ID,
		CurrentHp: character.CurrentHp,
		TempHp: amount,
	})
}

func (s *CharacterService) UpdateCharacterDeathSave(
	ctx context.Context, characterID uuid.UUID, success bool,
) (db.Character, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.Character{}, fmt.Errorf("failed to get character: %w", err)
	}

	successes := character.DeathSaveSuccesses
	failures := character.DeathSaveFailures

	if success {
		successes = min(successes + 1, 3)
	} else {
		failures = min(failures + 1, 3)
	}

	return s.queries.UpdateDeathSaves(ctx, db.UpdateDeathSavesParams{
		ID:	character.ID,
		DeathSaveSuccesses: int32(successes),
		DeathSaveFailures: int32(failures),
	})
}
