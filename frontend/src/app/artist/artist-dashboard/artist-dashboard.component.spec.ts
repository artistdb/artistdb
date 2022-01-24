import { ComponentFixture, TestBed } from '@angular/core/testing';

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
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render all data', () => {
    component.artists = data;
    fixture.detectChanges();
    const e: HTMLElement = fixture.nativeElement;
    expect(e.querySelectorAll('li').length).toEqual(data.length);
  })
});
