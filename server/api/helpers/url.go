package helpers

import (
	"context"
	"fmt"
)

func FromUrlParams(ctx context.Context, param string) string {
	return ctx.Value(fmt.Sprintf("URL_PARAM%v", param)).(string)
}
