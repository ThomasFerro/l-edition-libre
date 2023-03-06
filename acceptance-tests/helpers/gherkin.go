package helpers

import (
	"context"

	"github.com/cucumber/messages-go/v16"
)

func isAnErrorHandlingScenario(ctx context.Context) bool {
	tags := ctx.Value(TagsKey{}).([]*messages.PickleTag)
	for _, nextTag := range tags {
		if nextTag.Name == "@Error" {
			return true
		}
	}
	return false
}
