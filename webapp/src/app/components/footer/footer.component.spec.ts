import {ComponentFixture, TestBed} from '@angular/core/testing';

import {FooterComponent} from './footer.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
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

class FakeTcrBuildInfoService implements Partial<TcrBuildInfoService> {
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
      imports: [FooterComponent, HttpClientTestingModule],
      providers: [
        {provide: TcrBuildInfoService, useClass: FakeTcrBuildInfoService}
      ]
    }).compileComponents();

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
