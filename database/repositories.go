package database

import "github.com/fanky5g/ponzu/internal/domain/interfaces"

type Repositories struct {
	Analytics        interfaces.AnalyticsRepositoryInterface
	Config           interfaces.ConfigRepositoryInterface
	Users            interfaces.UserRepositoryInterface
	Content          interfaces.ContentRepositoryInterface
	CredentialHashes interfaces.CredentialHashRepositoryInterface
	RecoveryKeys     interfaces.RecoveryKeyRepositoryInterface
	Uploads          interfaces.ContentRepositoryInterface
}
