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

func (s *CombatService) ListEncounters(ctx context.Context) ([]db.CombatEncounter, error) {
	encounters, err := s.queries.ListEncounters(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list encounters: %w", err)
	}
	return encounters, nil
}

func (s *CombatService) CreateEncounter(ctx context.Context, name string) (db.CombatEncounter, error) {
	encounter, err := s.queries.CreateEncounter(ctx, name)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to create encounter: %w", err)
	}
	return encounter, nil
}

func (s *CombatService) GetEncounter(ctx context.Context, id uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.GetEncounter(ctx, id)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("encounter not found: %w", err)
	}
	return encounter, nil
}

func (s *CombatService) StartEncounter(ctx context.Context, id uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.StartEncounter(ctx, id)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to start encounter: %w", err)
	}
	return encounter, nil
}

func (s *CombatService) EndEncounter(ctx context.Context, id uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.EndEncounter(ctx, id)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to end encounter: %w", err)
	}
	return encounter, nil
}

func (s *CombatService) NextRound(ctx context.Context, id uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.NextRound(ctx, id)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to advance round: %w", err)
	}
	return encounter, nil
}

func (s *CombatService) DeleteEncounter(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteEncounter(ctx, id); err != nil {
		return fmt.Errorf("failed to delete encounter: %w", err)
	}
	return nil
}

func (s *CombatService) ListParticipants(ctx context.Context, encounterID uuid.UUID) ([]db.CombatParticipant, error) {
	participants, err := s.queries.ListParticipants(ctx, encounterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}
	return participants, nil
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

func (s *CombatService) RemoveParticipant(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.RemoveParticipant(ctx, id); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}
	return nil
}

// Apply damage to participant, knock out if current hp reaches 0.
// Temp HP is used first.
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

	p, err = s.queries.UpdateParticipantTempHP(ctx, db.UpdateParticipantTempHPParams{
		ID:     participantID,
		TempHp: int32(tempHP),
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update temp HP: %w", err)
	}

	p, err = s.queries.UpdateParticipantHP(ctx, db.UpdateParticipantHPParams{
		ID:        participantID,
		CurrentHp: int32(currentHP),
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update current HP: %w", err)
	}

	if currentHP == 0 {
		p, err = s.queries.DeactivateParticipant(ctx, participantID)
		if err != nil {
			return db.CombatParticipant{}, fmt.Errorf("failed to deactivate participant: %w", err)
		}
	}

	return p, nil
}

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

func (s *CombatService) AddParticipantTempHP(ctx context.Context, participantID uuid.UUID, amount int32) (db.CombatParticipant, error) {
	p, err := s.queries.UpdateParticipantTempHP(ctx, db.UpdateParticipantTempHPParams{
		ID:     participantID,
		TempHp: amount,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update temp HP: %w", err)
	}
	return p, nil
}

func (s *CombatService) UpdateParticipantInitiative(ctx context.Context, participantID uuid.UUID, initiative int32) (db.CombatParticipant, error) {
	p, err := s.queries.UpdateParticipantInitiative(ctx, db.UpdateParticipantInitiativeParams{
		ID:         participantID,
		Initiative: initiative,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update initiative: %w", err)
	}
	return p, nil
}

func (s *CombatService) UpdateParticipantConditions(ctx context.Context, participantID uuid.UUID, conditions []string) (db.CombatParticipant, error) {
	p, err := s.queries.UpdateParticipantConditions(ctx, db.UpdateParticipantConditionsParams{
		ID:         participantID,
		Conditions: conditions,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update conditions: %w", err)
	}
	return p, nil
}

func (s *CombatService) ToggleParticipantConcentration(ctx context.Context, participantID uuid.UUID) (db.CombatParticipant, error) {
	p, err := s.queries.ToggleParticipantConcentration(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to toggle concentration: %w", err)
	}
	return p, nil
}

func (s *CombatService) DeactivateParticipant(ctx context.Context, participantID uuid.UUID) (db.CombatParticipant, error) {
	p, err := s.queries.DeactivateParticipant(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to deactivate participant: %w", err)
	}
	return p, nil
}
