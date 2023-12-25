import {Routes} from '@angular/router';
import {TcrSessionInfoComponent} from "./components/tcr-session-info/tcr-session-info.component";
import {TcrAboutComponent} from "./components/tcr-about/tcr-about.component";
import {HomeComponent} from "./components/home/home.component";

export const routes: Routes = [
  // {path: '', redirectTo: '/session', pathMatch: 'full'},
  {path: '', component: HomeComponent},
  {path: 'session', component: TcrSessionInfoComponent},
  {path: 'about', component: TcrAboutComponent},
];
