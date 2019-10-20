/**
 * @Author: wei.tan
 * @Description:
 * @File:  lru_test
 * @Version: 1.0.0
 * @Date: 2019-10-18 21:53
 */

package lru

import (
	"fmt"
	"testing"
)

func TestLRUCache_Get(t *testing.T) {
	lru := NewLRUCache(2)
	lru.Put(1, 2)
	lru.Put(2, 4)
	lru.Put(3, 6)
	v := lru.Get(3)
	fmt.Println(v)
}