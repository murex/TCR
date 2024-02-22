import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RouterLink, RouterModule, RouterOutlet} from '@angular/router';
import {HttpClientModule} from "@angular/common/http";
import {TcrAboutComponent} from "./components/tcr-about/tcr-about.component";
import {TcrSessionInfoComponent} from "./components/tcr-session-info/tcr-session-info.component";
import {HeaderComponent} from "./components/header/header.component";
import {FooterComponent} from "./components/footer/footer.component";
import {NotificationComponent} from "./components/notification/notification.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    TcrAboutComponent,
    TcrSessionInfoComponent,
    RouterLink,
    RouterModule,
    HttpClientModule,
    HeaderComponent,
    FooterComponent,
    NotificationComponent,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'TCR';
}
