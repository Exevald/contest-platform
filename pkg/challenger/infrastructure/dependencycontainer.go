package infrastructure

func NewDependencyContainer(
	appID string,
) *DependencyContainer {
	return &DependencyContainer{}
}

type DependencyContainer struct {
}
