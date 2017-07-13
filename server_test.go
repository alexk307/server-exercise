package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDatabase_CheckForCompletedMessageSimple(t *testing.T) {
	transactionID := 1
	d := NewDatabase()
	d.Data[transactionID] = map[int]string{}
	d.Data[transactionID][1] = "byte"

	if len(d.CheckForCompletedMessage(transactionID)) != 0 {
		t.Fail()
	}
}

func TestDatabase_CheckForCompletedMessageOffset(t *testing.T) {
	transactionID := 1
	d := NewDatabase()
	d.Data[transactionID] = map[int]string{}
	d.Data[transactionID][10] = "byte"
	d.Data[transactionID][11] = "byte"
	d.Data[transactionID][12] = "byte"

	if len(d.CheckForCompletedMessage(transactionID)) != 0 {
		t.Fail()
	}
}

func TestDatabase_CheckForCompletedMessageFail(t *testing.T) {
	transactionID := 1
	d := NewDatabase()
	d.Data[transactionID] = map[int]string{}
	d.Data[transactionID][1] = "byte"
	d.Data[transactionID][3] = "byte"

	resp := d.CheckForCompletedMessage(transactionID)
	if len(resp) == 0 {
		t.Fail()
	}
	if reflect.DeepEqual(resp, []int{1, 3}) != true {
		t.Fail()
	}
}

func TestDatabase_CheckForCompletedMessageFailOffset(t *testing.T) {
	transactionID := 1
	d := NewDatabase()
	d.Data[transactionID] = map[int]string{}
	d.Data[transactionID][100] = "byte"
	d.Data[transactionID][101] = "byte"
	d.Data[transactionID][102] = "byte"
	d.Data[transactionID][104] = "byte"
	d.Data[transactionID][105] = "byte"

	resp := d.CheckForCompletedMessage(transactionID)
	if len(resp) == 0 {
		t.Fail()
	}
	if reflect.DeepEqual(resp, []int{102, 104}) != true {
		fmt.Println(resp)
		t.Fail()
	}
}

func TestDatabase_CheckForCompletedMessageFailOffsetMulti(t *testing.T) {
	transactionID := 1
	d := NewDatabase()
	d.Data[transactionID] = map[int]string{}
	d.Data[transactionID][100] = "byte"
	d.Data[transactionID][101] = "byte"
	d.Data[transactionID][103] = "byte"
	d.Data[transactionID][104] = "byte"
	d.Data[transactionID][105] = "byte"
	d.Data[transactionID][107] = "byte"

	resp := d.CheckForCompletedMessage(transactionID)
	if len(resp) == 0 {
		t.Fail()
	}
	if reflect.DeepEqual(resp, []int{101, 103, 105, 107}) != true {
		fmt.Println(resp)
		t.Fail()
	}
}
