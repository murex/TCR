# ----------------------------------------------------------------------------
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
# ----------------------------------------------------------------------------

# ----------------------------------------------------------------------------
# Example Makefile for Kotlin language
# ----------------------------------------------------------------------------

junit.ver = 1.8.2
junit.jar = junit-platform-console-standalone-$(junit.ver).jar
maven.url = https://repo1.maven.org/maven2/org/junit/platform/junit-platform-console-standalone
junit.url = $(maven.url)/$(junit.ver)/$(junit.jar)

src.dir = src
build.dir = build
build.lib.dir = $(build.dir)/lib

src.main.dir = $(src.dir)/main/kotlin
src.test.dir = $(src.dir)/test/kotlin

build.main.dir = $(build.dir)/classes
build.test.dir = $(build.dir)/test-classes

src.main.files := $(shell find $(src.main.dir) -type f -name '*.kt')
src.test.files := $(shell find $(src.test.dir) -type f -name '*.kt')

build.main.files = $(src.main.files:$(src.main.dir)/%.kt=$(build.main.dir)/%.class)
build.test.files = $(src.test.files:$(src.test.dir)/%.kt=$(build.test.dir)/%.class)

# Classpath separator is platform-dependant
os.kernel := $(shell uname -s)
ifeq ($(os.kernel),Darwin)
  classpath.sep = ":"
else ifeq ($(os.kernel),Linux)
  classpath.sep = ":"
else
  classpath.sep = ";"
endif

.PHONY: build
build: build-main build-test

.PHONY: test
test: build-test
	-@echo "Running tests"
	-@java -jar $(build.lib.dir)/$(junit.jar) \
         --class-path $(build.main.dir) \
         --class-path $(build.test.dir) \
         --scan-class-path --details=tree

.PHONY: clean
clean:
	-@echo "Cleaning up"
	-@rm -rf $(build.dir)

.PHONY: build-main
build-main: $(build.main.files)

$(build.main.files): | $(build.main.dir)

$(build.main.dir)/%.class: $(src.main.dir)/%.kt
	-@echo "Compiling $(<F)"
	-@kotlinc -d $(build.main.dir) $<

.PHONY: build-test
build-test: build-main download-junit $(build.test.files)

$(build.test.files): | $(build.test.dir)

$(build.test.dir)/%.class: $(src.test.dir)/%.kt
	-@echo "Compiling $(<F)"
	-@kotlinc -classpath $(build.lib.dir)/$(junit.jar)$(classpath.sep)$(build.main.dir)$(classpath.sep)$(src.main.dir) -d $(build.test.dir) $<

.PHONY: download-junit
download-junit: $(build.lib.dir)/$(junit.jar)

$(build.lib.dir)/$(junit.jar): | $(build.lib.dir)
	-@echo "Downloading $(@F)"
	-@curl --silent --fail \
         --dump-header $(build.lib.dir)/$(junit.jar).header \
         --output $(build.lib.dir)/$(junit.jar) \
         --location $(junit.url)

$(build.lib.dir) $(build.main.dir) $(build.test.dir):
	-@echo "Creating directory $@"
	-@mkdir -p $@
