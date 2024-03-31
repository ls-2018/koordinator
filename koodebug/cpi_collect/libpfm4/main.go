package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"go.uber.org/multierr"
	"golang.org/x/sys/unix"
	"k8s.io/klog/v2"
)

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -lpfm
#include <perfmon/pfmlib.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"

const (
	CYCLES       = "cycles"
	INSTRUCTIONS = "instructions"
)

var (
	initlibpfm    sync.Once
	closelibpfm   sync.Once
	perfValuePool sync.Pool
	BufPools      map[int]*sync.Pool
	attrMap       = make(map[string]*unix.PerfEventAttr)
)

func InitBufferPool() {
	BufPools = make(map[int]*sync.Pool)
	EventsMap = map[string][]string{
		"CPICollector": {"cycles", "instructions", "LONGEST_LAT_CACHE.MISS"},
	}
	for _, events := range EventsMap {
		BufPools[len(events)] = &sync.Pool{
			New: func() interface{} {
				// https://man7.org/linux/man-pages/man2/perf_event_open.2.html#Reading%20results
				// 24 means the size of nr, time_enabled, time_running
				// 16 means the size of value, id
				// struct read_format {
				//  u64 nr;            /* The number of events */
				//  u64 time_enabled;  /* if PERF_FORMAT_TOTAL_TIME_ENABLED */
				//  u64 time_running;  /* if PERF_FORMAT_TOTAL_TIME_RUNNING */
				//  struct {
				//      u64 value;     /* The value of the event */
				//      u64 id;        /* if PERF_FORMAT_ID */
				//      u64 lost;      /* if PERF_FORMAT_LOST */
				//  } values[nr];
				// };
				pool := make([]byte, 24+len(events)*16)
				return &pool
			},
		}
	}
}

func LibInit() {
	perfValuePool.New = func() interface{} {
		return &value{}
	}
	initlibpfm.Do(func() {
		if err := C.pfm_initialize(); err != C.PFM_SUCCESS {
			klog.Fatalf("Unable to init libpfm: %v", err)
		}
	})
	for _, events := range EventsMap {
		for _, event := range events {
			attr, err := createPerfConfig(event)
			if err != nil {
				panic(err)
			}
			attr.Read_format = unix.PERF_FORMAT_GROUP | unix.PERF_FORMAT_TOTAL_TIME_ENABLED | unix.PERF_FORMAT_TOTAL_TIME_RUNNING | unix.PERF_FORMAT_ID
			attr.Sample_type = unix.PERF_SAMPLE_IDENTIFIER
			attr.Size = uint32(unsafe.Sizeof(unix.PerfEventAttr{}))
			attr.Bits |= unix.PerfBitInherit
			attr.Bits |= unix.PerfBitDisabled
			attrMap[event] = attr
		}
	}
}

func LibFinalize() {
	klog.Info("termiate")
	closelibpfm.Do(func() {
		C.pfm_terminate()
	})
	for _, attr := range attrMap {
		C.free(unsafe.Pointer(attr))
	}
}

type PerfGroupCollector struct {
	cgroupFile     *os.File
	cpus           []int
	perfCollectors sync.Map
	idEventMap     map[uint64]string
	mu             sync.Mutex
	resultMap      map[string]float64
	valueCh        chan perfValue
	syscall6       func(trap uintptr, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr, a5 uintptr, a6 uintptr) (r1 uintptr, r2 uintptr, err syscall.Errno)
	closeCh        chan struct{}
}

func (pgc *PerfGroupCollector) GetPerfCollectors(cpu int) (*PerfCollector, bool) {
	pcInf, ok := pgc.perfCollectors.Load(cpu)
	if ok {
		return pcInf.(*PerfCollector), ok
	}
	return nil, false
}

type PerfCollector struct {
	cpu      int
	syscall6 func(trap uintptr, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr, a5 uintptr, a6 uintptr) (r1 uintptr, r2 uintptr, err syscall.Errno)
	leaderFd io.ReadCloser
	fds      []io.ReadCloser
}

