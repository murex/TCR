import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrTimerComponent } from './tcr-timer.component';

describe('TcrTimerComponent', () => {
  let component: TcrTimerComponent;
  let fixture: ComponentFixture<TcrTimerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrTimerComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrTimerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
