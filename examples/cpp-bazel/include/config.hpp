#ifndef HELLO_WORLD_CONFIG
#define HELLO_WORLD_CONFIG

#ifdef _MSC_VER
#ifdef HELLO_WORLD_EXPORTS
#define HELLO_WORLD_API __declspec(dllexport)
#else
#define HELLO_WORLD_API __declspec(dllimport)
#endif
#else
#define HELLO_WORLD_API
#endif

#endif // HELLO_WORLD_CONFIG
