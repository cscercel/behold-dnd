package service

import (
	"context"
	"fmt"

    "github.com/cscercel/beyond-dnd/internal/db"
    "github.com/google/uuid"
)


type InventoryService struct {
	queries	*db.Queries
}

func NewInventoryService(queries *db.Queries) *InventoryService {
	return &InventoryService{queries: queries}
}

func (s *InventoryService) AttuneItem(ctx context.Context, characterID uuid.UUID, itemID uuid.UUID) (db.InventoryItem, error) {
	char, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("character not found: %w", err)
	}

	// Count currently attuned items
	count, err := s.queries.CountAttunedItems(ctx, characterID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("failed to count attuned items: %w", err)
	}

	if count >= int64(char.AttunementSlots) {
		return db.InventoryItem{}, fmt.Errorf("attunement limit reached (%d/%d)", count, char.AttunementSlots)
	}

	// Make sure item belongs to character
	item, err := s.queries.GetInventoryItem(ctx, itemID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("item not found: %w", err)
	}
	if item.CharacterID != characterID {
		return db.InventoryItem{}, fmt.Errorf("item does not belong to this character")
	}
	if !item.RequiresAttunement {
		return db.InventoryItem{}, fmt.Errorf("item does not require attunement")
	}

	return s.queries.UpdateInventoryItem(ctx, db.UpdateInventoryItemParams{
		ID:	itemID,
		Name: item.Name,
		Quantity: item.Quantity,
		Weight:	item.Weight,
		Description: item.Description,
		IsEquipped:	item.IsEquipped,
		RequiresAttunement: item.RequiresAttunement,
		isAttuned:	true,
	})
}

func (s *InventoryService) UnattuneItem(ctx context.Context, characterID uuid.UUID, itemID uuid.UUID) (db.InventoryItem, error) {
	item, err := s.queries.GetInventoryItem(ctx, itemID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("item not found: %w", err)
	}
	if item.CharacterID != characterID {
		return db.InventoryItem{}, fmt.Errorf("item does not belong to this character")
	}
	if !item.RequiresAttunement {
		return db.InventoryItem{}, fmt.Errorf("item does not require attunement")
	}

	return s.queries.UpdateInventoryItem(ctx, db.UpdateInventoryItemParams{
		ID:	itemID,
		Name: item.Name,
		Quantity: item.Quantity,
		Weight:	item.Weight,
		Description: item.Description,
		IsEquipped:	item.IsEquipped,
		RequiresAttunement: item.RequiresAttunement,
		isAttuned:	false,
	})
}
