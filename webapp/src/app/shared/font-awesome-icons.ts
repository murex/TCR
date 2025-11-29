/*
TCR
Copyright (c) 2019-2025 Murex

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.
*/

import { FaIconLibrary } from "@fortawesome/angular-fontawesome";
import {
  faCogs,
  faCompass,
  faKeyboard,
  faGear,
  faFolderOpen,
  faFileCode,
  faCodeFork,
  faCircleExclamation,
  faClock,
  faIdCard,
  faDesktop,
  faQuestionCircle,
} from "@fortawesome/free-solid-svg-icons";
import { faGithub } from "@fortawesome/free-brands-svg-icons";

/**
 * Registers all FontAwesome icons used in the application
 * Call this function once in the app initialization
 */
export function registerFontAwesomeIcons(library: FaIconLibrary): void {
  // Add all icons used in the application
  library.addIcons(
    faCogs,
    faGithub,
    faCompass,
    faKeyboard,
    faGear,
    faFolderOpen,
    faFileCode,
    faCodeFork,
    faCircleExclamation,
    faClock,
    faIdCard,
    faDesktop,
    faQuestionCircle,
  );
}

/**
 * Icon name mappings for easy reference
 * Maps old FontAwesome 4.x names to new v7 names
 */
export const IconMappings = {
  "fa-cogs": "cogs",
  "fa-github": "github",
  "fa-compass": "compass",
  "fa-keyboard-o": "keyboard", // outline variant removed in v7
  "fa-gear": "gear",
  "fa-folder-open": "folder-open",
  "fa-file-code-o": "file-code", // outline variant removed in v7
  "fa-code-fork": "code-fork",
  "fa-exclamation-circle": "circle-exclamation", // renamed in v7
  "fa-clock-o": "clock", // outline variant removed in v7
  "fa-id-card-o": "id-card", // outline variant removed in v7
  "fa-desktop": "desktop",
  "fa-question-circle": "question-circle",
};
