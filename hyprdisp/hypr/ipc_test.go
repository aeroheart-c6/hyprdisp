package hypr

import "testing"

func TestOverflowConcat(t *testing.T) {
	var (
		overflow []byte = []byte("event>>d")
		data     []byte = []byte("ata1,data2,data3")
		result   string
		expect   string
	)

	result = string(append(overflow, data...))
	expect = "event>>data1,data2,data3"
	if result != expect {
		t.Fatalf("Expected: \"%s\" but got \"%s\"", expect, result)
	}

	var dataSize = 4
	result = string(append(overflow, data[:dataSize]...))
	expect = "event>>data1"
	if result != expect {
		t.Fatalf("Expected: \"%s\" but got \"%s\"", expect, result)
	}
}

func TestParseEvents(t *testing.T) {
	var payloadFull []byte = []byte("" +
		"event00>>data1,data2,data3\n" +
		"event00>>data4,data5,data6\n" +
		"event01>>data1\n")

	var (
		dataSize int    = 30
		buffer   []byte = payloadFull[:dataSize]
		overflow []byte
		events   []Event
	)

	events, overflow = parseEvents(buffer)

	if len(events) != 1 {
		t.Fatalf("Expected 1 item but got %d", len(events))
	}

	if events[0].Name != "event00" {
		t.Fatalf("Expected 1 item to be of \"event00\" but got \"%s\"", events[0].Name)
	}

	if string(overflow) != "eve" {
		t.Fatalf("Expected the left over \"eve\", but got %s", string(overflow))
	}
}
