cmake_minimum_required(VERSION 3.19.3)

project(HelloWorld CXX)

set(CMAKE_CXX_STANDARD 17)
enable_testing()

set_property(GLOBAL PROPERTY USE_FOLDERS ON)

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib CACHE PATH "Where to place compiled static libraries.")
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib CACHE PATH "Where to place compiled shared libraries.")
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin CACHE PATH "Where to place compiled executables.")

# We use CPM for 3rd-party dependencies management
include(cmake/CPM.cmake)
include(cmake/FindGoogleTest.cmake)

add_library(HelloWorld-shared SHARED)

target_compile_definitions(HelloWorld-shared PRIVATE HELLO_WORLD_EXPORTS)

target_sources(HelloWorld-shared
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src/HelloWorld.cpp
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include/HelloWorld.hpp
        ${CMAKE_CURRENT_SOURCE_DIR}/include/config.hpp
)

target_include_directories(HelloWorld-shared
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

add_executable(HelloWorld-test)

target_sources(HelloWorld-test
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/test/HelloWorldTest.cpp
)

target_link_libraries(HelloWorld-test
    PRIVATE
        HelloWorld-shared
        gtest_main)

add_test(
    NAME ctest-HelloWorld
    COMMAND HelloWorld-test --gtest_brief=1 --gtest_output=xml:test-results/HelloWorld-test.xml
)

set_property(DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR} PROPERTY VS_STARTUP_PROJECT HelloWorld-test)
set_property(DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR} PROPERTY XCODE_STARTUP_PROJECT HelloWorld-test)

set_property(TARGET HelloWorld-shared PROPERTY FOLDER "HelloWorld")
set_property(TARGET HelloWorld-test PROPERTY FOLDER "HelloWorld")

set_property(TARGET gtest gtest_main PROPERTY FOLDER "Third-Party/GoogleTest")

set_property(GLOBAL PROPERTY PREDEFINED_TARGETS_FOLDER "CMake Generated")
