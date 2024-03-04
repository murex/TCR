import {Component} from '@angular/core';
import {RouterLink} from "@angular/router";
import {TcrTimerComponent} from "../tcr-timer/tcr-timer.component";

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    RouterLink,
    TcrTimerComponent
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent {

}
