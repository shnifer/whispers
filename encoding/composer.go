package main

import (
	"encoding/binary"
	"github.com/pkg/errors"
)

type composer struct {
	//	data [][]byte
	names    []string
	result   []byte
	compress []byte
}

func NewComposer() *composer {
	return &composer{
		result: make([]byte, 2),
	}
}

func (c *composer) Reset() {
	c.result = c.result[:2]
	c.names = c.names[:0]
}

func (c *composer) Add(name string, data []byte) error {
	if name != "" {
		for i := range c.names {
			if c.names[i] == name {
				return errors.New("already exist")
			}
		}
	}
	c.names = append(c.names, name)

	c.result = append(c.result, byte(len(name)))
	c.result = append(c.result, []byte(name)...)

	var datalen [4]byte
	binary.LittleEndian.PutUint32(datalen[:], uint32(len(data)))
	c.result = append(c.result, datalen[:]...)
	c.result = append(c.result, data...)
	return nil
}

//Encode returns c.compress that should be used before next Encode call
func (c *composer) Encode() []byte {
	binary.LittleEndian.PutUint16(c.result[:2], uint16(len(c.names)))
	//c.compress = snappy.Encode(c.compress[:cap(c.compress)],c.result)
	//return c.compress
	return c.result
}

type decomposer struct {
	decompress []byte
	names      []string
	data       [][]byte
}

func NewDecomposer() *decomposer {
	return &decomposer{}
}

//decompress to internal slice, so v could be disposed after Decode
func (d *decomposer) Decode(v []byte) (err error) {
	if len(v) < 2 {
		return errors.New("can't decode empty")
	}
	//d.decompress, err = snappy.Decode(d.decompress[:cap(d.decompress)], v)
	//if err!=nil{
	//	return err
	//}
	//v = d.decompress
	l := int(binary.LittleEndian.Uint16(v[:2]))
	if len(d.names) < l {
		d.names = make([]string, l)
		d.data = make([][]byte, l)
	} else {
		d.names = d.names[:l]
		for i := range d.data {
			d.data[i] = nil
		}
		d.data = d.data[:l]
	}
	p := 2
	var sl int
	for i := 0; i < l; i++ {
		sl = int(v[p])
		d.names[i] = string(v[p+1 : p+1+sl])
		p += sl + 1
		sl = int(binary.LittleEndian.Uint32(v[p : p+4]))
		d.data[i] = v[p+4 : p+4+sl]
		p += sl + 4
	}
	return nil
}

func (d *decomposer) GetByName(name string) []byte {
	for i := range d.names {
		if d.names[i] == name {
			return d.data[i]
		}
	}
	return nil
}

func (d *decomposer) Count() int {
	return len(d.names)
}

func (d *decomposer) Get(i int) (name string, data []byte) {
	return d.names[i], d.data[i]
}
