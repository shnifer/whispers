package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/snappy"
	"github.com/mailru/easyjson"
	"github.com/shnifer/magellan/commons"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
)

var b = new(bytes.Buffer)

func Test_encode(t *testing.T) {
	ms := prepareMapStr()
	mb := prepareMapSlc()
	ss := prepareSlcStr()
	sb := prepareSlcSlc()

	t.Run("marshal", func(t *testing.T) {
		t.Run("bigData", func(t *testing.T) {
			o1, o2 := prepareBigStruct()
			t.Run("json", func(t *testing.T) {
				b1, _ := json.Marshal(o1)
				b2, _ := json.Marshal(o2)
				t.Log(len(b1) + len(b2))
			})
			t.Run("gob", func(t *testing.T) {
				b1 := marshalGob(o1)
				b2 := marshalGob(o2)
				t.Log(len(b1) + len(b2))
			})
			t.Run("gob full", func(t *testing.T) {
				Data := make(map[string][]byte)
				Data["bigPart1"] = marshalGob(o1)
				Data["bigPart2"] = marshalGob(o2)
				r := marshalGob(Data)
				t.Log(len(r))
			})
		})
		t.Run("small data", func(t *testing.T) {
			o1, o2, o3 := prepareSmallStruct()
			t.Run("json", func(t *testing.T) {
				b1, _ := json.Marshal(o1)
				b2, _ := json.Marshal(o2)
				b3, _ := json.Marshal(o3)
				t.Log(len(b1) + len(b2) + len(b3))
			})
			t.Run("gob", func(t *testing.T) {
				b1 := marshalGob(o1)
				b2 := marshalGob(o2)
				b3 := marshalGob(o3)
				t.Log(len(b1) + len(b2) + len(b3))
			})
			t.Run("gob full", func(t *testing.T) {
				Data := make(map[string][]byte)
				Data["bigPart1"] = marshalGob(o1)
				Data["bigPart2"] = marshalGob(o2)
				Data["bigPart3"] = marshalGob(o3)
				r := marshalGob(Data)
				t.Log(len(r))
			})
		})
	})

	t.Run("ZIP:", func(t *testing.T) {
		t.Run("ms+jsonzip", func(t *testing.T) {
			r, _ := encodeJSONZIP(ms, gzip.BestSpeed)
			t.Log(len(r))
		})
		t.Run("mb+jsonzip", func(t *testing.T) {
			r, _ := encodeJSONZIP(mb, gzip.BestSpeed)
			t.Log(len(r))
		})
		t.Run("ss+jsonzip", func(t *testing.T) {
			r, _ := encodeJSONZIP(ss, gzip.BestSpeed)
			t.Log(len(r))
		})
		t.Run("sb+jsonzip", func(t *testing.T) {
			r, _ := encodeJSONZIP(sb, gzip.BestSpeed)
			t.Log(len(r))
		})
	})
	t.Run("Only JSON:", func(t *testing.T) {
		t.Run("ms", func(t *testing.T) {
			r, _ := encodeJSON(ms)
			log.Println(len(r))
		})
		t.Run("ms", func(t *testing.T) {
			r, _ := encodeJSON(mb)
			log.Println(len(r))
		})
		t.Run("ms", func(t *testing.T) {
			r, _ := encodeJSON(ss)
			log.Println(len(r))
		})
		t.Run("ms", func(t *testing.T) {
			r, _ := encodeJSON(sb)
			log.Println(len(r))
		})
	})

	t.Run("Composer", func(t *testing.T) {
		var composed []byte
		b1, b2 := getFileBytes()
		t.Log("len1 = ", len(b1))
		t.Log("len2 = ", len(b2))
		t.Run("should compose", func(t *testing.T) {
			c := NewComposer()
			c.Add("bigPart1", b1)
			c.Add("bigPart2", b2)
			composed = c.Encode()
			t.Log("start len ", len(b1)+len(b2), "end len", len(composed))
		})
		t.Run("should decode", func(t *testing.T) {
			d := NewDecomposer()
			err := d.Decode(composed)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("names ", d.names)
			require.Equal(t, d.names, []string{"bigPart1", "bigPart2"})
			require.Equal(t, b1, d.data[0])
			require.Equal(t, b2, d.data[1])
		})
	})
}

