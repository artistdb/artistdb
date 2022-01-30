import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { ArtistDashboardComponent } from './artist-dashboard.component';
import { ArtistInput, MOCK_ARTISTS } from '../artist.model.component';

describe('ArtistDashboardComponent', () => {
  let component: ArtistDashboardComponent;
  let fixture: ComponentFixture<ArtistDashboardComponent>;
  let data: ArtistInput[] = MOCK_ARTISTS;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ArtistDashboardComponent ]    
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtistDashboardComponent);
    component = fixture.componentInstance;
    component.artists = data;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render all data', () => {
    const e = fixture.debugElement.queryAll(By.css('.artist'))
    expect(e.length).toEqual(data.length);
  })

  it('should select info on clicking on an artists name', () => {
    const e = fixture.debugElement.query(By.css('.artist'))
    e.triggerEventHandler('click', null);
    expect(component.selectedArtist).toEqual(data[0]);
  })

  it('should have a delete button', () => {
    const b = fixture.debugElement.query(By.css('.--delete-artist'));
    expect(b).toBeDefined();
  })
});
