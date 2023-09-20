package lb_factory

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	// Update 后期服务发现补充
	//Update()
}
