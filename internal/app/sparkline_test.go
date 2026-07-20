package app

import "testing"

func TestRenderBrailleSparklineEmpty(t *testing.T) {
	got := RenderBrailleSparkline(nil, 5, 2)
	if got != "" {
		t.Errorf("empty values should return empty string, got %q", got)
	}
}

func TestRenderBrailleSparklineZeroWidth(t *testing.T) {
	got := RenderBrailleSparkline([]float64{1, 2, 3}, 0, 2)
	if got != "" {
		t.Errorf("zero width should return empty string, got %q", got)
	}
}

func TestRenderBrailleSparklineZeroHeight(t *testing.T) {
	got := RenderBrailleSparkline([]float64{1, 2, 3}, 5, 0)
	if got != "" {
		t.Errorf("zero height should return empty string, got %q", got)
	}
}

func TestRenderBrailleSparklineSingleValue(t *testing.T) {
	got := RenderBrailleSparkline([]float64{42}, 3, 2)
	if got == "" {
		t.Error("single value should produce output")
	}
}

func TestRenderBrailleSparklineMultipleValues(t *testing.T) {
	got := RenderBrailleSparkline([]float64{1, 2, 3, 4, 5}, 5, 2)
	if got == "" {
		t.Error("multiple values should produce output")
	}
}

func TestRenderBrailleSparklineAllZeros(t *testing.T) {
	got := RenderBrailleSparkline([]float64{0, 0, 0}, 3, 1)
	if got == "" {
		t.Error("all zeros should still produce output (maxVal=1 fallback)")
	}
}
