import {ComponentFixture, TestBed} from '@angular/core/testing';
import {TcrRolesComponent} from './tcr-roles.component';
import {TcrRoleComponent} from "../tcr-role/tcr-role.component";
import {By} from "@angular/platform-browser";
import {MockComponent} from "ng-mocks";

describe('TcrRolesComponent', () => {
  let component: TcrRolesComponent;
  let fixture: ComponentFixture<TcrRolesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        TcrRolesComponent,
        MockComponent(TcrRoleComponent),
      ],
      providers: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TcrRolesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  describe('component instance', () => {
    it('should be created', () => {
      expect(component).toBeTruthy();
    });

    it('should have driver in list of roles', () => {
      expect(component.roles).toContain('driver');
    });

    it('should have navigator in list of roles', () => {
      expect(component.roles).toContain('navigator');
    });
  });

  describe('component DOM', () => {
    it('should contain 2 TcrRoleComponent children', () => {
      const roleElements = fixture.debugElement.queryAll(By.directive(TcrRoleComponent));
      expect(roleElements.length).toBe(2);
      expect(roleElements[0].componentInstance.name).toBe('driver');
      expect(roleElements[1].componentInstance.name).toBe('navigator');
    });
  });
});
