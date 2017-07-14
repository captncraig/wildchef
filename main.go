package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/JamesHovious/w32"
	"github.com/mitchellh/go-ps"
)

var handle w32.HANDLE
var itemNames map[string]string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dat, err := ioutil.ReadFile("items.json")
	e(err)
	e(json.Unmarshal(dat, &itemNames))
}

var ctx, cancel = context.WithCancel(context.Background())

func main() {
	go dumpToGithub()
	catchSigs()
	loadRecipes()
	prcs, err := ps.Processes()
	e(err)
	var pid uint32
	for _, prc := range prcs {
		if prc.Executable() == "CEMU.EXE" {
			pid = uint32(prc.Pid())
			fmt.Printf("Attaching to process %d\n", pid)
			break
		}
	}
	if pid == 0 {
		fmt.Println("Cemu Process Not Found")
		return
	}
	handle, err = w32.OpenProcess(w32.PROCESS_VM_OPERATION|w32.PROCESS_VM_READ|w32.PROCESS_VM_WRITE|w32.PROCESS_QUERY_INFORMATION, false, pid)
	e(err)
	defer w32.CloseHandle(handle)
	regionStart, regionSize := findRegionBySize()
	fmt.Printf("Memory Region found at 0x%0xd (0x%0xd bytes)\n", regionStart, regionSize)
	addr := findInventoryAddr(regionStart, regionSize)
	fmt.Printf("Inventory Address Found at 0x%0xd (%d)\n", addr, addr)
	items := readItems(addr, addr+regionSize)
	fmt.Printf("Found %d items\n", len(items))
	if len(items) > 0 {
		runAhk(genAhkClear(len(items)))
	}
	for i := 0; i < len(allRecipes); i++ {
		recipe := allRecipes[i]
		names := []string{}
		for _, ing := range recipe {
			names = append(names, ing.Name)
		}
		ingsIn := strings.Join(names, ",")
		if haveRecipe(ingsIn) {
			continue
		}
		fmt.Println("Cooking", i, ingsIn)
		runAhk(recipe.genAhk())
		time.Sleep(time.Second / 2)
		items := readItems(addr, addr+regionSize)
		found := false
		for _, item := range items {
			if !itemAddrs[item.addr] {
				res := getItemInfo(item.addr)
				fmt.Println(res)
				if res.IngredientsFromMem != ingsIn {
					log.Fatal("Ingredients no match. Uh oh.")
				}
				putRecipe(res, ingsIn)
				itemAddrs[item.addr] = true
				found = true
				break
			}
		}
		if !found {
			log.Fatal("Something wrong. Didn't find it.")
		}
		if len(items) > 14 {
			fmt.Println("Clearing Items")
			runAhk(genAhkClear(len(items)))
			itemAddrs = map[uint64]bool{}
			time.Sleep(time.Second / 2)
		}
	}
}

var itemAddrs = map[uint64]bool{}

const size = 0x8000

type ItemPtr struct {
	addr uint64
	id   string
}

func readItems(addr uint64, endAddr uint64) []ItemPtr {
	items := []ItemPtr{}
	d := func(a ...interface{}) {
		//fmt.Println(a...)
	}
	for ; addr < endAddr; addr += 0x220 {
		buf, err := w32.ReadProcessMemory(handle, addr, size)
		e(err)
		if offset := findSequenceMatch(buf, 0, targetSequence); offset != 0 {
			break
		}
		if isValidItem(buf, 8) {
			strBuf := buf[8:]
			i := 0
			for strBuf[i] != 0 {
				i++
			}
			name := string(strBuf[:i])
			if strings.HasPrefix(name, "Weapon") || strings.Contains(name, "Arrow") || strings.HasPrefix(name, "Armor") {
				continue
			}
			if !strings.HasPrefix(name, "Item_Cook") && !strings.HasPrefix(name, "Item_Roast") {
				continue
			}
			itemAddr := addr + 7
			d(itemAddr, name)
			items = append(items, ItemPtr{addr: itemAddr, id: name})
		}
	}
	//COST: 0x49
	//Bonus: 0x41

	return items
}

