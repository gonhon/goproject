package singleflight

import "sync"

//标识正在请求中或结束的请求
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//用于管理call
type Group struct {
	mutex sync.Mutex
	m     map[string]*call
}

//根据key进行加锁
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mutex.Lock()
	//初始化map
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if cl, ok := g.m[key]; ok {
		g.mutex.Unlock()
		cl.wg.Wait() //正在请求 等待
		return cl.val, cl.err
	}
	cl := new(call)
	cl.wg.Add(1)
	g.m[key] = cl
	g.mutex.Unlock()

	//执行func
	val, err := fn()
	//执行后释放
	cl.wg.Done()

	//执行后移除
	g.mutex.Lock()
	delete(g.m, key)
	g.mutex.Unlock()

	return val, err
}
