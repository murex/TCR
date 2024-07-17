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
import {TcrBuildInfo} from "../../interfaces/tcr-build-info";
import {TcrBuildInfoService} from "../../services/tcr-build-info.service";
import {DatePipe, NgIf} from "@angular/common";

@Component({
  selector: 'app-tcr-about',
  standalone: true,
  imports: [
    NgIf,
    DatePipe
  ],
  templateUrl: './tcr-about.component.html',
  styleUrl: './tcr-about.component.css'
})
export class TcrAboutComponent implements OnInit {
  title: string = "About TCR";
  @Input() buildInfo?: TcrBuildInfo;

  constructor(
    private buildInfoService: TcrBuildInfoService) {
  }

  ngOnInit(): void {
    this.getBuildInfo();
  }

  private getBuildInfo(): void {
    this.buildInfoService.getBuildInfo()
      .subscribe(buildInfo => this.buildInfo = buildInfo);
  }
}