func Benchmark_encode(b *testing.B) {
	ms := prepareMapStr()
	mb := prepareMapSlc()
	ss := prepareSlcStr()
	sb := prepareSlcSlc()

	ems := EMapOfStr{Map: ms}
	emb := EMapOfSlice{Map: mb}
	ess := ESliceOfStr{Slice: ss}
	esb := ESliceOfSlice{Slice: sb}

	b.Run("MARSHAL", func(b *testing.B) {
		b.Run("big data", func(b *testing.B) {
			o1, o2 := prepareBigStruct()
			e1, e2 := prepareBigStructEasy()
			b.ResetTimer()
			b.Run("json", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = json.Marshal(o1)
					_, _ = json.Marshal(o2)
				}
			})
			b.Run("json EASY", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = easyjson.Marshal(e1)
					_, _ = easyjson.Marshal(e2)
				}
			})
			b.Run("gob", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = marshalGob(o1)
					_ = marshalGob(o2)
				}
			})
			b.Run("gob full", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Data := make(map[string][]byte)
					Data["bigPart1"] = marshalGob(o1)
					Data["bigPart2"] = marshalGob(o2)
					_ = marshalGob(Data)
				}
			})
		})
		b.Run("small data", func(b *testing.B) {
			o1, o2, o3 := prepareSmallStruct()
			e1, e2, e3 := prepareSmallStructEasy()
			b.ResetTimer()
			b.Run("json", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = json.Marshal(o1)
					_, _ = json.Marshal(o2)
					_, _ = json.Marshal(o3)
				}
			})
			b.Run("json easy", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = easyjson.Marshal(e1)
					_, _ = easyjson.Marshal(e2)
					_, _ = easyjson.Marshal(e3)
				}
			})
			b.Run("gob", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = marshalGob(o1)
					_ = marshalGob(o2)
					_ = marshalGob(o3)
				}
			})
			b.Run("gob full", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Data := make(map[string][]byte)
					Data["bigPart1"] = marshalGob(o1)
					Data["bigPart2"] = marshalGob(o2)
					Data["bigPart3"] = marshalGob(o3)
					_ = marshalGob(Data)
				}
			})
		})
	})

	b.Run("ONLY json", func(b *testing.B) {
		b.Run("ms+json", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSON(ms)
			}
		})
		b.Run("mb+json", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSON(mb)
			}
		})
		b.Run("ss+json", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSON(ss)
			}
		})
		b.Run("sb+json", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSON(sb)
			}
		})
	})
	b.Run("ONLY json EASY", func(b *testing.B) {
		b.Run("ms+easyjson", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = encodeEasyJSON(ems)
			}
		})
		b.Run("mb+easyjson", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = encodeEasyJSON(emb)
			}
		})
		b.Run("ss+easyjson", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = encodeEasyJSON(ess)
			}
		})
		b.Run("sb+easyjson", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = encodeEasyJSON(esb)
			}
		})
	})
	b.Run("JSON+ZIP", func(b *testing.B) {
		b.Run("ms+jsonzip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSONZIP(ms, gzip.BestSpeed)
			}
		})
		b.Run("mb+jsonzip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSONZIP(mb, gzip.BestSpeed)
			}
		})
		b.Run("ss+jsonzip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSONZIP(ss, gzip.BestSpeed)
			}
		})
		b.Run("sb+jsonzip", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = encodeJSONZIP(sb, gzip.BestSpeed)
			}
		})
	})
}

func Benchmark_composer(b *testing.B) {
	composer := NewComposer()
	b1, b2 := getFileBytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composer.Reset()
		composer.Add("BigPart1", b1)
		composer.Add("BigPart2", b2)
		_ = composer.Encode()
	}
}

func Benchmark_decomposer(b *testing.B) {
	composer := NewComposer()
	b1, b2 := getFileBytes()
	composer.Add("BigPart1", b1)
	composer.Add("BigPart2", b2)
	composed := composer.Encode()
	d := NewDecomposer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Decode(composed)
	}
}

func BenchmarkEasy_compose(b *testing.B) {
	b.Run("BigData", func(b *testing.B) {
		s1, s2 := prepareBigStructEasy()
		c := NewComposer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Reset()
			c.Add("BigPart1", encodeEasyJSON(s1))
			c.Add("BigPart2", encodeEasyJSON(s2))
			c.Encode()
		}
	})
	b.Run("SmallData", func(b *testing.B) {
		s1, s2, s3 := prepareSmallStructEasy()
		c := NewComposer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.Reset()
			c.Add("BigPart1", encodeEasyJSON(s1))
			c.Add("BigPart2", encodeEasyJSON(s2))
			c.Add("BigPart3", encodeEasyJSON(s3))
			c.Encode()
		}
	})
}

