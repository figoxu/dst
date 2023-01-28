package dst

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

func Dup(err error) bool {
	if err == nil {
		return false
	}
	errStr := []string{"Error 1062", "Duplicate entry", "duplicate key value"}
	for _, s := range errStr {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}

func Chk(err error, entries ...*logrus.Entry) {
	if err != nil {
		Log(entries...).WithErr(err).Entry().Panicln(`panic `, err.Error())
	}
}

type Lg struct {
	sync.Mutex
	m  map[string]interface{}
	rs *logrus.Entry
}

func Log(entry ...*logrus.Entry) *Lg {
	rs := logrus.WithFields(logrus.Fields{})
	if len(entry) > 0 {
		rs = entry[0]
	}
	return &Lg{
		rs: rs,
		m:  map[string]interface{}{},
	}
}

func (p *Lg) WithErr(err error) *Lg {
	msg := fmt.Sprintf("%+v", err)
	p.rs = p.rs.WithField("err", msg)
	return p
}

func (p *Lg) Entry() *logrus.Entry {
	return p.rs
}
