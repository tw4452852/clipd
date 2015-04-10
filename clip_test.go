package main

import (
	"log"
	"testing"

	"github.com/atotto/clipboard"
)

func TestMain(m *testing.M) {
	ori, err := clipboard.ReadAll()
	if err != nil {
		log.Println(err)
		return
	}
	defer clipboard.WriteAll(ori)
	m.Run()
}

func TestClip(t *testing.T) {
	begin, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	c := NewClip()
	if c.current != begin {
		t.Errorf("NewClip error: current[%q] should be [%q]\n", c.current, begin)
	}

	const freshContent = "fresh"
	err = clipboard.WriteAll(freshContent)
	if err != nil {
		t.Fatal(err)
	}
	remote := <-c.Inform()
	if remote != freshContent {
		t.Errorf("fresh error: remote[%q] should be [%q]\n", remote, freshContent)
	}
	if c.current != freshContent {
		t.Errorf("fresh error: current[%q] should be [%q]\n", c.current, freshContent)
	}

	const updateContent = "update"
	err = c.Update(updateContent)
	if err != nil {
		t.Errorf("update error: %s\n", err)
	}
	if c.current != updateContent {
		t.Errorf("update error: current[%q] should be [%q]\n", c.current, updateContent)
	}
	now, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if now != updateContent {
		t.Errorf("update error: clipboard[%q] should be [%q]\n", now, updateContent)
	}
}
