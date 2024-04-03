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
