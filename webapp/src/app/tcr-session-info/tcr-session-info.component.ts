import {Component, Input, OnInit} from '@angular/core';
import {TcrSessionInfo} from "../tcr-session-info";
import {TcrSessionInfoService} from "../tcr-session-info.service";
import {DatePipe, NgIf} from "@angular/common";
import {OnOffPipe} from "../shared/pipes/on-off.pipe";
import {ShowEmptyPipe} from "../shared/pipes/show-empty.pipe";

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
