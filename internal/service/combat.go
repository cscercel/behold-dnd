package service

import (
	"context"
	"fmt"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

type CombatService struct {
	queries *db.Queries
}

func NewCombatService(queries *db.Queries) *CombatService {
	return &CombatService{queries: queries}
}

// Add character with initiative roll already done
func (s *CombatService) AddCharacterToEncounter(
	ctx context.Context,
	encounterID uuid.UUID,
	characterID uuid.UUID,
	initiative int32,
) (db.CombatParticipant, error) {
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("character not found: %w", err)
	}

	return s.queries.AddParticipant(ctx, db.AddParticipantParams{
		EncounterID: encounterID,
		CharacterID: characterID,
		Name:        character.Name,
		Initiative:  initiative,
		CurrentHp:   character.CurrentHp,
		MaxHp:       character.MaxHp,
		TempHp:      character.TempHp,
		ArmorClass:  character.ArmorClass,
		Speed:       character.Speed,
	})
}

// Apply damage to participant, knock out if current hp reaches 0
// Temp HP should be used first
func (s *CombatService) ApplyDamageToParticipant(
	ctx context.Context,
	participantID uuid.UUID,
	amount int,
) (db.CombatParticipant, error) {
	p, err := s.queries.GetParticipant(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("participant not found: %w", err)
	}

	currentHP := int(p.CurrentHp)
	tempHP := int(p.TempHp)

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
	currentHP = max(currentHP-amount, 0)

	// Update temp HP
	p, err = s.queries.UpdateParticipantTempHP(ctx, db.UpdateParticipantTempHPParams{
		ID:     participantID,
		TempHp: int32(tempHP),
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update temp HP: %w", err)
	}

	// Update current HP
	p, err = s.queries.UpdateParticipantHP(ctx, db.UpdateParticipantHPParams{
		ID:        participantID,
		CurrentHp: int32(currentHP),
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update current HP: %w", err)
	}

	// Knock out at 0 HP
	if currentHP == 0 {
		p, err = s.queries.DeactivateParticipant(ctx, participantID)
		if err != nil {
			return db.CombatParticipant{}, fmt.Errorf("failed to deactivate participant: %w", err)
		}
	}

	return p, nil
}

// Heal participant HP
func (s *CombatService) HealParticipant(
	ctx context.Context,
	participantID uuid.UUID,
	amount int,
) (db.CombatParticipant, error) {
	p, err := s.queries.GetParticipant(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("participant not found: %w", err)
	}

	newHP := min(int(p.CurrentHp)+amount, int(p.MaxHp))

	p, err = s.queries.UpdateParticipantHP(ctx, db.UpdateParticipantHPParams{
		ID:        participantID,
		CurrentHp: int32(newHP),
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update HP: %w", err)
	}

	return p, nil
}