func e(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getItemInfo(addr uint64) Result {
	r := Result{}
	raw, err := w32.ReadProcessMemory(handle, addr-0x27, 600)
	e(err)
	r.Raw = raw

	r.ID = readString(addr + 1)
	r.Name = itemNames[r.ID]
	r.Cost = readInt(addr + 0x49)
	r.Hearts = readInt(addr + 0x41)
	r.Duration = readInt(addr + 0x45)
	r.Effect = readInt(addr+0x4d) >> 16
	r.EffectString = effects[r.Effect]
	r.Strength = readInt(addr+0x51) >> 16
	r.StrengthString = strengths[r.Strength]

	strs := []string{}
	for idx := uint64(0); idx < 5; idx++ {
		addr2 := addr + 0x75 + 76*idx
		iname := readString(addr2)
		if iname == "" {
			break
		}
		if itemNames[iname] == "" {
			log.Fatal("Unknown Item ID", iname)
		}
		strs = append(strs, itemNames[iname])
	}
	r.IngredientsFromMem = strings.Join(strs, ",")
	return r
}

const effNone uint32 = 0xbf80

var effects = map[uint32]string{
	0x4000: "Hearty",
	0x4080: "Chilly",
	0x40a0: "Spicy",
	0x40c0: "Electro",
	0x4120: "Mighty",
	0x4130: "Tough",
	0x4140: "Sneaky",
	0x4150: "Hasty",
	0x4160: "Energizing",
	0x4170: "Endura",
	0x4180: "Fireproof",
}

var strengths = map[uint32]string{
	0x3f80: "Low",
	0x4000: "Mid",
	0x4040: "High",
}

func readString(addr uint64) string {
	buf, err := w32.ReadProcessMemory(handle, addr, 50)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for buf[i] != 0 {
		i++
	}
	return string(buf[:i])
}

var cnt = 0

func dumpMem(addr uint64, size uint, name string) {
	buf, err := w32.ReadProcessMemory(handle, addr, size)
	if err != nil {
		log.Fatal(err)
	}
	//if name == "Elixir" {
	for i := uint(0); i < size; i++ {
		rel := int(i) - 0x27
		fmt.Printf("%0x(%d): %0x(%d) %s %d\n", rel, rel, buf[i], buf[i], string(buf[i:i+1]), addr+uint64(i))
	}
	//}
	cnt++
	ioutil.WriteFile(name+fmt.Sprint(cnt)+".dat", buf, 0777)
}

func readInt(addr uint64) uint32 {
	buf, err := w32.ReadProcessMemory(handle, addr, 4)
	if err != nil {
		log.Fatal(err)
	}
	return binary.BigEndian.Uint32(buf)
}

func isValidItem(buf []byte, idx int) bool {
	v := buf[idx]
	if v < 41 || v > 90 {
		return false
	}
	for i := idx + 1; i < len(buf) && buf[i] != 0; i++ {
		v := buf[idx]
		if v < 21 || v > 0x7e {
			return false
		}
	}
	return true
}

var targetSequence = []int{0x10, -1, -1, -1, 0, 0, 0, 0x40}

func findInventoryAddr(regionStart, regionSize uint64) uint64 {
	d := func(a ...interface{}) {
		//fmt.Println(a...)
	}

	addr := regionStart
	endAddr := addr + regionSize
	d(addr, endAddr, addr < endAddr)
	max := 62000
	iter := 0
	for addr < endAddr {
		if iter == max {
			break
		}
		iter++
		d("Reading addr", addr)
		buf, err := w32.ReadProcessMemory(handle, addr, size)
		if err != nil {
			log.Fatal(err)
		}
		offset := findSequenceMatch(buf, 0, targetSequence)
		if offset != -1 {
			d("FOUND Initial Sequence!", offset)
			addr += uint64(offset)
			d(buf[offset : offset+10])
			d("Looking deeper", addr, iter)
			buf, err = w32.ReadProcessMemory(handle, addr, size)
			if err != nil {
				log.Fatal(err)
			}
			d(buf[:4])
			if buf[1] == 30 || buf[1] == 0x1f {
				d("FOUND Second Match")
				offset2 := findSequenceMatch(buf, 0x220, targetSequence)
				d("OFFSET", offset2, iter)
				if offset2 == 0x220 {
					d("GOT IT!", addr)
					return addr
				}
				addr += 2
			} else {
				addr += 2
			}
		} else {
			addr += size - 2
		}
	}
	log.Fatal("Couldn't find inventory")
	return 0
}

func findSequenceMatch(buf []byte, start int, seq []int) int {
	maxSearch := len(buf) - len(seq)
	for idx := start; idx < maxSearch; idx++ {
		target := seq[0]
		at := buf[idx]
		if (target == -1 && at != 0) || (target >= 0 && at == byte(target)) {
			var bad bool
			for idx2 := 1; idx2 < len(seq); idx2++ {
				at = buf[idx2+idx]
				target = seq[idx2]
				if target == -1 || byte(target) == at {
					continue
				}
				bad = true
				break
			}
			if !bad {
				return idx
			}
		}
	}
	return -1
}

func findRegionBySize() (start, size uint64) {
	d := func(a ...interface{}) {
		//fmt.Println(a...)
	}
	const targetSize uint64 = 1416757248
	const maxAddr uint64 = 0x7fffffffffffffff
	var addr uint64
	var info w32.MEMORY_BASIC_INFORMATION
	for {

		d("READING", addr)
		res := w32.VirtualQueryEx(handle, uintptr(addr), &info, int(unsafe.Sizeof(info)))
		if res == 0 {
			log.Fatal("Problem with virtualQueryEx")
		}
		d(info.RegionSize, info, info.BaseAddress)
		base := addr
		if info.BaseAddress != nil {
			base = uint64(uintptr(info.BaseAddress))
		}
		size := uint64(info.RegionSize)
		if size == targetSize && info.Protect == w32.PAGE_READWRITE && info.State == w32.MEM_COMMIT {
			return base, size
		}
		addr = base + size
		if addr > maxAddr {
			log.Fatal("Suitable region not found")
		}
	}
}

func catchSigs() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		cancel()
		time.Sleep(50 * time.Millisecond)
		os.Exit(3)
	}()
}
