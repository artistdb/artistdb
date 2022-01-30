import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { ArtistComponent } from './artist.component';

describe('ArtistComponent', () => {
  let component: ArtistComponent;
  let fixture: ComponentFixture<ArtistComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ArtistComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtistComponent);
    component = fixture.componentInstance;
    expect(component).toBeDefined();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show artist input form upon clicking button', () => {
    const b = fixture.debugElement.query(By.css('.--new-artist'));
    b.triggerEventHandler('click', null);
    const form: HTMLElement = fixture.nativeElement.querySelector('app-artist-form')
    expect(form).toBeDefined();
  })
});