type perfValue struct {
	Value float64
	ID    uint64
}

type value struct {
	Value uint64
	ID    uint64
}

type perfValueHeader struct {
	Nr          uint64 // number of events
	TimeEnabled uint64 // time event active
	TimeRunning uint64 // time event on CPU
}

func CreatePerfGroupCollector(cgroupFile *os.File, cpus []int, events []string, syscallFunc func(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)) (collector *PerfGroupCollector, err error) {
	if len(events) == 0 {
		err = errors.New("events cannot be empty")
		return nil, err
	}
	collector = &PerfGroupCollector{
		cgroupFile:     cgroupFile,
		cpus:           cpus,
		perfCollectors: sync.Map{},
		idEventMap:     make(map[uint64]string),
		mu:             sync.Mutex{},
		resultMap:      make(map[string]float64),
		valueCh:        make(chan perfValue),
		syscall6:       syscallFunc,
		closeCh:        make(chan struct{}),
	}
	var wg sync.WaitGroup
	wg.Add(len(cpus))
	var errMutex sync.Mutex
	for _, cpu := range cpus {
		go func(cpu int) {
			// create perf group
			var pc PerfCollector
			pc.syscall6 = syscallFunc
			pc.fds = make([]io.ReadCloser, 0, len(events)-1)
			pc.cpu = cpu
			defer wg.Done()
			attr := attrMap[events[0]]
			defaultFd := -1
			r1, _, e1 := collector.syscall6(syscall.SYS_PERF_EVENT_OPEN, uintptr(unsafe.Pointer(attr)),
				cgroupFile.Fd(), uintptr(cpu), uintptr(defaultFd), uintptr(unix.PERF_FLAG_PID_CGROUP|unix.PERF_FLAG_FD_CLOEXEC), uintptr(0))
			if !errors.Is(e1, syscall.Errno(0)) {
				errMutex.Lock()
				err = multierr.Append(err, fmt.Errorf("failed to create perf fd, Error: %s, cpu: %d, event: %s", unix.ErrnoName(e1), cpu, events[0]))
				errMutex.Unlock()
				return
			}
			leaderFd := os.NewFile(r1, fmt.Sprintf("%s_%d", events[0], cpu))
			pc.leaderFd = leaderFd
			var id uint64
			_, _, e1 = collector.syscall6(syscall.SYS_IOCTL, r1, uintptr(unix.PERF_EVENT_IOC_ID), uintptr(unsafe.Pointer(&id)), 0, 0, 0)
			if !errors.Is(e1, syscall.Errno(0)) {
				errMutex.Lock()
				err = multierr.Append(err, fmt.Errorf("failed to get perf id, Error: %s, cpu: %d, event: %s", unix.ErrnoName(e1), cpu, events[0]))
				errMutex.Unlock()
				return
			}
			collector.mu.Lock()
			collector.idEventMap[id] = events[0]
			collector.mu.Unlock()
			for i := 1; i < len(events); i++ {
				attr := attrMap[events[i]]
				r1, _, e1 := collector.syscall6(syscall.SYS_PERF_EVENT_OPEN, uintptr(unsafe.Pointer(attr)),
					cgroupFile.Fd(), uintptr(cpu), r1, uintptr(unix.PERF_FLAG_PID_CGROUP|unix.PERF_FLAG_FD_CLOEXEC), uintptr(0))
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to create perf fd, Error: %s, cpu: %d, event: %s", unix.ErrnoName(e1), cpu, events[i]))
					errMutex.Unlock()
					return
				}
				fd := os.NewFile(r1, fmt.Sprintf("%s_%d", events[i], cpu))
				pc.fds = append(pc.fds, fd)
				var id uint64
				_, _, e1 = collector.syscall6(syscall.SYS_IOCTL, r1, uintptr(unix.PERF_EVENT_IOC_ID), uintptr(unsafe.Pointer(&id)), 0, 0, 0)
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to get perf id, Error: %s, cpu: %d, event: %s", unix.ErrnoName(e1), cpu, events[i]))
					errMutex.Unlock()
					return
				}
				collector.mu.Lock()
				collector.idEventMap[id] = events[i]
				collector.mu.Unlock()
			}
			collector.perfCollectors.Store(cpu, &pc)
			runtime.LockOSThread()
		}(cpu)
	}
	wg.Wait()
	return collector, err
}

