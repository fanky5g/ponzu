package root

type ViewModel struct {
	PublicPath string
	AppName    string
	Logo       string
}

func NewRootViewModel(appNameProvider AppNameProvider, publicPath string) (*ViewModel, error) {
	appName, err := appNameProvider.GetAppName()
	if err != nil {
		return nil, err
	}

	return &ViewModel{
		PublicPath: publicPath,
		AppName:    appName,
		Logo:       appName,
	}, nil
}
