// Code generated by $GOPATH/src/go-common/app/tool/cache/mc. DO NOT EDIT.

/*
  Package dao is a generated mc cache package.
  It is generated from:
  type _mc interface {
		//mc: -key=tokenKey -type=get
		TokenCache(c context.Context, key string) (*model.Report, error)

		//mc: -key=tokenKey -expire=d.mcReportExpire
		AddTokenCache(c context.Context, key string, value *model.Report) error
		//mc: -key=tokenKey -expire=d.mcReportExpire
		AddTokensCache(c context.Context, values map[string]*model.Report) error

		//mc: -key=tokenKey
		DelTokenCache(c context.Context, key string) error
	}
*/

package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/push/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

var _ _mc

// TokenCache get data from mc
func (d *Dao) TokenCache(c context.Context, id string) (res *model.Report, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := tokenKey(id)
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:TokenCache")
		log.Errorv(c, log.KV("TokenCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &model.Report{}
	err = conn.Scan(reply, res)
	if err != nil {
		prom.BusinessErrCount.Incr("mc:TokenCache")
		log.Errorv(c, log.KV("TokenCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddTokenCache Set data to mc
func (d *Dao) AddTokenCache(c context.Context, id string, val *model.Report) (err error) {
	if val == nil {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	key := tokenKey(id)
	item := &memcache.Item{Key: key, Object: val, Expiration: d.mcReportExpire, Flags: memcache.FlagJSON}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:AddTokenCache")
		log.Errorv(c, log.KV("AddTokenCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddTokensCache Set data to mc
func (d *Dao) AddTokensCache(c context.Context, values map[string]*model.Report) (err error) {
	if len(values) == 0 {
		return
	}
	conn := d.mc.Get(c)
	defer conn.Close()
	for id, val := range values {
		key := tokenKey(id)
		item := &memcache.Item{Key: key, Object: val, Expiration: d.mcReportExpire, Flags: memcache.FlagJSON}
		if err = conn.Set(item); err != nil {
			prom.BusinessErrCount.Incr("mc:AddTokensCache")
			log.Errorv(c, log.KV("AddTokensCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
	}
	return
}

// DelTokenCache delete data from mc
func (d *Dao) DelTokenCache(c context.Context, id string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := tokenKey(id)
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:DelTokenCache")
		log.Errorv(c, log.KV("DelTokenCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
