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

func NewCombarService(queries *db.Queries) *CombatService {
	return &CombatService{queries: queries}
}

func (s *CombatService) CreateEncouter(ctx context.Context, name string) (db.CombatEncounter, error) {
	encounter, err := s.queries.CreateEncounter(ctx, name)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to create encounter: %w", err)
	}

	return encounter, nil
}

func (s *CombatService) GetEncounter(ctx context.Context, combatID uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.GetEncounter(ctx, combatID)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to get encounter: %w", err)
	}
	
	return encounter, nil
}

func (s *CombatService) ListEncouters(ctx context.Context) ([]db.CombatEncounter, error) {
	encounters, err := s.queries.ListEncounters(ctx)
	if err != nil {
		return []db.CombatEncounter{}, fmt.Errorf("failed to list encounters: %w", err)
	}
	
	return encounters, nil
}

func (s *CombatService) GetActiveEncouter(ctx context.Context) (db.CombatEncounter, error) {
	encounter, err := s.queries.GetActiveEncounter(ctx)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to list active encounter: %w", err)
	}

	return encounter, nil
}

func (s *CombatService) StartEncounter(ctx context.Context, combatID uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.StartEncounter(ctx, combatID)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to start encounter: %w", err)
	}

	return encounter, nil
}

func (s *CombatService) EndEncounter(ctx context.Context, combatID uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.EndEncounter(ctx, combatID)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to end encounter: %w", err)
	}

	return encounter, nil
}

func (s *CombatService) NextRound(ctx context.Context, combatID uuid.UUID) (db.CombatEncounter, error) {
	encounter, err := s.queries.NextRound(ctx, combatID)
	if err != nil {
		return db.CombatEncounter{}, fmt.Errorf("failed to advance to next round: %w", err)
	}

	return encounter, nil
}

func (s *CombatService) DeleteEncounter(ctx context.Context, combatID uuid.UUID) error {
	if err := s.queries.DeleteEncounter(ctx, combatID); err != nil {
		return fmt.Errorf("failed to delete encounter: %w", err)
	}

	return nil
}

func (s *CombatService) GetParticipant(
	ctx context.Context, participantID uuid.UUID,
) (db.CombatParticipant, error) {
	participant, err := s.queries.GetParticipant(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to get participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) ListParticipants(
	ctx context.Context, combatID uuid.UUID,
) ([]db.CombatParticipant, error) {
	participants, err := s.queries.ListParticipants(ctx, combatID)
	if err != nil {
		return []db.CombatParticipant{}, fmt.Errorf("failed to list participants: %w", err)
	}

	return participants, nil
}

func (s *CombatService) ListActiveParticipants(
	ctx context.Context, combatID uuid.UUID,
) ([]db.CombatParticipant, error) {
	participants, err := s.queries.ListActiveParticipants(ctx, combatID)
	if err != nil {
		return []db.CombatParticipant{}, fmt.Errorf("failed to list active participants: %w", err)
	}

	return participants, nil
}

func (s *CombatService) AddParticipant(
	ctx context.Context, combatID, characterID uuid.UUID, initiative int32,
) (db.CombatParticipant, error) {
	// Get character first
	character, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to get character: %w", err)
	}

	participant, err := s.queries.AddParticipant(ctx, db.AddParticipantParams{
		EncounterID: combatID,
		CharacterID: characterID,
		Initiative: initiative,
		Name: character.Name,
		CurrentHp: character.CurrentHp,
		MaxHp: character.MaxHp,
		TempHp: character.TempHp,
		ArmorClass: character.ArmorClass,
		Speed: character.Speed,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to add participant: %w", err)
	}
	
	return participant, nil
}

func (s *CombatService) UpdateParticipantInitiative(
	ctx context.Context, participantID uuid.UUID, initiative int32, 
) (db.CombatParticipant, error) {
	participant, err := s.queries.UpdateParticipantInitiative(ctx, db.UpdateParticipantInitiativeParams{
		ID: participantID,
		Initiative: initiative,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) UpdateParticipantHP(
	ctx context.Context, participantID uuid.UUID, hp int32, 
) (db.CombatParticipant, error) {
	participant, err := s.queries.UpdateParticipantHP(ctx, db.UpdateParticipantHPParams{
		ID: participantID,
		CurrentHp: hp,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) UpdateParticipantTempHP(
	ctx context.Context, participantID uuid.UUID, hp int32, 
) (db.CombatParticipant, error) {
	participant, err := s.queries.UpdateParticipantTempHP(ctx, db.UpdateParticipantTempHPParams{
		ID: participantID,
		TempHp: hp,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) UpdateParticipantConditions(
	ctx context.Context, participantID uuid.UUID, conditions []string, 
) (db.CombatParticipant, error) {
	participant, err := s.queries.UpdateParticipantConditions(ctx, db.UpdateParticipantConditionsParams{
		ID: participantID,
		Conditions: conditions,
	})
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to update participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) ToggleParticipantConcentration(
	ctx context.Context, participantID uuid.UUID, 
) (db.CombatParticipant, error) {
	participant, err := s.queries.ToggleParticipantConcentration(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to toggle participant concentration: %w", err)
	}

	return participant, nil
}

func (s *CombatService) DeactivateParticipant(
	ctx context.Context, participantID uuid.UUID,
) (db.CombatParticipant, error) {
	participant, err := s.queries.DeactivateParticipant(ctx, participantID)
	if err != nil {
		return db.CombatParticipant{}, fmt.Errorf("failed to deactivate participant: %w", err)
	}

	return participant, nil
}

func (s *CombatService) RemoveParticipant(ctx context.Context, participantID uuid.UUID) error {
	if err := s.queries.RemoveParticipant(ctx, participantID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	return nil
}
