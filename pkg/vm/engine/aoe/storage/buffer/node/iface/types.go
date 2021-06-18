package iface

import (
	"fmt"
	"io"
	buf "matrixone/pkg/vm/engine/aoe/storage/buffer"
	"sync"
	"sync/atomic"
)

var (
	NODE_HEAD_SIZE  uint64 = 32
	NODE_ALLOC_SIZE uint64 = 256 * 1024
	NODE_DATA_SIZE         = NODE_ALLOC_SIZE - NODE_HEAD_SIZE
)

type NodeState = uint32

const (
	NODE_UNLOAD NodeState = iota
	NODE_LOADING
	NODE_ROOLBACK
	NODE_COMMIT
	NODE_UNLOADING
	NODE_LOADED
)

func NodeStateString(state NodeState) string {
	switch state {
	case NODE_UNLOAD:
		return "UNLOAD"
	case NODE_LOADING:
		return "LOADING"
	case NODE_ROOLBACK:
		return "ROLLBACK"
	case NODE_COMMIT:
		return "COMMIT"
	case NODE_UNLOADING:
		return "UNLOADING"
	case NODE_LOADED:
		return "LOADED"
	}
	panic(fmt.Sprintf("unsupported: %d", state))
}

func AtomicLoadState(addr *NodeState) NodeState {
	return atomic.LoadUint32(addr)
}

func AtomicStoreState(addr *NodeState, val NodeState) {
	atomic.StoreUint32(addr, val)
}

func AtomicCASState(addr *NodeState, old, new NodeState) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

type NodeRTState = uint32

const (
	NODE_RT_RUNNING NodeRTState = iota
	NODE_RT_CLOSED
)

func AtomicLoadRTState(addr *NodeRTState) NodeRTState {
	return atomic.LoadUint32(addr)
}

func AtomicStoreRTState(addr *NodeRTState, val NodeRTState) {
	atomic.StoreUint32(addr, val)
}

func AtomicCASRTState(addr *NodeRTState, old, new NodeRTState) bool {
	return atomic.CompareAndSwapUint32(addr, old, new)
}

type BufferType uint8

const (
	STATIC_SIZED BufferType = iota
	DYNAMIC_SIZED
)

type INodeBuffer interface {
	buf.IBuffer
	GetID() uint64
	// GetType() BufferType
}

type INodeHandle interface {
	sync.Locker
	io.Closer
	GetID() uint64
	Unload()
	// Loadable() bool
	Unloadable() bool
	// GetBuff() buf.IBuffer
	PrepareLoad() bool
	RollbackLoad()
	CommitLoad() error
	MakeHandle() IBufferHandle
	GetState() NodeState
	GetCapacity() uint64
	// Size() uint64
	// IsDestroyable() bool
	IsClosed() bool
	Ref()
	// If the current Refs is already 0, it returns false, else true
	UnRef() bool
	// If the current Refs is not 0, it returns true, else false
	HasRef() bool
	SetBuffer(buffer buf.IBuffer) error
	Iteration() uint64
	IncIteration() uint64
	GetBuffer() buf.IBuffer
	IsSpillable() bool
	Clean() error
}

type IBufferHandle interface {
	io.Closer
	GetID() uint64
}