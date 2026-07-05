package service

import (
	"context"
	"fmt"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

type InventoryService struct {
	queries *db.Queries
}

func NewInventoryService(queries *db.Queries) *InventoryService {
	return &InventoryService{queries: queries}
}

func (s *InventoryService) ListItems(ctx context.Context, characterID uuid.UUID) ([]db.InventoryItem, error) {
	items, err := s.queries.ListInventoryItems(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	return items, nil
}

func (s *InventoryService) CreateItem(ctx context.Context, params db.CreateInventoryItemParams) (db.InventoryItem, error) {
	item, err := s.queries.CreateInventoryItem(ctx, params)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("failed to create inventory item: %w", err)
	}
	return item, nil
}

func (s *InventoryService) UpdateItem(ctx context.Context, params db.UpdateInventoryItemParams) (db.InventoryItem, error) {
	item, err := s.queries.UpdateInventoryItem(ctx, params)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("failed to update inventory item: %w", err)
	}
	return item, nil
}

func (s *InventoryService) DeleteItem(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteInventoryItem(ctx, id); err != nil {
		return fmt.Errorf("failed to delete inventory item: %w", err)
	}
	return nil
}

func (s *InventoryService) AttuneItem(ctx context.Context, characterID uuid.UUID, itemID uuid.UUID) (db.InventoryItem, error) {
	char, err := s.queries.GetCharacter(ctx, characterID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("character not found: %w", err)
	}

	count, err := s.queries.CountAttunedItems(ctx, characterID)
	if err != nil {
		return db.InventoryItem{}, fmt.Errorf("failed to count attuned items: %w", err)
	}
	if count >= int64(char.AttunementSlots) {
		return db.InventoryItem{}, fmt.Errorf("attunement limit reached (%d/%d)", count, char.AttunementSlots)
	}

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

	attuned := true
	return s.queries.UpdateInventoryItem(ctx, db.UpdateInventoryItemParams{
		ID:        itemID,
		IsAttuned: &attuned,
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

	attuned := false
	return s.queries.UpdateInventoryItem(ctx, db.UpdateInventoryItemParams{
		ID:        itemID,
		IsAttuned: &attuned,
	})
}
