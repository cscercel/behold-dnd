package service

import (
	"context"
	"fmt"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

type SpellService struct {
	queries *db.Queries
}

func NewSpellService(queries *db.Queries) *SpellService {
	return &SpellService{queries: queries}
}

func (s *SpellService) ListSpells(ctx context.Context, characterID uuid.UUID) ([]db.Spell, error) {
	spells, err := s.queries.ListSpells(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list spells: %w", err)
	}
	return spells, nil
}

func (s *SpellService) CreateSpell(ctx context.Context, params db.CreateSpellParams) (db.Spell, error) {
	spell, err := s.queries.CreateSpell(ctx, params)
	if err != nil {
		return db.Spell{}, fmt.Errorf("failed to create spell: %w", err)
	}
	return spell, nil
}

func (s *SpellService) UpdateSpell(ctx context.Context, params db.UpdateSpellParams) (db.Spell, error) {
	spell, err := s.queries.UpdateSpell(ctx, params)
	if err != nil {
		return db.Spell{}, fmt.Errorf("failed to update spell: %w", err)
	}
	return spell, nil
}

func (s *SpellService) DeleteSpell(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteSpell(ctx, id); err != nil {
		return fmt.Errorf("failed to delete spell: %w", err)
	}
	return nil
}

func (s *SpellService) ToggleSpellPrepared(ctx context.Context, id uuid.UUID) (db.Spell, error) {
	spell, err := s.queries.ToggleSpellPrepared(ctx, id)
	if err != nil {
		return db.Spell{}, fmt.Errorf("failed to toggle spell preparation: %w", err)
	}
	return spell, nil
}

func (s *SpellService) ListSpellSlots(ctx context.Context, characterID uuid.UUID) ([]db.SpellSlot, error) {
	slots, err := s.queries.ListSpellSlots(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list spell slots: %w", err)
	}
	return slots, nil
}

func (s *SpellService) UpsertSpellSlot(ctx context.Context, params db.UpsertSpellSlotParams) (db.SpellSlot, error) {
	slot, err := s.queries.UpsertSpellSlot(ctx, params)
	if err != nil {
		return db.SpellSlot{}, fmt.Errorf("failed to upsert spell slot: %w", err)
	}
	return slot, nil
}

func (s *SpellService) UseSpellSlot(ctx context.Context, params db.UseSpellSlotParams) (db.SpellSlot, error) {
	slot, err := s.queries.UseSpellSlot(ctx, params)
	if err != nil {
		return db.SpellSlot{}, fmt.Errorf("no spell slots remaining at level %d: %w", params.SpellLevel, err)
	}
	return slot, nil
}

func (s *SpellService) LongRestSlots(ctx context.Context, characterID uuid.UUID) error {
	return s.queries.ResetSpellSlots(ctx, characterID)
}
