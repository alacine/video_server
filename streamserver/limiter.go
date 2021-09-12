// 限制连接数量
// https://zhuanlan.zhihu.com/p/64124008
package main

import "log"

// ConnLimiter 限流器
type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

// NewConnLimiter ...
func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

// GetConnLimiter ...
func (cl *ConnLimiter) GetConnLimiter() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation.")
		return false
	}
	cl.bucket <- 1
	return true
}

// ReleaseConn ...
func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("New connection coming: %d", c)
}
