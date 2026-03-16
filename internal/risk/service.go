package risk

type Service interface {
	GetClientsStats() ([]ClientStats, error)
	GetClientProfile(clientID string) (*ClientProfileResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetClientsStats() ([]ClientStats, error) {
	return s.repo.GetClientsStats()
}

func (s *service) GetClientProfile(clientID string) (*ClientProfileResponse, error) {
	return s.repo.GetClientProfile(clientID)
}
