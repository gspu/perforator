package uprobe

import (
	"fmt"
	"os"
	"sync"

	"github.com/yandex/perforator/perforator/pkg/xelf"
)

type Key struct {
	BuildID string
	Offset  uint64
}

type SymbolInfo struct {
	Name        string
	LocalOffset uint64
}

type UprobeInfo struct {
	SymbolInfo
	SampleType string
}

type Registry struct {
	mutex   sync.RWMutex
	uprobes map[Key]*UprobeInfo
}

func NewRegistry() *Registry {
	return &Registry{
		uprobes: make(map[Key]*UprobeInfo),
	}
}

func extractOffset(binaryFile *os.File, symbol string, localOffset uint64) (uint64, error) {
	offsets, err := xelf.GetSymbolFileOffsets(binaryFile, symbol)
	if err != nil {
		return 0, fmt.Errorf("failed to get symbol offset: %w", err)
	}

	symbolOffset, ok := offsets[symbol]
	if !ok {
		return 0, fmt.Errorf("symbol not found: %s", symbol)
	}

	return symbolOffset + localOffset, nil
}

func extractKey(file *os.File, offset uint64) (Key, error) {
	buildID, err := xelf.ReadBuildID(file)
	if err != nil {
		return Key{}, fmt.Errorf("failed to get build ID: %w", err)
	}

	return Key{
		Offset:  offset,
		BuildID: buildID,
	}, nil
}

type options struct {
	pid uint32
}

type Option func(*options)

func WithPID(pid uint32) Option {
	return func(o *options) {
		o.pid = pid
	}
}

func (r *Registry) Create(config Config, optionAppliers ...Option) Uprobe {
	opts := &options{}
	for _, opt := range optionAppliers {
		opt(opts)
	}

	return &uprobe{
		config: config,
		reg:    r,
		opts:   opts,
	}
}

func (r *Registry) removeResolveInfo(key Key) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.uprobes, key)
}

func (r *Registry) addResolveInfo(key Key, symbolName string, localOffset uint64, sampleType string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.uprobes[key] = &UprobeInfo{
		SymbolInfo: SymbolInfo{
			Name:        symbolName,
			LocalOffset: localOffset,
		},
		SampleType: sampleType,
	}
}

func (r *Registry) ResolveUprobe(key Key) *UprobeInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.uprobes[key]
}