func ResetAndEnablePerfGroupCollector(collector *PerfGroupCollector, cGroupFile *os.File, cpus []int, events []string, syscallFunc func(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)) (*PerfGroupCollector, error) {
	var err error
	if len(events) == 0 {
		err = errors.New("events cannot be empty")
		return nil, err
	}
	var wg sync.WaitGroup
	wg.Add(len(cpus))
	var errMutex sync.Mutex
	for _, cpu := range cpus {
		go func(cpu int) {
			defer wg.Done()
			if pcInf, ok := collector.perfCollectors.Load(cpu); ok {
				pc := pcInf.(*PerfCollector)
				leaderFd := pc.leaderFd.(*os.File)
				_, _, e1 := collector.syscall6(syscall.SYS_IOCTL, leaderFd.Fd(), uintptr(unix.PERF_EVENT_IOC_RESET), uintptr(unix.PERF_IOC_FLAG_GROUP), 0, 0, 0)
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to reset perf group, Error: %s, cpu: %d, events: %s", unix.ErrnoName(e1), cpu, events))
					errMutex.Unlock()
					return
				}
				_, _, e1 = collector.syscall6(syscall.SYS_IOCTL, leaderFd.Fd(), uintptr(unix.PERF_EVENT_IOC_ENABLE), unix.PERF_IOC_FLAG_GROUP, 0, 0, 0)
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to enable perf group, Error: %s, cpu: %d, events: %s", unix.ErrnoName(e1), cpu, events))
					errMutex.Unlock()
					return
				}
				runtime.LockOSThread()
			} else {
				klog.Info("perf collector on cpu not found")
			}
		}(cpu)
	}
	wg.Wait()
	// collect and statistic perf result
	go func() {
		for value := range collector.valueCh {
			collector.resultMap[collector.idEventMap[value.ID]] += value.Value
		}
		// changzhi comment
		// Don't close this channel, because it will be used in the next round of perf collection
		//close(collector.closeCh)
	}()

	return collector, err
}

func XResetAndEnablePerfGroupCollector(collector *PerfGroupCollector) (*PerfGroupCollector, error) {
	var err error
	var wg sync.WaitGroup
	cpus := make([]int, runtime.GOMAXPROCS(0))
	for i := range cpus {
		cpus[i] = i
	}

	wg.Add(len(cpus))
	var errMutex sync.Mutex
	// renew the closeCh and valueCh
	collector.closeCh = make(chan struct{})
	collector.valueCh = make(chan perfValue)
	collector.resultMap = make(map[string]float64)
	for _, cpu := range cpus {
		go func(cpu int) {
			defer wg.Done()
			if pcInf, ok := collector.perfCollectors.Load(cpu); ok {
				pc := pcInf.(*PerfCollector)
				leaderFd := pc.leaderFd.(*os.File)
				_, _, e1 := collector.syscall6(syscall.SYS_IOCTL, leaderFd.Fd(), uintptr(unix.PERF_EVENT_IOC_RESET), uintptr(unix.PERF_IOC_FLAG_GROUP), 0, 0, 0)
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to reset perf group, Error: %s, cpu: %d", unix.ErrnoName(e1), cpu))
					errMutex.Unlock()
					return
				}
				_, _, e1 = collector.syscall6(syscall.SYS_IOCTL, leaderFd.Fd(), uintptr(unix.PERF_EVENT_IOC_ENABLE), unix.PERF_IOC_FLAG_GROUP, 0, 0, 0)
				if !errors.Is(e1, syscall.Errno(0)) {
					errMutex.Lock()
					err = multierr.Append(err, fmt.Errorf("failed to enable perf group, Error: %s, cpu: %d", unix.ErrnoName(e1), cpu))
					errMutex.Unlock()
					return
				}
				runtime.LockOSThread()
			} else {
				klog.Info("perf collector on cpu not found")
			}
		}(cpu)
	}
	wg.Wait()
	// collect and statistic perf result
	go func() {
		for value := range collector.valueCh {
			collector.mu.Lock()
			collector.resultMap[collector.idEventMap[value.ID]] += value.Value
			collector.mu.Unlock()
		}
		close(collector.closeCh)
	}()

	return collector, err
}

