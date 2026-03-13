package rules

type Service interface {
	GetAllRules() ([]RuleResponse, error)
	CreateRule(req CreateRuleRequest) (RuleResponse, error)
	UpdateRule(id uint, req UpdateRuleRequest) (RuleResponse, error)
	DeleteRule(id uint) error
	ToggleRule(id uint) (RuleResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllRules() ([]RuleResponse, error) {
	rules, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]RuleResponse, 0, len(rules))
	for _, r := range rules {
		res = append(res, NewRuleResponse(r))
	}
	return res, nil
}

func (s *service) CreateRule(req CreateRuleRequest) (RuleResponse, error) {
	rule := RiskRule{
		Name:              req.Name,
		Description:       req.Description,
		EventType:         req.EventType,
		ConditionField:    req.ConditionField,
		ConditionOperator: req.ConditionOperator,
		ConditionValue:    req.ConditionValue,
		Score:             req.Score,
		Enabled:           req.Enabled,
	}

	if err := s.repo.Create(&rule); err != nil {
		return RuleResponse{}, err
	}

	return NewRuleResponse(rule), nil
}

func (s *service) UpdateRule(id uint, req UpdateRuleRequest) (RuleResponse, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return RuleResponse{}, err
	}

	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = *req.Description
	}
	if req.EventType != nil {
		rule.EventType = *req.EventType
	}
	if req.ConditionField != nil {
		rule.ConditionField = *req.ConditionField
	}
	if req.ConditionOperator != nil {
		rule.ConditionOperator = *req.ConditionOperator
	}
	if req.ConditionValue != nil {
		rule.ConditionValue = *req.ConditionValue
	}
	if req.Score != nil {
		rule.Score = *req.Score
	}
	if req.Enabled != nil {
		rule.Enabled = *req.Enabled
	}

	if err := s.repo.Update(rule); err != nil {
		return RuleResponse{}, err
	}

	return NewRuleResponse(*rule), nil
}

func (s *service) DeleteRule(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) ToggleRule(id uint) (RuleResponse, error) {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return RuleResponse{}, err
	}

	rule.Enabled = !rule.Enabled

	if err := s.repo.Update(rule); err != nil {
		return RuleResponse{}, err
	}

	return NewRuleResponse(*rule), nil
}
