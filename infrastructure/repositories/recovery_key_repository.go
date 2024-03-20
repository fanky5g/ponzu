package repositories

type RecoveryKeyRepositoryInterface interface {
	SetRecoveryKey(email, key string) error
	GetRecoveryKey(email string) (string, error)
}
