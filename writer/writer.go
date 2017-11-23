package writer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"os"
)

var stdout = "/dev/stdout"

func BinWriter(info interface{}) {
	fp, _ := os.Create(stdout)
	defer fp.Close()
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, info)
	fp.Write(buf.Bytes())
	fp.Sync()
}

func TextWriter(info interface{}) {
	c, _ := json.Marshal(info)
	ioutil.WriteFile(stdout, c, os.ModeAppend)
	ioutil.WriteFile(stdout, []byte("\n"), os.ModeAppend)
}
