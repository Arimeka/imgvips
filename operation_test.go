package imgvips_test

import (
	"errors"
	"log"
	"testing"

	"github.com/Arimeka/imgvips"
)

func TestNewOperation(t *testing.T) {
	op, err := imgvips.NewOperation("jpegload")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	op.Free()

	_, err = imgvips.NewOperation("non_exists")
	if err == nil {
		t.Fatal("Expected to return error, got nil")
	}
}

func TestOperation_ExecFree(t *testing.T) {
	op, err := imgvips.NewOperation("jpegload")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	// Check multiply free
	op.Free()
	op.Free()

	if err := op.Exec(); err != imgvips.ErrOperationAlreadyFreed {
		t.Fatalf("Expected error %v, got %v", imgvips.ErrOperationAlreadyFreed, err)
	}
}

func TestOperation_Exec(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	out, loadOp := webpLoad(t)
	defer loadOp.Free()

	resizeOut, resizeOp := resize(t, out)
	defer resizeOp.Free()

	save(t, resizeOut)
}

func webpLoad(t *testing.T) (*imgvips.GValue, *imgvips.Operation) {
	op, err := imgvips.NewOperation("webpload")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	out := imgvips.GNullVipsImage()
	op.AddInput("filename", imgvips.GString("./tests/fixtures/img.webp"))
	op.AddInput("scale", imgvips.GDouble(0.1))
	op.AddOutput("out", out)

	if err := op.Exec(); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	return out, op
}

func resize(t *testing.T, in *imgvips.GValue) (*imgvips.GValue, *imgvips.Operation) {
	resizeOp, err := imgvips.NewOperation("resize")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	image, ok := in.Image()
	if !ok {
		t.Fatalf("Expected *C.VipsImage in out")
	}

	w := imgvips.ImageWidth(image)
	h := imgvips.ImageHeight(image)

	hScale := float64(350) / float64(h)
	wScale := float64(650) / float64(w)

	resizeOut := imgvips.GNullVipsImage()
	resizeOp.AddInput("in", in)
	resizeOp.AddInput("scale", imgvips.GDouble(wScale))
	resizeOp.AddInput("vscale", imgvips.GDouble(hScale))
	resizeOp.AddOutput("out", resizeOut)

	if err := resizeOp.Exec(); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	image, ok = resizeOut.Image()
	if !ok {
		t.Fatalf("Expected *C.VipsImage in out")
	}

	if imgvips.ImageHeight(image) != 350 {
		t.Errorf("Expected height %d, got %d", 350, imgvips.ImageHeight(image))
	}
	if imgvips.ImageWidth(image) != 650 {
		t.Errorf("Expected width %d, got %d", 650, imgvips.ImageWidth(image))
	}

	return resizeOut, resizeOp
}

func save(t *testing.T, in *imgvips.GValue) {
	saveOp, err := imgvips.NewOperation("jpegsave")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	defer saveOp.Free()

	saveOp.AddInput("in", in)
	saveOp.AddInput("filename", imgvips.GString("/dev/null"))
	saveOp.AddInput("Q", imgvips.GInt(50))

	if err := saveOp.Exec(); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func BenchmarkOperation_Exec(b *testing.B) {
	imgvips.VipsDetectMemoryLeak(true)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op, err := imgvips.NewOperation("webpload")
		if err != nil {
			b.Fatalf("Unexpected error %v", err)
		}

		out := imgvips.GNullVipsImage()
		op.AddInput("filename", imgvips.GString("./tests/fixtures/img.webp"))
		op.AddInput("scale", imgvips.GDouble(0.1))
		op.AddOutput("out", out)

		if err := op.Exec(); err != nil {
			op.Free()
			b.Fatalf("Unexpected error %v", err)
		}

		image, ok := out.Image()
		if !ok {
			op.Free()
			b.Fatalf("Expected *C.VipsImage in out")
		}

		resizeOp, err := imgvips.NewOperation("resize")
		if err != nil {
			b.Fatalf("Unexpected error %v", err)
		}

		w := imgvips.ImageWidth(image)
		h := imgvips.ImageHeight(image)

		hScale := float64(350) / float64(h)
		wScale := float64(650) / float64(w)

		resizeOut := imgvips.GNullVipsImage()
		resizeOp.AddInput("in", out)
		resizeOp.AddInput("scale", imgvips.GDouble(wScale))
		resizeOp.AddInput("vscale", imgvips.GDouble(hScale))
		resizeOp.AddOutput("out", resizeOut)

		if err := resizeOp.Exec(); err != nil {
			op.Free()
			resizeOp.Free()
			b.Fatalf("Unexpected error %v", err)
		}

		saveOp, err := imgvips.NewOperation("pngsave")
		if err != nil {
			op.Free()
			resizeOp.Free()
			b.Fatalf("Unexpected error %v", err)
		}

		saveOp.AddInput("in", resizeOut)
		saveOp.AddInput("filename", imgvips.GString("/dev/null"))

		if err := saveOp.Exec(); err != nil {
			op.Free()
			resizeOp.Free()
			saveOp.Free()
			b.Fatalf("Unexpected error %v", err)
		}
		op.Free()
		resizeOp.Free()
		saveOp.Free()
	}
}

func ExampleOperation_Exec() {
	op, err := imgvips.NewOperation("webpload")
	if err != nil {
		log.Println(err)
		return
	}
	defer op.Free()

	out := imgvips.GNullVipsImage()
	op.AddInput("filename", imgvips.GString("./tests/fixtures/img.webp"))
	op.AddInput("scale", imgvips.GDouble(0.1))
	op.AddOutput("out", out)

	if err := op.Exec(); err != nil {
		log.Println(err)
		return
	}

	resizeOp, err := imgvips.NewOperation("resize")
	if err != nil {
		log.Println(err)
		return
	}
	defer resizeOp.Free()

	image, ok := out.Image()
	if !ok {
		log.Println(errors.New("out is not *C.VipsImage"))
		return
	}

	w := imgvips.ImageWidth(image)
	h := imgvips.ImageHeight(image)

	hScale := float64(350) / float64(h)
	wScale := float64(650) / float64(w)

	resizeOut := imgvips.GNullVipsImage()
	resizeOp.AddInput("in", out)
	resizeOp.AddInput("scale", imgvips.GDouble(wScale))
	resizeOp.AddInput("vscale", imgvips.GDouble(hScale))
	resizeOp.AddOutput("out", resizeOut)

	if err := resizeOp.Exec(); err != nil {
		log.Println(err)
		return
	}

	saveOp, err := imgvips.NewOperation("webpsave")
	if err != nil {
		log.Println(err)
		return
	}
	defer saveOp.Free()

	saveOp.AddInput("in", resizeOut)
	saveOp.AddInput("filename", imgvips.GString("./tests/fixtures/resized.webp"))
	saveOp.AddInput("Q", imgvips.GInt(50))
	saveOp.AddInput("strip", imgvips.GBoolean(true))

	if err := saveOp.Exec(); err != nil {
		log.Println(err)
		return
	}
}
