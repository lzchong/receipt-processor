package receipt

type Service interface {
	Points(id string) int64
	SetPoints() string
}

type serviceImpl struct{}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) Points(id string) int64 {
	return 32
}

func (s *serviceImpl) SetPoints() string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
