/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import {Component, Input, OnInit} from '@angular/core';
import {TcrSessionInfo} from "../../interfaces/tcr-session-info";
import {TcrSessionInfoService} from "../../services/tcr-session-info.service";
import {DatePipe, NgIf} from "@angular/common";
import {OnOffPipe} from "../../pipes/on-off.pipe";
import {ShowEmptyPipe} from "../../pipes/show-empty.pipe";

@Component({
  selector: 'app-tcr-session-info',
  standalone: true,
  imports: [
    DatePipe,
    NgIf,
    OnOffPipe,
    ShowEmptyPipe
  ],
  templateUrl: './tcr-session-info.component.html',
  styleUrl: './tcr-session-info.component.css'
})
export class TcrSessionInfoComponent implements OnInit {
  title: string = "TCR Session Information"
  @Input() sessionInfo?: TcrSessionInfo;

  constructor(
    private sessionInfoService: TcrSessionInfoService) {
  }

  ngOnInit(): void {
    this.getSessionInfo();
  }

  private getSessionInfo(): void {
    this.sessionInfoService.getSessionInfo()
      .subscribe(sessionInfo => this.sessionInfo = sessionInfo);
  }
}