func GetAndStartPerfGroupCollectorOnContainer(cgroupFile *os.File, cpus []int, events []string) (*PerfGroupCollector, error) {
	cg, err := CreatePerfGroupCollector(cgroupFile, cpus, events, syscall.Syscall6)
	if err != nil {
		klog.Info("err ", err)
	}
	collector, err := ResetAndEnablePerfGroupCollector(cg, cgroupFile, cpus, events, syscall.Syscall6)

	if err != nil {
		return nil, err
	}
	return collector, nil
}

func GetContainerPerfResult(collector *PerfGroupCollector) (map[string]float64, error) {
	var err error
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(collector.cpus))
	for _, cpu := range collector.cpus {
		go func(cpu int) {
			defer wg.Done()
			if pcInf, ok := collector.perfCollectors.Load(cpu); ok {
				pc := pcInf.(*PerfCollector)
				if err = pc.collect(collector.valueCh); err != nil {
					mu.Lock()
					err = multierr.Append(err, err)
					mu.Unlock()
				}
			}
		}(cpu)
	}
	wg.Wait()
	// changzhi comment
	err = multierr.Append(err, collector.cleanUp())
	<-collector.closeCh

	return collector.resultMap, err
}

func GetContainerCyclesAndInstructionsGroup(collector *PerfGroupCollector) (float64, float64, float64, error) {
	resMap, err := GetContainerPerfResult(collector)
	if err != nil {
		return 0, 0, 0, err
	}
	return resMap[CYCLES], resMap[INSTRUCTIONS], resMap["LONGEST_LAT_CACHE.MISS"], nil
}

func (pgc *PerfGroupCollector) cleanUp() error {
	// err := c.cgroupFile.Close()
	// if err != nil {
	// 	return fmt.Errorf("close cgroupFile %v, err : %v", c.cgroupFile.Name(), err)
	// }
	close(pgc.valueCh)
	return nil
}

// caller must free the memory
func createPerfConfig(event string) (*unix.PerfEventAttr, error) {
	// https://pkg.go.dev/cmd/cgo OOM instread of check malloc error
	perfEventAttrPtr := C.malloc(C.ulong(unsafe.Sizeof(unix.PerfEventAttr{})))
	C.memset(perfEventAttrPtr, 0, C.ulong(unsafe.Sizeof(unix.PerfEventAttr{})))
	if err := pfmGetOsEventEncoding(event, perfEventAttrPtr); err != nil {
		return nil, err
	}

	return (*unix.PerfEventAttr)(perfEventAttrPtr), nil
}

// pfmPerfEncodeArgT represents structure that is used to parse perf event name
// into perf_event_attr using libpfm4.
type pfmPerfEncodeArgT struct {
	attr unsafe.Pointer
	fstr unsafe.Pointer
	size C.size_t
	_    C.int // idx
	_    C.int // cpu
	_    C.int // flags
}

// https://man7.org/linux/man-pages/man3/pfm_get_os_event_encoding.3.html
func pfmGetOsEventEncoding(event string, perfEventAttrPtr unsafe.Pointer) error {
	arg := pfmPerfEncodeArgT{}
	arg.attr = perfEventAttrPtr
	fstr := C.CString("")
	defer C.free(unsafe.Pointer(fstr))
	arg.size = C.ulong(unsafe.Sizeof(arg))
	eventCStr := C.CString(event)
	defer C.free(unsafe.Pointer(eventCStr))
	if err := C.pfm_get_os_event_encoding(eventCStr, C.PFM_PLM0|C.PFM_PLM3, C.PFM_OS_PERF_EVENT, unsafe.Pointer(&arg)); err != C.PFM_SUCCESS {
		return fmt.Errorf("failed to get event encoding: %v", err)
	}
	return nil
}

