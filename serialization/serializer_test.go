package serialization

import (
	"testing"

	"github.com/merenguessss/dracarys/serialization/testdata"
)

var pbHello = &testdata.PHello{
	Msg: "123456789",
}

var rfHello = &testdata.RHello{
	Msg: "123456789",
}

func TestJsonMarshalAndUnmarshal(t *testing.T) {
	testReflectDataMarshalAndUnmarshal(Get(Json), t)
}

func TestMsgPackMarshalAndUnmarshal(t *testing.T) {
	testReflectDataMarshalAndUnmarshal(Get(Msgpack), t)
}

func TestGencodeMarshalAndUnmarshal(t *testing.T) {
	testReflectDataMarshalAndUnmarshal(Get(Gencode), t)
}

func testReflectDataMarshalAndUnmarshal(serializer Serialization, t testing.TB) {
	b, err := serializer.Marshal(rfHello)
	if err != nil {
		t.Error(err)
	}

	res := &testdata.RHello{}
	err = serializer.Unmarshal(b, res)
	if err != nil {
		t.Error(err)
	}
	if res.Msg != rfHello.Msg {
		t.Error("serializer error")
	}
}

func TestProtoMarshalAndUnmarshal(t *testing.T) {
	serializer := Get(Proto)
	b, err := serializer.Marshal(pbHello)
	if err != nil {
		t.Error(err)
	}

	res := &testdata.PHello{}
	err = serializer.Unmarshal(b, res)
	if err != nil {
		t.Error(err)
	}
	if res.Msg != pbHello.Msg {
		t.Error("proto error")
	}
}

var num = 100000

func BenchmarkJsonSerializer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < num; i++ {
		testReflectDataMarshalAndUnmarshal(Get(Json), b)
	}
}

func BenchmarkMsgPackSerializer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < num; i++ {
		testReflectDataMarshalAndUnmarshal(Get(Msgpack), b)
	}
}

func BenchmarkGencodeSerializer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < num; i++ {
		testReflectDataMarshalAndUnmarshal(Get(Gencode), b)
	}
}

func BenchmarkProtoSerializer(b *testing.B) {
	b.ResetTimer()
	serializer := Get(Proto)
	for i := 0; i < num; i++ {
		bytes, err := serializer.Marshal(pbHello)
		if err != nil {
			b.Error(err)
		}

		res := &testdata.PHello{}
		err = serializer.Unmarshal(bytes, res)
		if err != nil {
			b.Error(err)
		}
		if res.Msg != pbHello.Msg {
			b.Error("proto error")
		}
	}
}
