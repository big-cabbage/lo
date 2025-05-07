package cache

type IStorage interface {
	Get(key string) (any, bool)
	Set(key string, value any) error
}
