package cache

//根据key获取相应的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

//从对应的组中查询缓存值
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
