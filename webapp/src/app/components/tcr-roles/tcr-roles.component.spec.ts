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

import { ComponentFixture, TestBed } from "@angular/core/testing";
import { TcrRolesComponent } from "./tcr-roles.component";
import { Component, Input } from "@angular/core";
import { By } from "@angular/platform-browser";
import { TcrRoleComponent } from "../tcr-role/tcr-role.component";
import { TcrRolesService } from "../../services/trc-roles.service";
import { Observable, of } from "rxjs";
import { TcrRole } from "../../interfaces/tcr-role";
import { TcrMessage } from "../../interfaces/tcr-message";

// Mock service for testing
class FakeTcrRolesService {
  message$ = new Observable<TcrMessage>();

  getRole(name: string): Observable<TcrRole> {
    return of({
      name: name,
      description: `${name} role`,
      active: false,
    });
  }

  activateRole(name: string, state: boolean): Observable<TcrRole> {
    return of({
      name: name,
      description: `${name} role`,
      active: state,
    });
  }
}

// Mock component for testing - use a different selector to avoid conflicts
@Component({
  selector: "app-mock-tcr-role",
  template: '<div class="mock-role">{{ name }}</div>',
  standalone: true,
})
class MockTcrRoleComponent {
  @Input() name = "";
}

describe("TcrRolesComponent", () => {
  let component: TcrRolesComponent;
  let fixture: ComponentFixture<TcrRolesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TcrRolesComponent],
      providers: [{ provide: TcrRolesService, useClass: FakeTcrRolesService }],
    })
      .overrideComponent(TcrRolesComponent, {
        remove: {
          imports: [TcrRoleComponent],
        },
        add: {
          imports: [MockTcrRoleComponent],
          template: `
            <section>
              <div class="container">
                <div class="row mbr-justify-content-center">
                  @for (role of roles; track role) {
                    <app-mock-tcr-role
                      [name]="role"
                      class="col-lg-6 mbr-col-md-10"></app-mock-tcr-role>
                  }
                </div>
              </div>
            </section>
          `,
        },
      })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrRolesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe("component instance", () => {
    it("should be created", () => {
      expect(component).toBeTruthy();
    });

    it("should have driver in list of roles", () => {
      expect(component.roles).toContain("driver");
    });

    it("should have navigator in list of roles", () => {
      expect(component.roles).toContain("navigator");
    });
  });

  describe("component DOM", () => {
    it("should contain 2 TcrRoleComponent children", () => {
      const roleElements = fixture.debugElement.queryAll(
        By.css("app-mock-tcr-role"),
      );
      expect(roleElements.length).toBe(2);

      // Check that the mock components received the correct inputs
      const mockComponents = roleElements.map(
        (el) => el.componentInstance as MockTcrRoleComponent,
      );
      expect(mockComponents[0].name).toBe("driver");
      expect(mockComponents[1].name).toBe("navigator");
    });
  });
});
