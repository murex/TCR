import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrSessionInfoComponent } from './tcr-session-info.component';

describe('TcrSessionInfoComponent', () => {
  let component: TcrSessionInfoComponent;
  let fixture: ComponentFixture<TcrSessionInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrSessionInfoComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrSessionInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
