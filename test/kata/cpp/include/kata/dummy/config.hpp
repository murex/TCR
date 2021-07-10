#ifndef DUMMY_CONFIG
#define DUMMY_CONFIG

#ifdef _MSC_VER
#ifdef DUMMY_EXPORTS
#define DUMMY_API __declspec(dllexport)
#else
#define DUMMY_API __declspec(dllimport)
#endif
#else
#define DUMMY_API
#endif

#endif // DUMMY_CONFIG
