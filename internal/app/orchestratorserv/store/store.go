package store

type Store interface {
	Expression() ExpressionRepository
	Task() TaskRepository
	StackReadyTask() StackReadyTaskRepository
}
