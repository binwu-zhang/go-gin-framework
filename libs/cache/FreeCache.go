package cache

import (
	"github.com/coocood/freecache"
)

var FreeCache *freecache.Cache

func init (){
	FreeCache = freecache.NewCache(1024 * 1024 * 256)
}


