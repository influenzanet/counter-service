package internal

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestInfluenzanetConfigString(t *testing.T) {
	s := "myStudy"
	d, _ := GetInfluenzanetStudy(s)
	if d.StudyKey != s {
		t.Fatalf("Study key is not '%s'", s)
	}
}

func TestInfluenzanetConfigJSON(t *testing.T) {
	// Customize with JSON
	d, err := GetInfluenzanetStudy(`{"studykey":"myTest01", "active_surveys":["s1","s2"], "active_delay":"256h", "update_delay":"1h"}`)
	fmt.Println(d, err)
	if d.StudyKey != "myTest01" {
		t.Fatal("Study key is not 'myTest01'")
	}
	if d.ActiveParticipantDelay.Duration != time.Hour*256 {
		t.Fatal("Active participant delay should be 256h")
	}
	if d.UpdateDelay.Duration != time.Hour {
		t.Fatal("Active participant delay should be 1h")
	}
	if !reflect.DeepEqual(d.ActiveParticipantSurveys, []string{"s1", "s2"}) {
		t.Fatal("Active participant surveys should be s1,s2")
	}

}
