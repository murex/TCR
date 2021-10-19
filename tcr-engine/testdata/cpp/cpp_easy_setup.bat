REM Copyright (c) 2021 Murex
REM
REM Permission is hereby granted, free of charge, to any person obtaining a copy
REM of this software and associated documentation files (the "Software"), to deal
REM in the Software without restriction, including without limitation the rights
REM to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
REM copies of the Software, and to permit persons to whom the Software is
REM furnished to do so, subject to the following conditions:
REM
REM The above copyright notice and this permission notice shall be included in all
REM copies or substantial portions of the Software.
REM
REM THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
REM IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
REM FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
REM AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
REM LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
REM OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
REM SOFTWARE.

setlocal

set OS=windows
set ARCH=x86_64
set ARCHIVE_EXTENSION=zip
set CMAKE_BIN_DIR=bin
set CMAKE=cmake.exe
set CTEST=ctest.exe
set CMAKE_GENERATOR_OPTIONS=-G "Visual Studio 15 2017 Win64"

set CMAKE_VERSION=3.21.0
set CMAKE_EXPECTED_DIR=cmake-%CMAKE_VERSION%-%OS%-%ARCH%
set CMAKE_EXPECTED_ARCHIVE_FILE=%CMAKE_EXPECTED_DIR%.%ARCHIVE_EXTENSION%
set CMAKE_ARCHIVE_URL="http://github.com/Kitware/CMake/releases/download/v%CMAKE_VERSION%/%CMAKE_EXPECTED_ARCHIVE_FILE%"
set CMAKE_HOME=cmake-%OS%-%ARCH%

set BUILD_DIR=build
if not exist %BUILD_DIR% (
    mkdir %BUILD_DIR%
)
pushd %BUILD_DIR%

set CMAKE_BUILD_DIR=cmake
if not exist %CMAKE_BUILD_DIR% (
    mkdir %CMAKE_BUILD_DIR%
)
pushd %CMAKE_BUILD_DIR%

if not exist %CMAKE_EXPECTED_ARCHIVE_FILE% (
	powershell -command "Invoke-WebRequest %CMAKE_ARCHIVE_URL% -OutFile %CMAKE_EXPECTED_ARCHIVE_FILE%"
    powershell -command "Expand-Archive -Force '%~dp0\%BUILD_DIR%\%CMAKE_BUILD_DIR%\%CMAKE_EXPECTED_ARCHIVE_FILE%' '%~dp0\%BUILD_DIR%\%CMAKE_BUILD_DIR%'"
	powershell -command "Rename-Item %CMAKE_EXPECTED_DIR% %CMAKE_HOME%"
)

pushd ..

set CMAKE_BIN_PATH=%CMAKE_BUILD_DIR%\%CMAKE_HOME%\%CMAKE_BIN_DIR%

%CMAKE_BIN_PATH%\%CMAKE% %CMAKE_GENERATOR_OPTIONS% -S .. -B .
%CMAKE_BIN_PATH%\%CMAKE% --build . --config Debug
%CMAKE_BIN_PATH%\%CTEST% --output-on-failure -C Debug

popd

popd

popd

