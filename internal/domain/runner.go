package domain

type NodeRunner interface {
	Run(node Node, task Task) error
}
