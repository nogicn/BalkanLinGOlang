package server

func (s *FiberServer) RegisterFiberRoutes() {

	s.RegisterDictionaryRouter()
	s.RegisterIndexRouter()
	s.RegisterLocaleRouter()
	s.RegisterUserRouter()
	s.RegisterWordRouter()
}
