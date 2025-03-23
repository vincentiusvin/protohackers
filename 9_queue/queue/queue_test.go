package queue_test

import (
	"bytes"
	"context"
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

func intPtr(v int) *int {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func rawMsgPtr(v json.RawMessage) *json.RawMessage {
	return &v
}

func TestEncode(t *testing.T) {
	type encodeCases struct {
		in  any
		exp string
	}

	runTest := func(t *testing.T, cs []encodeCases) {
		for _, c := range cs {
			b := new(bytes.Buffer)
			enc := json.NewEncoder(b)
			err := enc.Encode(c.in)
			if err != nil {
				t.Fatal(err)
			}
			out := b.String()
			out = strings.TrimSpace(out)
			if out != c.exp {
				t.Fatalf("decode fail exp %v got %v", c.exp, out)
			}
		}
	}

	t.Run("get cases", func(t *testing.T) {
		cases := []encodeCases{
			{
				in:  queue.GetResponse{Status: queue.StatusNoJob},
				exp: "{\"status\":\"no-job\"}",
			},
			{
				in: queue.GetResponse{
					Status: queue.StatusOK,
					Id:     intPtr(12345),
					Pri:    intPtr(123),
					Queue:  strPtr("queue1"),
					Job:    rawMsgPtr([]byte("{\"f\":1}")),
				},
				exp: "{\"status\":\"ok\",\"id\":12345,\"job\":{\"f\":1},\"pri\":123,\"queue\":\"queue1\"}",
			},
		}
		runTest(t, cases)
	})
}

func TestJobCenter(t *testing.T) {
	ctx := context.Background()
	jc := queue.NewJobCenter(ctx)
	req := &queue.PutRequest{
		Request: "put",
		Queue:   "test",
		Pri:     123,
		Job:     []byte("{\"f\":1}"),
	}
	resp := jc.Put(req)

	if resp.Status != queue.StatusOK {
		t.Fatalf("failed to put")
	}

}
