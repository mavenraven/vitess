package provision

type factory = func(map[string]string) (Provisioner, error)
var (
	provisioners = make(map[string]factory)
)

func NewProvisioner(provisionerType string, config map[string]string) (Provisioner, error) {
	f, ok := provisioners[provisionerType]
	if !ok {
		return nil, ErrInvalidProvisionerType
	}
	return f(config)
}
