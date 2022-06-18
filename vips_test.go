package imgvips_test

import (
	"testing"

	"github.com/Arimeka/imgvips"
)

func initVips(t testing.TB) {
	err := imgvips.Initialize(imgvips.VipsCacheSetMaxMem(-10), imgvips.VipsCacheSetMax(-10),
		imgvips.VipsVectorSetEnables(false), imgvips.VipsDetectMemoryLeak(true),
		imgvips.VipsConcurrencySet(0))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestGetMemCacheOFF(t *testing.T) {
	initVips(t)

	_, op := webpLoadBytes(t)

	if imgvips.GetMem() <= 0 {
		t.Error("expected take memory")
	}

	op.Free()

	if imgvips.GetMem() != 0 {
		t.Errorf("expected free memory, got %f", imgvips.GetMem())
	}
}

func TestGetMemCacheON(t *testing.T) {
	err := imgvips.Initialize(imgvips.VipsCacheSetMaxMem(1024*1024*100), imgvips.VipsCacheSetMax(1024),
		imgvips.VipsVectorSetEnables(false), imgvips.VipsDetectMemoryLeak(true),
		imgvips.VipsConcurrencySet(-10))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	_, op := webpLoadBytes(t)

	if imgvips.GetMem() <= 0 {
		t.Error("expected take memory")
	}

	beforeFree := imgvips.GetMem()
	op.Free()

	if imgvips.GetMem() == 0 || imgvips.GetMem() != beforeFree {
		t.Errorf("expected memory %f, got %f", beforeFree, imgvips.GetMem())
	}
}

func TestGetMemHighwaterCacheOFF(t *testing.T) {
	initVips(t)

	_, op := webpLoadBytes(t)

	if imgvips.GetMemHighwater() <= 0 {
		t.Error("expected memory high-water bigger than 0")
	}

	op.Free()

	if imgvips.GetMemHighwater() <= 0 {
		t.Error("expected memory high-water bigger than 0")
	}
}

func TestGetAllocsCacheOFF(t *testing.T) {
	initVips(t)

	_, op := webpLoadBytes(t)

	if imgvips.GetAllocs() <= 0 {
		t.Error("expected memory allocations bigger than 0")
	}

	op.Free()

	if imgvips.GetAllocs() != 0 {
		t.Errorf("expected free memory allocations, got %f", imgvips.GetAllocs())
	}
}

func TestGetAllocsCacheON(t *testing.T) {
	err := imgvips.Initialize(imgvips.VipsCacheSetMaxMem(1024*1024*100), imgvips.VipsCacheSetMax(1024),
		imgvips.VipsVectorSetEnables(false), imgvips.VipsDetectMemoryLeak(true),
		imgvips.VipsConcurrencySet(-10))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	_, op := webpLoadBytes(t)

	if imgvips.GetAllocs() <= 0 {
		t.Error("expected memory allocations bigger than 0")
	}
	beforeFree := imgvips.GetAllocs()
	op.Free()

	if imgvips.GetAllocs() == 0 || imgvips.GetAllocs() != beforeFree {
		t.Errorf("expected memory allocations %f, got %f", beforeFree, imgvips.GetAllocs())
	}
}
