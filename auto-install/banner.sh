#!/usr/bin/env bash
#
# Copyright (c) 2025 Murex
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

print_tcr_banner() {
    green='\e[1;32m'
    yellow='\e[1;33m'
    white='\e[1;37m'
    red='\e[1;31m'
    italic='\e[3m'
    reset='\e[1;0m'

    printf "${yellow} ███████████████████████╗\n" >&2
    printf "${yellow} ╚══██╔═════════════════╝\n" >&2
    printf "${yellow}    ██║ ${green} ██████╗ ${red}██████╗ \n" >&2
    printf "${yellow}    ██║ ${green}██╔════╝ ${red}██╔══██╗\n" >&2
    printf "${yellow}    ██║ ${green}██║      ${red}██████╔╝\n" >&2
    printf "${yellow}    ██║ ${green}██║      ${red}██╔══██╗\n" >&2
    printf "${yellow}    ██║ ${green}╚██████╗ ${red}██║  ██║\n" >&2
    printf "${yellow}    ╚═╝ ${green} ╚═════╝ ${red}╚═╝  ╚═╝\n" >&2
    printf "${italic}${yellow} Test ${white}&& ${green}Commit ${white}|| ${red}Revert\n" >&2
    printf "${reset}" >&2
}

# Entry point
print_tcr_banner
