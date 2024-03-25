import {ComponentFixture, TestBed} from '@angular/core/testing';

import {TcrRoleComponent} from './tcr-role.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {Observable, Subject} from "rxjs";
import {TcrMessage} from "../../interfaces/tcr-message";
import {TcrRolesService} from "../../services/trc-roles.service";
import {TcrRole} from "../../interfaces/tcr-role";

class TcrRolesServiceFake implements Partial<TcrRolesService> {
  webSocket$: Subject<TcrMessage> = new Subject<TcrMessage>();

  getRole(name: string): Observable<TcrRole> {
    return new Observable<TcrRole>();
  }
}

describe('TcrRoleComponent', () => {
  let component: TcrRoleComponent;
  let fixture: ComponentFixture<TcrRoleComponent>;
  let serviceFake: TcrRolesService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrRoleComponent, HttpClientTestingModule],
      providers: [
        {provide: TcrRolesService, useClass: TcrRolesServiceFake},
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(TcrRoleComponent);
    component = fixture.componentInstance;
    serviceFake = TestBed.inject(TcrRolesService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

