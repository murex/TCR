import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterLink, RouterOutlet} from '@angular/router';
import {HttpClientModule} from "@angular/common/http";
import {TcrBuildInfoComponent} from "./tcr-build-info/tcr-build-info.component";
import {TcrSessionInfoComponent} from "./tcr-session-info/tcr-session-info.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    TcrBuildInfoComponent,
    TcrSessionInfoComponent,
    RouterLink,
    HttpClientModule,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'TCR - Test && Commit || Revert';
}
