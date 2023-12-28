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
  title = "TCR Session Information"
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
