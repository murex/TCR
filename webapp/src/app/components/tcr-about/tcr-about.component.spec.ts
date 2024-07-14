import {Observable, of} from 'rxjs';
import {TcrBuildInfoService} from '../../services/tcr-build-info.service';
import {TcrAboutComponent} from "./tcr-about.component";
import {ComponentFixture, TestBed} from "@angular/core/testing";
import {provideHttpClientTesting} from "@angular/common/http/testing";
import {TcrBuildInfo} from "../../interfaces/tcr-build-info";
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';

const sample: TcrBuildInfo = {
  version: "1.0.0",
  os: "some-os",
  arch: "some-arch",
  commit: "abc123",
  date: "2024-01-01T00:00:00Z",
  author: "some-author",
};

class FakeTcrBuildInfoService {
  buildInfo: TcrBuildInfo = sample;

  getBuildInfo(): Observable<TcrBuildInfo> {
    return of(this.buildInfo);
  }
}

describe('TcrAboutComponent', () => {
  let component: TcrAboutComponent;
  let fixture: ComponentFixture<TcrAboutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrAboutComponent],
      providers: [
        {provide: TcrBuildInfoService, useClass: FakeTcrBuildInfoService},
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting()
      ]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrAboutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have title "About TCR"', () => {
    expect(component.title).toEqual('About TCR');
  });

  it('should fetch TCR build info on init', () => {
    expect(component.buildInfo).toEqual(sample);
  });
});
