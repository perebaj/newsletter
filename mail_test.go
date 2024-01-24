package newsletter

import (
	"context"
	"testing"
)

type MailClientMockImpl struct{}

func (m MailClientMockImpl) Send(_ []string, _ string) error { return nil }

func TestEmailTrigger(t *testing.T) {
	ctx := context.Background()
	s := NewStorageMock()
	e := MailClientMockImpl{}
	err := EmailTrigger(ctx, s, e)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
