package mocks

//go:generate mockgen -destination=domain_mock.go -package=mocks github.com/Vla8islav/gophemart/internal/domain GophermartRepository,GophemartService,GophermartAccrualClient
