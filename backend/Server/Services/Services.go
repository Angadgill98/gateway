package services



type Services struct{
	AppService *ProxyService
}

func CreateService()*Services{
	return &Services{
		AppService: CreateProxyService(),
	}
}