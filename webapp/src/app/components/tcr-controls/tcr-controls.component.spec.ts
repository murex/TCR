import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrControlsComponent } from './tcr-controls.component';

describe('TcrControlsComponent', () => {
  let component: TcrControlsComponent;
  let fixture: ComponentFixture<TcrControlsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrControlsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
