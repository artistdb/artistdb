import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { of } from 'rxjs';
import { FieldBase } from '../dynamic-form/field-base';
import { FieldService } from '../dynamic-form/field.service';
import { ReactiveFormsModule, FormGroup } from '@angular/forms';

import { EventComponent } from './event.component';
import { DynamicFormComponent } from '../dynamic-form/dynamic-form.component';
import { DynamicFormFieldComponent } from '../dynamic-form/dynamic-form-field/dynamic-form-field.component';
import { ArtistDashboardComponent } from '../artist/artist-dashboard/artist-dashboard.component';
import { ArtistFormComponent } from '../artist/artist-form/artist-form.component';

describe('EventComponent', () => {
  
  let component: EventComponent;
  let fixture: ComponentFixture<EventComponent>;
  let fieldServiceStub: Partial<FieldService>;

  fieldServiceStub = {
    
    getFields() {
      var fields: FieldBase<string>[] = [
        { 
          key: 'Name',
          label: 'Name',
          value: '',
          required: true,
          controlType:'',
          type: ''
        }
      ];

      return of(fields);
    } 
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ ReactiveFormsModule ],
      declarations: [ ArtistDashboardComponent, 
        ArtistFormComponent, 
        DynamicFormComponent, 
        DynamicFormFieldComponent, 
        EventComponent ],
      providers: [ { provide: FieldService, use: fieldServiceStub } ]
    }) 
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EventComponent);
    component = fixture.componentInstance;
    fixture.debugElement.injector.get(FieldService)
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show event input form upon clicking button', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    fixture.detectChanges();
    const form: HTMLElement = fixture.nativeElement.querySelector('app-dynamic-form')
    expect(form).toBeTruthy();
  })

  it('should have a name field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    fixture.detectChanges();
    const e = fixture.debugElement.nativeElement;
    const f = e.querySelector('label[for="Name"]')
    expect(f).toBeTruthy()
  })

  it('should have a time field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    fixture.detectChanges();
    const e = fixture.debugElement.nativeElement;
    const f = e.querySelector('label[for="Start Date/Time"]')
    expect(f).toBeTruthy()
  })

  it('should have a location field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    fixture.detectChanges();
    const e = fixture.debugElement.nativeElement;
    const f = e.querySelector('label[for="Location"]')
    expect(f).toBeTruthy()
  })

  it('should have an artists field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    fixture.detectChanges();
    const e = fixture.debugElement.nativeElement;
    const f = e.querySelector('label[for="Invited Artists"]')
    expect(f).toBeTruthy()
  })
});
