import {Component, Input, OnInit} from '@angular/core';
import {TcrBuildInfo} from "../tcr-build-info";
import {TcrBuildInfoService} from "../tcr-build-info.service";
import {DatePipe, NgIf} from "@angular/common";

@Component({
  selector: 'app-tcr-build-info',
  standalone: true,
  imports: [
    NgIf,
    DatePipe
  ],
  templateUrl: './tcr-build-info.component.html',
  styleUrl: './tcr-build-info.component.css'
})
export class TcrBuildInfoComponent implements OnInit {
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
