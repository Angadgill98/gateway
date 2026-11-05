package services



type Services struct{
	ProxyService *ProxyService
}

func CreateService()*Services{
	return &Services{
		ProxyService: CreateProxyService(),
	}
}