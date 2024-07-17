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
import {HeaderComponent} from './header.component';
import {RouterModule} from "@angular/router";
import {TcrTimerComponent} from "../tcr-timer/tcr-timer.component";
import {MockComponent} from "ng-mocks";

describe('HeaderComponent', () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        HeaderComponent,
        MockComponent(TcrTimerComponent),
        RouterModule.forRoot([]),
      ],
      providers: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });
  });

  describe('component DOM', () => {
    [
      {
        description: 'a navigation bar',
        selector: '.navbar'
      },
      {
        description: 'a brand element called TCR',
        selector: '.navbar-brand',
        text: 'TCR'
      },
      {
        description: 'a link to the home page',
        selector: 'a[href="/"]',
        text: 'Home',
      },
      {
        description: 'a link to the session page',
        selector: 'a[href="/session"]',
        text: 'Session',
      },
      {
        description: 'a link to the console page',
        selector: 'a[href="/console"]',
        text: 'Console',
      },
      {
        description: 'a link to the about page',
        selector: 'a[href="/about"]',
        text: 'About',
      },
      {
        description: 'a timer component',
        selector: 'app-tcr-timer',
      },
    ].forEach(({selector, description, text}) => {
      it(`should contain ${description}`, () => {
        const element = fixture.nativeElement.querySelector(selector);
        expect(element).toBeTruthy();
        if (text) {
          expect(element.textContent).toContain(text);
        }
      });
    });
  });
});
