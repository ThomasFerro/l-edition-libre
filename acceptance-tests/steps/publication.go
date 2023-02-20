package steps

import (
	"acceptance-tests/helpers"
	"context"
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api"
)

func isPublished(ctx context.Context, publicationID string) (context.Context, error) {
	ctx, publication, err := getPublicationStatus(ctx, publicationID)
	if err != nil {
		return ctx, fmt.Errorf("cannot check publication's status: %v", err)
	}

	if publication.Status != "Available" {
		return ctx, fmt.Errorf("publication should be available instead of %v", publication.Status)
	}
	return ctx, nil
}

func getPublicationStatus(ctx context.Context, publicationID string) (context.Context, api.PublicationDto, error) {
	url := fmt.Sprintf("http://localhost:8080/api/publications/%v", publicationID)
	var publication api.PublicationDto
	ctx, err := helpers.Call(ctx, url, http.MethodGet, nil, &publication)
	if err != nil {
		return ctx, api.PublicationDto{}, fmt.Errorf("unable to get publication's status: %v", err)
	}
	return ctx, publication, nil
}
