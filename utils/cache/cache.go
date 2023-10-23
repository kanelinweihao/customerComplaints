package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.lwh.com/linweihao/customerComplaints/config/env"
	"go.lwh.com/linweihao/customerComplaints/utils/err"
	"time"
	// "go.lwh.com/linweihao/customerComplaints/utils/time"
	"go.lwh.com/linweihao/customerComplaints/factory/factoryOfSSH"
	"go.lwh.com/linweihao/customerComplaints/utils/rfl"
	"go.lwh.com/linweihao/customerComplaints/utils/ssh"
)

var MaxRetries int = 0
var TimeoutDial time.Duration = time.Second * time.Duration(5)
var TimeoutRead time.Duration = time.Second * time.Duration(1)
var TimeoutWrite time.Duration = time.Second * time.Duration(5)

type EntityCache struct {
	CacheRedis *redis.Client
	EntitySSH  *ssh.EntitySSH
}
type EntityConfigRedis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

/*
Init
*/

func InitCache() (entityCache *EntityCache) {
	cacheRedis := &redis.Client{}
	entitySSH := getEntitySSH()
	r := getEntityConfigRedis()
	cacheRedis = initCacheRedis(r, entitySSH)
	// time.ShowTimeAndMsg("PING")
	// resRedisPing, errRedisPing := cacheRedis.Ping().Result()
	_, errRedisPing := cacheRedis.Ping().Result()
	err.ErrCheck(errRedisPing)
	// time.ShowTimeAndMsg(resRedisPing)
	entityCache = &EntityCache{
		CacheRedis: cacheRedis,
		EntitySSH:  entitySSH,
	}
	return entityCache
}

func getEntitySSH() (entitySSH *ssh.EntitySSH) {
	entitySSH = nil
	isNeedSSH := isNeedSSH()
	if isNeedSSH {
		// entitySSH = ssh.InitSSHForRedis()
		entitySSH = factoryOfSSH.MakeEntityOfSSHForRedis()
	}
	return entitySSH
}

func isNeedSSH() (isNeedSSH bool) {
	isNeedSSH = env.IsNeedSSH()
	return isNeedSSH
}

func getEntityConfigRedis() (r *EntityConfigRedis) {
	paramsRedis := env.GetParamsRedis()
	r = &EntityConfigRedis{}
	rfl.ToEntityFromAttr(paramsRedis, r)
	return r
}

func initCacheRedis(r *EntityConfigRedis, entitySSH *ssh.EntitySSH) (cacheRedis *redis.Client) {
	redisHost := r.Host
	redisPort := r.Port
	redisPassword := r.Password
	redisDb := r.DB
	redisAddr := fmt.Sprintf(
		"%s:%s",
		redisHost,
		redisPort)
	maxRetries := MaxRetries
	timeoutDial := TimeoutDial
	timeoutRead := TimeoutRead
	timeoutWrite := TimeoutWrite
	entityRedisOptions := redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           redisDb,
		MaxRetries:   maxRetries,
		DialTimeout:  timeoutDial,
		ReadTimeout:  timeoutRead,
		WriteTimeout: timeoutWrite,
	}
	isNeedSSH := isNeedSSH()
	if isNeedSSH {
		entitySSH.SetAddress(redisAddr)
		funcDialerRedis := entitySSH.DialForRedis
		entityRedisOptions.Dialer = funcDialerRedis
	}
	// fmt.Println(entityRedisOptions)
	// rfl.ShowType(entityRedisOptions)
	cacheRedis = redis.NewClient(&entityRedisOptions)
	// fmt.Println(cacheRedis)
	// time.ShowTimeAndMsg("Redis connect success")
	return cacheRedis
}

/*
Exec
*/

func (self *EntityCache) CloseCache() {
	cacheRedis := self.CacheRedis
	if cacheRedis != nil {
		cacheRedis.Close()
	}
	self.CacheRedis = nil
	// time.ShowTimeAndMsg("Cache close success")
	isNeedSSH := env.IsNeedSSH()
	if isNeedSSH {
		entitySSH := self.EntitySSH
		if entitySSH != nil {
			entitySSH.CloseSSH()
			// time.ShowTimeAndMsg("SSH close success")
		}
	}
	self.EntitySSH = nil
	return
}

func (self *EntityCache) SetToCache(cacheKey string, cacheValue interface{}, ttl time.Duration) {
	cacheRedis := self.CacheRedis
	errRedisSet := cacheRedis.Set(cacheKey, cacheValue, ttl).Err()
	err.ErrCheck(errRedisSet)
	return
}
