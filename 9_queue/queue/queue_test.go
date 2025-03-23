package queue_test

import (
	"bytes"
	"encoding/json"
	"protohackers/9_queue/queue"
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	in := "{\"request\":\"get\",\"queues\":[\"queue1\",\"queue2\"],\"wait\":true}"
	exp := queue.GetRequest{
		Request: "get",
		Queues:  []string{"queue1", "queue2"},
		Wait:    true,
	}

	var out queue.GetRequest

	b := bytes.NewBufferString(in)
	dec := json.NewDecoder(b)

	err := dec.Decode(&out)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(exp, out) {
		t.Fatalf("decode fail exp %v got %v", exp, out)
	}
}

func TestEncode(t *testing.T) {
	in := queue.GetResponse{
		Status: queue.StatusNoJob,
	}
	exp := "{\"status\":\"no-job\"}"

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	err := enc.Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	out := b.String()
	out = strings.TrimSpace(out)
	if out != exp {
		t.Fatalf("decode fail exp %v got %v", exp, out)
	}
}
