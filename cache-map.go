package decorator

type CacheMap interface{
	Store(key, value any) error
	Load(key any) (value any, ok bool, err error)
}