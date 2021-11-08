cmake_minimum_required(VERSION 3.16)
project(ready C)

set(CMAKE_C_STANDARD 11)

file(GLOB_RECURSE LIB_SOURCES_SO RELATIVE ${CMAKE_SOURCE_DIR} "libs/*.so")
file(GLOB_RECURSE LIB_SOURCES RELATIVE ${CMAKE_SOURCE_DIR} "libs/*.c")
file(GLOB_RECURSE SRC_SOURCES RELATIVE ${CMAKE_SOURCE_DIR} "src/*.c")

message("-- Cmake Source Dirs: ${CMAKE_SOURCE_DIR}")
message("-- Including Libs Directories: ${LIB_SOURCES}")
message("-- Including Libs Directories: ${LIB_SOURCES_SO}")
message("-- Including Src Directories: ${SRC_SOURCES}")

include(${PROJECT_SOURCE_DIR}/conan/conanbuildinfo.cmake)
conan_basic_setup()

add_executable(ready main.c ${LIB_SOURCES} ${SRC_SOURCES})
target_link_libraries(ready ${CONAN_LIBS} ${CMAKE_SOURCE_DIR}/libs/ready/libready.so)