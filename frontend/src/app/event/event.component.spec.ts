import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { EventComponent } from './event.component';

describe('EventComponent', () => {
  let component: EventComponent;
  let fixture: ComponentFixture<EventComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EventComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(EventComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show event input form upon clicking button', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    const form: HTMLElement = fixture.nativeElement.querySelector('app-event-form')
    expect(form).toBeDefined();
  })

  it('should have a name field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    const f = fixture.debugElement.query(By.css('label[for="name"]'))
    expect(f).toBeDefined()
  })

  it('should have a time field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    const f = fixture.debugElement.query(By.css('label[for="Start Date/Time"]'))
    expect(f).toBeDefined()
  })

  it('should have a location field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    const f = fixture.debugElement.query(By.css('label[for="Location"]'))
    expect(f).toBeDefined()
  })

  it('should have a artists field', () => {
    const b = fixture.debugElement.query(By.css('.--new-event'));
    b.triggerEventHandler('click', null);
    const f = fixture.debugElement.query(By.css('label[for="Invited Artists"]'))
    expect(f).toBeDefined()
  })
});
