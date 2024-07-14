import {ComponentFixture, TestBed} from '@angular/core/testing';

import {TcrControlsComponent} from './tcr-controls.component';
import {TcrControlsService} from "../../services/tcr-controls.service";
import {Observable, of} from "rxjs";

class FakeTcrControlsService {
  abortCommand(): Observable<unknown> {
    return of({});
  }
}

describe('TcrControlsComponent', () => {
  let component: TcrControlsComponent;
  let fixture: ComponentFixture<TcrControlsComponent>;
  let serviceFake: TcrControlsService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrControlsComponent],
      providers: [
        {provide: TcrControlsService, useClass: FakeTcrControlsService},
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    serviceFake = TestBed.inject(TcrControlsService);
    fixture = TestBed.createComponent(TcrControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });
  });

  describe('abort command button', () => {
    it('should trigger abort command from controls service', () => {
      const abortCommandFunction = spyOn(serviceFake, 'abortCommand').and.callThrough();
      // Trigger the abortCommand call
      component.abortCommand();
      // Verify that the service received the request
      expect(abortCommandFunction).toHaveBeenCalledTimes(1);
    });
  });
});
