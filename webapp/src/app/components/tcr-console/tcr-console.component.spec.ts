import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrConsoleComponent } from './tcr-console.component';

describe('TcrConsoleComponent', () => {
  let component: TcrConsoleComponent;
  let fixture: ComponentFixture<TcrConsoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrConsoleComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrConsoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
