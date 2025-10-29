#!/usr/bin/env bash
# Prints TCR Banner

print_tcr_banner() {
    green=$(printf '\e[1;32m')
    yellow=$(printf '\e[1;33m')
    white=$(printf '\e[1;37m')
    red=$(printf '\e[1;31m')
    italic=$(printf '\e[3m')
    reset=$(printf '\e[1;0m')

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
