package provider

type Provider interface {
	//组件注册
	Register(args ...interface{}) error

	//输出所有的注册实例
	Provides() []string

	//相关组件的关闭
	Close() error
}
