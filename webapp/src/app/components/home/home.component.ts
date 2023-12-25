import {Component, OnInit} from '@angular/core';
import {Router} from "@angular/router";

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent implements OnInit {
  ngOnInit() {
  }

  constructor(private router: Router) {
  }

  navigateTo(url: string) {
    this.router.navigateByUrl(url)
      .then(r => {
        if (!r) {
          window.alert(`Page not found: ${url}`)
        }
      }
    );
  }
}
