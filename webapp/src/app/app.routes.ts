import {Routes} from '@angular/router';
import {TcrSessionInfoComponent} from "./tcr-session-info/tcr-session-info.component";
import {TcrBuildInfoComponent} from "./tcr-build-info/tcr-build-info.component";

export const routes: Routes = [
  {path: '', redirectTo: '/session', pathMatch: 'full'},
  {path: 'session', component: TcrSessionInfoComponent},
  {path: 'about', component: TcrBuildInfoComponent},
];
