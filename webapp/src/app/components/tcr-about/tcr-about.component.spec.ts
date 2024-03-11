import {ComponentFixture, TestBed} from '@angular/core/testing';

import {TcrAboutComponent} from './tcr-about.component';

xdescribe('TcrAboutComponent', () => {
  let component: TcrAboutComponent;
  let fixture: ComponentFixture<TcrAboutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrAboutComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(TcrAboutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
