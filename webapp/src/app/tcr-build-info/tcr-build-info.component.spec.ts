import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TcrBuildInfoComponent } from './tcr-build-info.component';

describe('TcrBuildInfoComponent', () => {
  let component: TcrBuildInfoComponent;
  let fixture: ComponentFixture<TcrBuildInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrBuildInfoComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TcrBuildInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
