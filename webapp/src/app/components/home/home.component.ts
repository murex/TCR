import {Component} from '@angular/core';
import {Router} from "@angular/router";

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent {
  title: string = 'TCR - Test && Commit || Revert';

  constructor(private router: Router) {
  }

  async navigateTo(url: string): Promise<void> {
    this.router.navigateByUrl(url)
      .then(r => {
        if (!r)
          window.alert(`Page not found: ${url}`);
      });
  }
}
