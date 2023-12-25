import {Component, Input} from '@angular/core';
import {TcrBuildInfo} from "../../interfaces/tcr-build-info";
import {TcrBuildInfoService} from "../../services/tcr-build-info.service";
import {DatePipe, NgIf} from "@angular/common";

@Component({
  selector: 'app-footer',
  standalone: true,
  imports: [
    NgIf,
    DatePipe
  ],
  templateUrl: './footer.component.html',
  styleUrl: './footer.component.css'
})
export class FooterComponent {
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
