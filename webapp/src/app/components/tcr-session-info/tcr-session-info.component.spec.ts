import {Observable, of} from 'rxjs';
import {ComponentFixture, TestBed} from "@angular/core/testing";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {TcrSessionInfo} from "../../interfaces/tcr-session-info";
import {TcrSessionInfoService} from "../../services/tcr-session-info.service";
import {TcrSessionInfoComponent} from "./tcr-session-info.component";

const sample: TcrSessionInfo = {
  baseDir: "/my/base/dir",
  commitOnFail: false,
  gitAutoPush: false,
  language: "java",
  messageSuffix: "my-suffix",
  toolchain: "gradle",
  vcsName: "git",
  vcsSession: "my VCS session",
  workDir: "/my/work/dir"
};

class FakeTcrSessionInfoService implements Partial<TcrSessionInfoService> {
  sessionInfo: TcrSessionInfo = sample;

  getSessionInfo(): Observable<TcrSessionInfo> {
    return of(this.sessionInfo);
  }
}

describe('TcrSessionInfoComponent', () => {
  let component: TcrSessionInfoComponent;
  let fixture: ComponentFixture<TcrSessionInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrSessionInfoComponent, HttpClientTestingModule],
      providers: [
        {provide: TcrSessionInfoService, useClass: FakeTcrSessionInfoService}
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(TcrSessionInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have title "TCR Session Information"', () => {
    expect(component.title).toEqual('TCR Session Information');
  });

  it('should fetch TCR session info on init', () => {
    expect(component.sessionInfo).toEqual(sample);
  });
});
