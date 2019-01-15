package health

import (
	"errors"
	"gopkg.in/redis.v3"
	"regexp"
	"strings"
	"time"
)

type RedisCheck struct {
	Name   string
	Type   string
	Client interface{}
}

func (rc RedisCheck) Check() *DependencyInfo {
	var err error
	var t float64
	ver := ""
	sha1 := ""
	state := DependencyState{Status: OK}

	sTime := time.Now()
	switch c := rc.Client.(type) {
	case *redis.Client:
		if s, er1 := c.Info("Server").Result(); er1 == nil {
			ver = getMatch(s, regexp.MustCompile(`redis_version:\s*(\w|\-|\.)+`))
			sha1 = getMatch(s, regexp.MustCompile(`redis_git_sha1:\s*(\w|\-|\.)+`))
		} else {
			err = errors.New("Error querying Redis info: " + er1.Error())
		}
		t = time.Since(sTime).Seconds()
	default:
		err = errors.New("Unknown type of Redis client")
	}

	if err != nil {
		state.Status = CRIT
		state.Details = err.Error()
	}
	return &DependencyInfo{
		Name:         rc.Name,
		Type:         rc.Type,
		Version:      ver,
		Revision:     sha1,
		State:        state,
		ResponseTime: t,
	}
}

func (rc RedisCheck) GetName() string {
	return rc.Name
}

func (rc RedisCheck) GetType() string {
	return rc.Type
}

func getValueFromPair(s, sep string) string {
	if i := strings.Index(s, sep); i >= 0 {
		return s[i+1:]
	}
	return ""
}

func getMatch(s string, reg *regexp.Regexp) string {
	if m := reg.FindString(s); m != "" {
		return getValueFromPair(m, ":")
	}
	return ""
}
