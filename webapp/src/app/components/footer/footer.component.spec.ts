/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
