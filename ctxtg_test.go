package ctxtg

import (
	"reflect"
	"testing"
	"time"

	"context"
)

func TestAddDataToEmptyContext(t *testing.T) {
	ctx := context.Background()
	k := "key"
	v := "value"
	ctx = WithDataValue(ctx, k, v)
	data := ctx.Value(DataKey).(map[string]interface{})
	if v != data[k] {
		t.Error("Valud not saved")
	}
}

func TestAddDataToContext(t *testing.T) {
	k1 := "key1"
	v1 := 1

	data := map[string]interface{}{
		k1: v1,
	}
	k := "key"
	v := "value"

	ctx := context.WithValue(context.Background(), DataKey, data)
	ctx = WithDataValue(ctx, k, v)
	d := ctx.Value(DataKey).(map[string]interface{})
	if d[k1] != v1 && d[k] != v {
		t.Error("Data add incorrectly")
	}
}

func TestToContext(t *testing.T) {
	token := Token("tokentest")
	deadline := time.Now().Add(10 * time.Second).Unix()
	trackingID := "123123"
	data := map[string]interface{}{
		"1": 123,
		"2": "string",
		"3": 3 * time.Second,
	}
	c := Context{
		Token:     token,
		Deadline:  deadline,
		TracingID: trackingID,
		Data:      data,
	}

	ctx, cancel := c.ToContext()
	defer cancel()
	if cancel == nil {
		t.Errorf("Invalid cancel func")
	}
	if d, ok := ctx.Deadline(); !ok || d.Unix() != deadline {
		t.Errorf("Invalid deadline %v %v", d.Unix(), ok)
	}
	if tok := ctx.Value(TokenKey); tok != token {
		t.Errorf("Invalid token %v", tok)
	}
	if tID := ctx.Value(TracingIDKey); tID != trackingID {
		t.Errorf("Invalid trackingID %v", tID)
	}
	if d := ctx.Value(DataKey); !reflect.DeepEqual(d, data) {
		t.Errorf("Invalid data %v", d)
	}
}

func TestEmptyToContext(t *testing.T) {
	var c Context
	ctx, cancel := c.ToContext()
	defer cancel()
	if cancel == nil {
		t.Errorf("Invalid cancel func")
	}
	if _, ok := ctx.Deadline(); ok {
		t.Error("Invalid deadline")
	}
	if tok := ctx.Value(TokenKey); tok != Token("") {
		t.Errorf("Invalid token %v", tok)
	}
	if tID := ctx.Value(TracingIDKey); tID != "" {
		t.Errorf("Invalid trackingID %v", tID)
	}
	if d := ctx.Value(DataKey); d != nil {
		t.Errorf("Invalid data %v", d)
	}
}

func TestFromContext(t *testing.T) {
	token := Token("tokentest")
	deadline := time.Now().Add(10 * time.Second).Unix()
	trackingID := "123123"
	data := map[string]interface{}{
		"1": 123,
		"2": "string",
		"3": 3 * time.Second,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, TokenKey, token)
	ctx = context.WithValue(ctx, TracingIDKey, trackingID)
	ctx = context.WithValue(ctx, DataKey, data)
	ctx = context.WithValue(ctx, key(1000000000), 234)
	ctx, cancel := context.WithDeadline(ctx, time.Unix(deadline, 0))
	defer cancel()

	c := FromContext(ctx)

	if c.Deadline != deadline {
		t.Errorf("Invalid deadline %v", c.Deadline)
	}
	if c.Token != token {
		t.Errorf("Invalid token %v", c.Token)
	}
	if c.TracingID != trackingID {
		t.Errorf("Invalid trackingID %v", c.TracingID)
	}
	if !reflect.DeepEqual(c.Data, data) {
		t.Errorf("Invalid data %v", c.Data)
	}
}

func TestFromToContext(t *testing.T) {
	token := "tokentest"
	deadline := time.Now().Add(10 * time.Second).Unix()
	trackingID := "123123"
	data := map[string]interface{}{
		"1": 123,
		"2": "string",
		"3": 3 * time.Second,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, TokenKey, token)
	ctx = context.WithValue(ctx, TracingIDKey, trackingID)
	ctx = context.WithValue(ctx, DataKey, data)
	ctx = context.WithValue(ctx, key(10325234562343), 234)
	ctx, cancel := context.WithDeadline(ctx, time.Unix(deadline, 0))
	defer cancel()

	c := FromContext(ctx)
	ctx2, _ := c.ToContext()
	c2 := FromContext(ctx2)

	if !reflect.DeepEqual(c, c2) {
		t.Errorf("Should be the same %v != %v", c, c2)
	}
}

func TestEmptyFromContext(t *testing.T) {
	c := FromContext(context.Background())

	if c.Deadline != 0 {
		t.Errorf("Invalid deadline %v", c.Deadline)
	}
	if c.Token != "" {
		t.Errorf("Invalid token %v", c.Token)
	}
	if c.TracingID != "" {
		t.Errorf("Invalid trackingID %v", c.TracingID)
	}
	if c.Data != nil {
		t.Errorf("Invalid data %v", c.Data)
	}
}

func TestValueFromData(t *testing.T) {
	k1 := "k"
	v1 := 123
	k2 := "k2"
	v2 := "12345"
	testData := map[string]interface{}{
		k1: v1,
		k2: v2,
	}
	ctx := context.WithValue(context.Background(), DataKey, testData)
	if ValueFromData(ctx, k1) != v1 {
		t.Errorf("Invalid v1 %v", v1)
	}
	if ValueFromData(ctx, k2) != v2 {
		t.Errorf("Invalid v2 %v", v2)
	}
	if ValueFromData(ctx, "no such key") != nil {
		t.Errorf("Invalid key should return nil")
	}
}

func TestValueFromDataWithEmptyCtx(t *testing.T) {
	if ValueFromData(context.Background(), "no such key") != nil {
		t.Errorf("Context without data should return nil")
	}
}
