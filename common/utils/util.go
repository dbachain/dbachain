package utils

import (
	"dbachain/common/log"
	"encoding/json"
	"io"
	"io/ioutil"
)

func JSONUnmarshal(r io.Reader, v interface{}) error {
	bin, err := ioutil.ReadAll(r)
	if nil != err {
		log.Error(err.Error())
		return err
	}

	err = json.Unmarshal(bin, v)
	if nil != err {
		log.Errorf("json unmarshal failed, data:%s,err:%s", string(bin), err.Error())
		return err
	}

	return nil
}

// XXX reference the common declaration of this function
func Subspace(prefix []byte) (start, end []byte) {
	end = make([]byte, len(prefix))
	copy(end, prefix)
	end[len(end)-1]++
	return prefix, end
}
