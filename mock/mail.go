package mock

type MailClientMockImpl struct{}

func (m MailClientMockImpl) Send(_ []string, _ string) error { return nil }
