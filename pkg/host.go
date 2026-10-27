package goxash3d_fwgs

/*
#cgo LDFLAGS: -L. -lxash -lpublic -lbuild_vcs -lm -lbacktrace -llibrary_suffix
#include "xash.h"
#include <stdlib.h>

// Сообщаем Си, что сама функция объявлена где-то в Go
extern void go_stats_callback_cgo(engine_stats_t* stats);

// Делаем инлайновой и статической, чтобы линкер не ругался на дубликаты
static inline int host_Main( int argc, char **argv, const char *progname, int bChangeGame ) {
	return Host_Main( argc, argv, progname, bChangeGame, Sys_ChangeGame);
}

// Тоже делаем статичной
static inline void register_engine_callbacks() {
	Sys_SetStatsCallback(go_stats_callback_cgo);
}
*/
import "C"
import (
	"sync"
	"unsafe"
)

// ServerStats хранит расширенную статистику, обновляемую движком
type ServerStats struct {
	ServerName  string  `json:"server_name"`
	CurrentMap  string  `json:"current_map"`
	NextMap     string  `json:"next_map"`
	OnlineUsers int     `json:"online_users"`
	MaxPlayers  int     `json:"max_players"`
	TimeLeft    string  `json:"timeleft"`
}

var (
	// StatsMutex защищает GlobalStats от гонок данных при чтении из HTTP/WebRTC
	StatsMutex  sync.RWMutex
	GlobalStats ServerStats
)


//export go_stats_callback_cgo
func go_stats_callback_cgo(stats *C.engine_stats_t) {
	StatsMutex.Lock()
	defer StatsMutex.Unlock()

	GlobalStats.ServerName = C.GoString(&stats.server_name[0])
	GlobalStats.CurrentMap = C.GoString(&stats.map_name[0])
	GlobalStats.NextMap = C.GoString(&stats.next_map[0])
	GlobalStats.OnlineUsers = int(stats.online_users)
	GlobalStats.MaxPlayers = int(stats.max_players)
	GlobalStats.TimeLeft = C.GoString(&stats.timeleft[0])
}

// GetServerStats позволяет твоему эндпоинту /ServerInfo безопасно забрать копию текущей статы
func GetServerStats() ServerStats {
	StatsMutex.RLock()
	defer StatsMutex.RUnlock()
	return GlobalStats
}

// HostMain Runs Xash3D main loop
func (x *Xash3D) HostMain(args []string, gameDir string, bChangeGame int) int {
	// Автоматически регистрируем мост в Си перед стартом главного цикла движка
	C.register_engine_callbacks()

	argsCount := len(args)
	argc := C.int(argsCount)
	argv := make([]*C.char, argsCount+1) // +1 for NULL terminator
	for i, arg := range args {
		argv[i] = C.CString(arg)
	}
	argv[argsCount] = nil

	cGameDir := C.CString(gameDir)
	defer func() {
		for _, arg := range argv {
			if arg != nil {
				C.free(unsafe.Pointer(arg))
			}
		}
		C.free(unsafe.Pointer(cGameDir))
	}()

	return int(C.host_Main(argc, &argv[0], cGameDir, C.int(bChangeGame))) //
}

// DropClient кикает/отключает игрока по его индексу слота (индексу сессии)
func DropClientByAddr(ip [4]byte, port uint16) {
	C.Sys_DropClientByAddr(
		C.uchar(ip[0]),
		C.uchar(ip[1]),
		C.uchar(ip[2]),
		C.uchar(ip[3]),
		C.ushort(port),
	)
}