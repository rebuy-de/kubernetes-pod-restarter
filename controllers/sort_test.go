package controllers

import (
	"sort"
	"testing"
	"time"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSortPodByAge(t *testing.T) {
	ages := []string{
		"2006-01-02T16:00:00Z",
		"2006-01-02T14:00:00Z",
		"2006-01-02T15:00:00Z",
	}

	podList := new(core.PodList)

	for _, age := range ages {
		ts, err := time.Parse(time.RFC3339, age)
		if err != nil {
			t.Fatal(err)
		}

		podList.Items = append(podList.Items, core.Pod{
			ObjectMeta: meta.ObjectMeta{
				CreationTimestamp: meta.NewTime(ts),
			},
		})
	}

	sort.Sort(PodsByAge(*podList))

	ts0 := podList.Items[0].ObjectMeta.CreationTimestamp.Time.Format(time.RFC3339)
	ts1 := podList.Items[1].ObjectMeta.CreationTimestamp.Time.Format(time.RFC3339)
	ts2 := podList.Items[2].ObjectMeta.CreationTimestamp.Time.Format(time.RFC3339)

	if ts0 != "2006-01-02T14:00:00Z" {
		t.Fatal("Item #0 is wrong")
	}

	if ts1 != "2006-01-02T15:00:00Z" {
		t.Fatal("Item #1 is wrong")
	}

	if ts2 != "2006-01-02T16:00:00Z" {
		t.Fatal("Item #2 is wrong")
	}
}
