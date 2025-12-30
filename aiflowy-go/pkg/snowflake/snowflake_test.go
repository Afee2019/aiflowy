package snowflake

import (
	"sync"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		workerID     int64
		datacenterID int64
		expectError  bool
	}{
		{"valid IDs", 1, 1, false},
		{"zero IDs", 0, 0, false},
		{"max valid IDs", 31, 31, false},
		{"worker ID too high", 32, 1, true},
		{"datacenter ID too high", 1, 32, true},
		{"negative worker ID", -1, 1, true},
		{"negative datacenter ID", 1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset singleton for each test
			once = sync.Once{}
			instance = nil

			err := Init(tt.workerID, tt.datacenterID)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestNextID(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	err := Init(1, 1)
	if err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	// Generate IDs
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id, err := sf.NextID()
		if err != nil {
			t.Fatalf("failed to generate ID: %v", err)
		}
		if id <= 0 {
			t.Errorf("generated ID should be positive, got: %d", id)
		}
		if ids[id] {
			t.Errorf("duplicate ID generated: %d", id)
		}
		ids[id] = true
	}
}

func TestNextIDConcurrency(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	err := Init(1, 1)
	if err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	const numGoroutines = 100
	const idsPerGoroutine = 100

	var mu sync.Mutex
	ids := make(map[int64]bool)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := sf.NextID()
				if err != nil {
					t.Errorf("failed to generate ID: %v", err)
					return
				}
				mu.Lock()
				if ids[id] {
					t.Errorf("duplicate ID generated: %d", id)
				}
				ids[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(ids) != numGoroutines*idsPerGoroutine {
		t.Errorf("expected %d unique IDs, got %d", numGoroutines*idsPerGoroutine, len(ids))
	}
}

func TestNextIDString(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	err := Init(1, 1)
	if err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	idStr, err := sf.NextIDString()
	if err != nil {
		t.Fatalf("failed to generate ID string: %v", err)
	}
	if idStr == "" {
		t.Error("generated ID string should not be empty")
	}
}

func TestParseID(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	workerID := int64(5)
	datacenterID := int64(10)

	err := Init(workerID, datacenterID)
	if err != nil {
		t.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	id, err := sf.NextID()
	if err != nil {
		t.Fatalf("failed to generate ID: %v", err)
	}

	ts, dc, worker, seq := ParseID(id)

	if dc != datacenterID {
		t.Errorf("expected datacenterID %d, got %d", datacenterID, dc)
	}
	if worker != workerID {
		t.Errorf("expected workerID %d, got %d", workerID, worker)
	}
	if seq < 0 {
		t.Errorf("sequence should be non-negative, got %d", seq)
	}
	if ts.Before(time.Now().Add(-time.Second)) || ts.After(time.Now().Add(time.Second)) {
		t.Errorf("timestamp should be approximately now, got %v", ts)
	}
}

func TestGenerateID(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	id, err := GenerateID()
	if err != nil {
		t.Fatalf("failed to generate ID: %v", err)
	}
	if id <= 0 {
		t.Errorf("generated ID should be positive, got: %d", id)
	}
}

func TestGenerateIDString(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	idStr, err := GenerateIDString()
	if err != nil {
		t.Fatalf("failed to generate ID string: %v", err)
	}
	if idStr == "" {
		t.Error("generated ID string should not be empty")
	}
}

func TestMustGenerateID(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustGenerateID should not panic: %v", r)
		}
	}()

	id := MustGenerateID()
	if id <= 0 {
		t.Errorf("generated ID should be positive, got: %d", id)
	}
}

func TestMustGenerateIDString(t *testing.T) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustGenerateIDString should not panic: %v", r)
		}
	}()

	idStr := MustGenerateIDString()
	if idStr == "" {
		t.Error("generated ID string should not be empty")
	}
}

func BenchmarkNextID(b *testing.B) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	err := Init(1, 1)
	if err != nil {
		b.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sf.NextID()
	}
}

func BenchmarkNextIDConcurrent(b *testing.B) {
	// Reset singleton
	once = sync.Once{}
	instance = nil

	err := Init(1, 1)
	if err != nil {
		b.Fatalf("failed to init: %v", err)
	}

	sf := GetInstance()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = sf.NextID()
		}
	})
}
