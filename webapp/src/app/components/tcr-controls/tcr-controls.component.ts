import {Component} from '@angular/core';
import {NgForOf, NgIf} from "@angular/common";
import {TcrRoleComponent} from "../tcr-role/tcr-role.component";
import {TcrControlsService} from "../../services/tcr-controls.service";

@Component({
  selector: 'app-tcr-controls',
  standalone: true,
  imports: [
    NgForOf,
    TcrRoleComponent,
    NgIf
  ],
  templateUrl: './tcr-controls.component.html',
  styleUrl: './tcr-controls.component.css'
})
export class TcrControlsComponent {
  abortCommandDescription: string = `Abort Current Command`;

  constructor(private controlsService: TcrControlsService) {
  }

  abortCommand() {
    this.controlsService.abortCommand().subscribe(_ => {
      console.log(`Sent abort command request`);
    })
  }
}
