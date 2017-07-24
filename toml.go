// tdc means toml based dynamic config
package tdc

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"
)

type ResourceSolver interface {
	ConfigurationRefByName(name string) ([]byte, error)
}

func NewDynamicToml(solver ResourceSolver) *dynamicToml {
	return &dynamicToml{
		Solver: solver,
	}
}

type dynamicToml struct {
	Solver ResourceSolver
}

func (d *dynamicToml) Load(filePath string) ([]byte, error) {
	tomlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	tomlFile, err = d.parse(tomlFile)
	if err != nil {
		return nil, err
	}
	return tomlFile, nil
}
func (d *dynamicToml) parse(data []byte) (out []byte, err error) {
	output := bytes.NewBuffer(nil)
	lines := bytes.Split(data, []byte('\n'))
	for i, line := range lines {
		line = bytes.TrimSpace(line)
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

	if len(words) != 2 || bytes.Equal(words[1], []byte("refs")) {
		return data, nil
	}

	source, err := d.Solver.ConfigurationRefByName(string(words[1]))
	if err != nil {
		return data, err
	}
	out = make([]byte, len(data)+len(out)+1)

	copy(out, data)

	out[len(data)] = '\n'

	copy(out[:len(data)+1], source)

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
