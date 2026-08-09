package service

import (
	"context"
	"fmt"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/google/uuid"
)

type FeatureService struct {
	queries *db.Queries
}

func NewFeatureService(queries *db.Queries) *FeatureService {
	return &FeatureService{queries: queries}
}

func (s *FeatureService) ListFeatures(ctx context.Context, characterID uuid.UUID) ([]db.Feature, error) {
	features, err := s.queries.ListFeatures(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to list features: %w", err)
	}
	return features, nil
}

func (s *FeatureService) CreateFeature(ctx context.Context, params db.CreateFeatureParams) (db.Feature, error) {
	feature, err := s.queries.CreateFeature(ctx, params)
	if err != nil {
		return db.Feature{}, fmt.Errorf("failed to create feature: %w", err)
	}
	return feature, nil
}

func (s *FeatureService) UpdateFeature(ctx context.Context, params db.UpdateFeatureParams) (db.Feature, error) {
	feature, err := s.queries.UpdateFeature(ctx, params)
	if err != nil {
		return db.Feature{}, fmt.Errorf("failed to update feature: %w", err)
	}
	return feature, nil
}

func (s *FeatureService) DeleteFeature(ctx context.Context, id uuid.UUID) error {
	if err := s.queries.DeleteFeature(ctx, id); err != nil {
		return fmt.Errorf("failed to delete feature: %w", err)
	}
	return nil
}
