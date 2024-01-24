package newsletter

import (
	"context"
	"testing"

	"github.com/perebaj/newsletter/mock"
)

func TestEmailTrigger(t *testing.T) {
	ctx := context.Background()
	s := mock.NewStorageMock()
	e := mock.MailClientMockImpl{}
	err := EmailTrigger(ctx, s, e)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