func (p *PerfCollector) collect(ch chan perfValue) error {
	// changzhi comment
	if err := p.stop(); err != nil {
		return err
	}
	bufPool := BufPools[len(p.fds)+1]
	buf := bufPool.Get().(*[]byte)
	defer bufPool.Put(buf)
	_, err := p.leaderFd.Read(*buf)
	if err != nil {
		return err
	}

	header := &perfValueHeader{}
	reader := bytes.NewReader(*buf)
	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		return err
	}
	scalingRatio := 1.0
	if header.TimeRunning != 0 && header.TimeEnabled != 0 {
		scalingRatio = float64(header.TimeRunning) / float64(header.TimeEnabled)
	}
	for i := 0; i < int(header.Nr); i++ {
		v := perfValuePool.Get().(*value)
		defer perfValuePool.Put(v)
		value := &perfValue{}
		if err := binary.Read(reader, binary.LittleEndian, v); err != nil {
			return err
		}
		value.Value = float64(v.Value) / scalingRatio
		value.ID = v.ID
		ch <- *value
	}
	return nil
}

// stop stops perf group counter
func (p *PerfCollector) stop() error {
	f := p.leaderFd.(*os.File)
	// DISABLE -> stop
	_, _, e1 := p.syscall6(syscall.SYS_IOCTL, f.Fd(), uintptr(unix.PERF_EVENT_IOC_DISABLE), unix.PERF_IOC_FLAG_GROUP, 0, 0, 0)
	if !errors.Is(e1, syscall.Errno(0)) {
		return errors.New(unix.ErrnoName(e1))
	}
	return nil
}

// close closes all perf fds
func (p *PerfCollector) close() error {
	var err error
	if leaderErr := p.leaderFd.Close(); leaderErr != nil {
		err = multierr.Append(err, leaderErr)
	}
	for _, fd := range p.fds {
		if closeErr := fd.Close(); closeErr != nil {
			err = multierr.Append(err, closeErr)
		}
	}
	return err
}

var (
	EventsMap = map[string][]string{
		"CPICollector": {"cycles", "instructions", "LONGEST_LAT_CACHE.MISS"},
	}
)

func createPerfContainerCollector(containerCgroupFilePath string, events []string) (*PerfGroupCollector, error) {
	cpus := make([]int, runtime.GOMAXPROCS(0))
	for i := range cpus {
		cpus[i] = i
	}
	containerCgroupFile, err := os.OpenFile(containerCgroupFilePath, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return nil, err
	}

	return CreatePerfGroupCollector(containerCgroupFile, cpus, events, syscall.Syscall6)
}

func main() {
	containerCgroupFilePath := os.Args[1]
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/kubepods-pod7bb341b7_92a7_4e41_919c_97b7e6fc68b7.slice/docker-48dcb0f76e8de43b8cd330b5e744aa92b54b47074712bd89cf382727b7999f54.scope"
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/kubepods-pod7bb341b7_92a7_4e41_919c_97b7e6fc68b7.slice"
	//containerCgroupFilePath := "/sys/fs/cgroup/perf_event/kubepods.slice/"
	LibInit()
	InitBufferPool()
	collectorOnSingleContainer, err := createPerfContainerCollector(containerCgroupFilePath, EventsMap["CPICollector"])
	if err != nil {
		klog.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	XResetAndEnablePerfGroupCollector(collectorOnSingleContainer)
	cycles, instructions, cacheMisses, err := GetContainerCyclesAndInstructionsGroup(collectorOnSingleContainer)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info("cache-misses: ", cacheMisses, " cycles: ", cycles, " instructions: ", instructions)
	LibFinalize()
}
