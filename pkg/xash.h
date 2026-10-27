#ifndef ENGINE_H
#define ENGINE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
	char server_name[64];
    char map_name[64];
    char next_map[64];
    int online_users;
    int max_players;
    char timeleft[64];
} engine_stats_t;

// Тип функции-коллбэка, которую мы зарегистрируем из Go
typedef void (*pfnUpdateGoStats)(engine_stats_t* stats);

// Функция в Си-движке для регистрации этого коллбэка
void Sys_SetStatsCallback(pfnUpdateGoStats func);

// Функция в Си-движке, которую будем дёргать из GO
void Sys_DropClientByAddr(unsigned char ip0, unsigned char ip1, unsigned char ip2, unsigned char ip3, unsigned short port);

typedef void( *pfnChangeGame )( const char *progname );

static void Sys_ChangeGame( const char *progname ) {}

int Host_Main( int argc, char **argv, const char *progname, int bChangeGame, pfnChangeGame func );

#ifdef __cplusplus
}
#endif

#endif // ENGINE_H