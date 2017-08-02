// tdc means toml based dynamic config
package tdc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
)

type Listener func(name string, data []byte, version uint64, exist bool)

type ResourceSolver interface {
	ConfigurationRefByName(name string, listener Listener) ([]byte, error)
}

func NewDynamicToml(solver ResourceSolver, listener Listener) *dynamicToml {
	return &dynamicToml{
		Solver:   solver,
		listener: listener,
	}
}

type dynamicToml struct {
	Solver   ResourceSolver
	listener Listener
}

func (d *dynamicToml) Load(filePath string) ([]byte, error) {
	tomlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return d.Parse(tomlFile)
}
func (d *dynamicToml) Parse(data []byte) (out []byte, err error) {
	output := bytes.NewBuffer(nil)
	lines := bytes.Split(data, []byte{'\n'})
	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		data, err := d.tryExec(line)
		if err != nil {
			return nil, errors.New(strconv.Itoa(i) + "line:" + string(line) + " query resource failed.caused by :" + err.Error())
		}
		output.Write(data)
		output.Write([]byte{'\n'})
	}
	return output.Bytes(), nil
}

// 尝试执行
func (d *dynamicToml) tryExec(data []byte) (out []byte, err error) {
	if data[0] != '#' {
		return data, nil
	}
	words := bytes.Split(data, []byte{' '})

	words = filterEmpty(words)

	if len(words) != 2 || !bytes.Equal(words[0], []byte("####")) {
		return data, nil
	}

	source, err := d.Solver.ConfigurationRefByName(string(words[1]), d.listener)
	if err != nil {
		return data, err
	}
	out = make([]byte, len(data)+len(source)+1)

	copy(out, data)

	out[len(data)] = '\n'

	copy(out[len(data)+1:], source)

	return out, nil
}
func filterEmpty(arr [][]byte) [][]byte {
	out := [][]byte{}
	for _, v := range arr {
		if len(v) != 0 {
			out = append(out, v)
		}
	}
	return out
}
