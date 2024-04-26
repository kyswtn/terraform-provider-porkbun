package consts

func GetDefaultNameservers() []string {
	// Based on the data from https://kb.porkbun.com/article/63-how-to-switch-to-porkbuns-nameservers.
	return []string{
		"maceio.ns.porkbun.com",
		"curitiba.ns.porkbun.com",
		"salvador.ns.porkbun.com",
		"fortaleza.ns.porkbun.com",
	}
}

const (
	ErrUnknownConfigurationValue         = "Unknown configuration value"
	ErrInvalidConfigurationValue         = "Invalid configuration value"
	ErrUnexpectedDataSourceConfigureType = "Unexpected data source configure type"
	ErrUnexpectedResourceConfigureType   = "Unexpected resource configure type"
)