func BenchmarkCollapse(b *testing.B) {
	b1, b2 := prepareBigStructEasy()
	c := NewComposer()
	c.Add("BigPart1", encodeEasyJSON(b1))
	c.Add("BigPart2", encodeEasyJSON(b2))
	bigPayload := c.Encode()
	s1, s2, s3 := prepareSmallStructEasy()
	c.Reset()
	c.Add("BigPart1", encodeEasyJSON(s1))
	c.Add("BigPart2", encodeEasyJSON(s2))
	c.Add("BigPart3", encodeEasyJSON(s3))
	smallPayload := c.Encode()
	_ = smallPayload
	b.ResetTimer()
	b.Run("gzip", func(b *testing.B) {
		b.Run("bigload", func(b *testing.B) {
			buf := &bytes.Buffer{}

			for i := 0; i < b.N; i++ {
				buf.Reset()
				enc, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
				enc.Write(bigPayload)
				enc.Close()
			}
			b.Log("size: ", buf.Len())
		})
		b.Run("smallload", func(b *testing.B) {
			buf := &bytes.Buffer{}
			for i := 0; i < b.N; i++ {
				buf.Reset()
				enc, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
				enc.Write(smallPayload)
				enc.Close()
			}
			b.Log("size: ", buf.Len())
		})
	})
	b.Run("snappy", func(b *testing.B) {
		var dest []byte
		b.Run("bigload", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dest = snappy.Encode(dest[:cap(dest)], bigPayload)
			}
			b.Log("size: ", len(bigPayload), "->", len(dest))
		})
		b.Run("smallload", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dest = snappy.Encode(dest[:cap(dest)], smallPayload)
			}
			b.Log("size: ", len(smallPayload), "->", len(dest))
		})
	})
	b.Run("snappy parts", func(b *testing.B) {
		var dest []byte
		b.Run("bigload", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dest = snappy.Encode(dest[:cap(dest)], bigPayload)
			}
			b.Log("size: ", len(bigPayload), "->", len(dest))
		})
		b.Run("smallload", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dest = snappy.Encode(dest[:cap(dest)], smallPayload)
			}
			b.Log("size: ", len(smallPayload), "->", len(dest))
		})
	})
}

func encodeJSONZIP(v interface{}, level int) (result []byte, err error) {
	*b = bytes.Buffer{}
	enc, _ := gzip.NewWriterLevel(b, level)
	err = json.NewEncoder(enc).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encodeJSON(v interface{}) (result []byte, err error) {
	return json.Marshal(v)
}

func encodeEasyJSON(v easyjson.Marshaler) (result []byte) {
	result, _ = easyjson.Marshal(v)
	return
}

func prepareMapStr() MapOfStr {
	b1, b2 := getFileBytes()
	Big := make(MapOfStr)
	Big["bigPart1"] = string(b1)
	Big["bigPart2"] = string(b2)
	return Big
}
func prepareMapSlc() MapOfSlice {
	b1, b2 := getFileBytes()
	Big := make(MapOfSlice)
	Big["bigPart1"] = b1
	Big["bigPart2"] = b2
	return Big
}
func prepareSlcStr() SliceOfStr {
	b1, b2 := getFileBytes()
	Big := make(SliceOfStr, 2)
	Big[0] = struct {
		Name string
		Data string
	}{Name: "bigPart1", Data: string(b1)}
	Big[1] = struct {
		Name string
		Data string
	}{Name: "bigPart2", Data: string(b2)}
	return Big
}
func prepareSlcSlc() SliceOfSlice {
	b1, b2 := getFileBytes()
	Big := make(SliceOfSlice, 2)
	Big[0] = struct {
		Name string
		Data []byte
	}{Name: "bigPart1", Data: b1}
	Big[1] = struct {
		Name string
		Data []byte
	}{Name: "bigPart2", Data: b2}
	return Big
}
func getFileBytes() ([]byte, []byte) {
	b1, err := ioutil.ReadFile("galaxy_warp.json")
	if err != nil {
		panic(err)
	}
	b2, err := ioutil.ReadFile("galaxy_IV5.json")
	if err != nil {
		panic(err)
	}
	return b1, b2
}

func prepareBigStruct() (s1 commons.Galaxy, s2 commons.Galaxy) {
	b1, b2 := getFileBytes()
	err := json.Unmarshal(b1, &s1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b2, &s2)
	if err != nil {
		panic(err)
	}
	return s1, s2
}

func prepareSmallStruct() (s1 commons.PilotData, s2 commons.NaviData, s3 commons.EngiData) {
	return
}

func prepareBigStructEasy() (s1 Galaxy, s2 Galaxy) {
	b1, b2 := getFileBytes()
	err := json.Unmarshal(b1, &s1)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b2, &s2)
	if err != nil {
		panic(err)
	}
	return s1, s2
}

func prepareSmallStructEasy() (s1 PilotData, s2 NaviData, s3 EngiData) {
	return
}
func marshalGob(v interface{}) []byte {
	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	enc.Encode(v)
	return b.Bytes()
}
