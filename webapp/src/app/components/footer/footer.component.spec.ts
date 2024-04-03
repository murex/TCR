import {ComponentFixture, TestBed} from '@angular/core/testing';

import {FooterComponent} from './footer.component';
import {TcrBuildInfo} from "../../interfaces/tcr-build-info";
import {TcrBuildInfoService} from "../../services/tcr-build-info.service";
import {Observable, of} from "rxjs";
import {DatePipe} from "@angular/common";

const sample: TcrBuildInfo = {
  version: "1.0.0",
  os: "",
  arch: "",
  commit: "",
  date: "2024-03-02T00:00:00Z",
  author: "",
};

class FakeTcrBuildInfoService {
  buildInfo: TcrBuildInfo = sample;

  getBuildInfo(): Observable<TcrBuildInfo> {
    return of(this.buildInfo);
  }
}

describe('FooterComponent', () => {
  let component: FooterComponent;
  let fixture: ComponentFixture<FooterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FooterComponent],
      providers: [
        {provide: TcrBuildInfoService, useClass: FakeTcrBuildInfoService}
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FooterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should show TCR version, year and month', () => {
    const element = fixture.nativeElement.querySelector('.footer-copyright');
    expect(element).toBeTruthy();
    const expected = "TCR version " + sample.version + " "
      + "("
      + new DatePipe('en-US').transform(sample.date, 'MMM yyyy')
      + ")";
    expect(element.textContent).toContain(expected);
  });

  it('should fetch TCR build info on init', () => {
    expect(component.buildInfo).toEqual(sample);
  });
});
