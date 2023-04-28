package discovery

import (
	"sync"
	"sync/atomic"
)

// A Subscriber is used to subscribe the given key on an etcd cluster.
type Subscriber struct {
	endpoints []string
	exclusive bool
	container *container
}

type container struct {
	//独占 如果是独占模式那么 kv是 强制1:1的 如果不是那么 多个 key 可能有同一个 value
	exclusive bool
	values    map[string][]string
	mapping   map[string]string
	//快照
	snapshot atomic.Value
	//是否被写入了
	dirty     atomic.Bool
	listeners []func()
	lock      sync.Mutex
}

func newContainer(exclusive bool) *container {
	return &container{
		exclusive: exclusive,
		values:    make(map[string][]string),
		mapping:   make(map[string]string),
		dirty:     atomic.Bool{},
	}
}

func (c *container) OnAdd(kv KV) {
	c.addKv(kv.Key, kv.Val)
	c.notifyChange()
}

func (c *container) OnDelete(kv KV) {
	c.removeKey(kv.Key)
	c.notifyChange()
}

// addKv adds the kv, returns if there are already other keys associate with the value
func (c *container) addKv(key, value string) ([]string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Store(true)
	keys := c.values[value]
	previous := append([]string(nil), keys...)
	//values 里面有值了
	early := len(keys) > 0
	//且容器是独占模式的
	if c.exclusive && early {
		for _, each := range keys {
			//把所有的 key 都删掉
			c.doRemoveKey(each)
		}
	}
	c.values[value] = append(c.values[value], key)
	c.mapping[key] = value

	if early {
		return previous, true
	}

	return nil, false
}

func (c *container) doRemoveKey(key string) {
	server, ok := c.mapping[key]
	if !ok {
		return
	}

	delete(c.mapping, key)
	keys := c.values[server]
	remain := keys[:0]

	for _, k := range keys {
		if k != key {
			remain = append(remain, k)
		}
	}

	if len(remain) > 0 {
		c.values[server] = remain
	} else {
		delete(c.values, server)
	}
}

func (c *container) removeKey(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Store(true)
	c.doRemoveKey(key)
}

func (c *container) notifyChange() {
	c.lock.Lock()
	listeners := append(([]func())(nil), c.listeners...)
	c.lock.Unlock()

	for _, listener := range listeners {
		listener()
	}
}
