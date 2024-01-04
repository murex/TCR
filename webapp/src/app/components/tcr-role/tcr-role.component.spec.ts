import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrRoleComponent } from './tcr-role.component';

describe('TcrRoleComponent', () => {
  let component: TcrRoleComponent;
  let fixture: ComponentFixture<TcrRoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrRoleComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrRoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
