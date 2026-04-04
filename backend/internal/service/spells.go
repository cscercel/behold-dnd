package service

import (
	"context"
	"fmt"

    "github.com/cscercel/beyond-dnd/internal/db"
    "github.com/google/uuid"
)


type SpellService struct {
	queries *db.Queries
}

func NewSpellService(queries *db.Queries) *SpellService {
	return &SpellService{queries: queries}
}

func (s *SpellService) UseSlot(ctx context.Context, characterID uuid.UUID, level int32) (db.SpellSlot, error) {
	slot, err := s.queries.UseSpellSlot(ctx, db.UseSpellSlotParams{
		CharacterID: characterID,
		SpellLevel:	level,
	})
	if err != nil {
		return db.SpellSlot{}, fmt.Errorf("no spell slots remaining at level %d", level)
	}

	return slot, nil
}

func (s *SpellService) LongRestSlots(ctx context.Context, characterID uuid.UUID) error {
	return s.queries.ResetSpellSlots(ctx, characterID)
}
