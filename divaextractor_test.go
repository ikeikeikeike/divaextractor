package divaextractor

import "testing"

func TestDocument(t *testing.T) {
	c := NewWikipedia()

	err := c.Do("明日花キララ")
	if err != nil {
		t.Fatal(err)
	}

	if c.Birthday().Unix() != 591753600 {
		t.Errorf("Unexpected Wikipedia.Birthday: %s", c.Birthday())
	}

	if c.Blood() != "A" {
		t.Errorf("Unexpected Wikipedia.Blood: %s", c.Blood())
	}

	if c.HW() != "162  / ―" {
		t.Errorf("Unexpected Wikipedia.H/Weight: %s", c.HW())
	}
	if c.Height() != 162 {
		t.Errorf("Unexpected Wikipedia.Height: %d", c.Height())
	}
	if c.Weight() != 0 {
		t.Errorf("Unexpected Wikipedia.Weight: %d", c.Weight())
	}

	if c.BWH() != "90 - 58 - 85" {
		t.Errorf("Unexpected Wikipedia.BWH: %s", c.BWH())
	}
	if c.Bust() != 90 {
		t.Errorf("Unexpected Wikipedia.Bust: %d", c.Bust())
	}
	if c.Waste() != 58 {
		t.Errorf("Unexpected Wikipedia.Waste: %d", c.Waste())
	}
	if c.Hip() != 85 {
		t.Errorf("Unexpected Wikipedia.Hip: %d", c.Hip())
	}

	if c.Bracup() != "G" {
		t.Errorf("Unexpected Wikipedia.Bracup: %s", c.Bracup())
	}
}
