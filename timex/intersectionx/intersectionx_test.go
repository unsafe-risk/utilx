package intersectionx_test

import (
	"testing"

	"github.com/unsafe-risk/utilx/timex/intersectionx"
)

func TestIntersection(t *testing.T) {
	t.Run("one interval", func(t *testing.T) {
		lower, upper, ok := intersectionx.Intersection(
			intersectionx.Interval{
				Offset:         0,
				ConfidenceBand: 10,
			},
		)

		if !ok {
			t.Fatal("expected ok")
		}

		if lower != -10 {
			t.Fatalf("expected lower to be -10, got %d", lower)
		}

		if upper != 10 {
			t.Fatalf("expected upper to be 10, got %d", upper)
		}
	})

	t.Run("two same intervals", func(t *testing.T) {
		lower, upper, ok := intersectionx.Intersection(
			intersectionx.Interval{
				Offset:         0,
				ConfidenceBand: 10,
			},
			intersectionx.Interval{
				Offset:         0,
				ConfidenceBand: 10,
			},
		)

		if !ok {
			t.Fatal("expected ok")
		}

		if lower != -10 {
			t.Fatalf("expected lower to be -10, got %d", lower)
		}

		if upper != 10 {
			t.Fatalf("expected upper to be 10, got %d", upper)
		}
	})

	t.Run("two different intervals", func(t *testing.T) {
		lower, upper, ok := intersectionx.Intersection(
			intersectionx.Interval{
				Offset:         0,
				ConfidenceBand: 10,
			},
			intersectionx.Interval{
				Offset:         2,
				ConfidenceBand: 5,
			},
		)

		if !ok {
			t.Fatal("expected ok")
		}

		if lower != -3 {
			t.Fatalf("expected lower to be -3, got %d", lower)
		}

		if upper != 7 {
			t.Fatalf("expected upper to be 7, got %d", upper)
		}
	})

	t.Run("three different intervals", func(t *testing.T) {
		lower, upper, ok := intersectionx.Intersection(
			intersectionx.Interval{
				Offset:         5,
				ConfidenceBand: 10,
			},
			intersectionx.Interval{
				Offset:         2,
				ConfidenceBand: 5,
			},
			intersectionx.Interval{
				Offset:         3,
				ConfidenceBand: 6,
			},
		)

		if !ok {
			t.Fatal("expected ok")
		}

		if lower != -3 {
			t.Fatalf("expected lower to be -3, got %d", lower)
		}

		if upper != 7 {
			t.Fatalf("expected upper to be 7, got %d", upper)
		}
	})
}
