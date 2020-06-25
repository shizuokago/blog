package datastore

import (
	"context"
	"strings"

	"cloud.google.com/go/datastore"
)

func init() {
}

func createClient(ctx context.Context) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, "shizuoka-go")
	return client, err
}

func CreateSubTitle(src string) string {

	dst := strings.Replace(src, "\n", "", -1)
	dst = strings.Replace(dst, "*", "", -1)

	if len(dst) > 600 {
		dst = string([]rune(dst)[0:200]) + "..."
	}
	return dst
}
