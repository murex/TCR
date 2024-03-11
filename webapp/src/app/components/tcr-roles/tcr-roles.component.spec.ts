import {ComponentFixture, TestBed} from '@angular/core/testing';

import {TcrRolesComponent} from './tcr-roles.component';

xdescribe('TcrRolesComponent', () => {
  let component: TcrRolesComponent;
  let fixture: ComponentFixture<TcrRolesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrRolesComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(TcrRolesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
