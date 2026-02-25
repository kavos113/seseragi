package domain

type NodeRunner interface {
	Run(node Node) error
}
