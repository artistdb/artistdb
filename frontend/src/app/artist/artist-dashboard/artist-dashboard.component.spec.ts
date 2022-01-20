import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ArtistDashboardComponent } from './artist-dashboard.component';
import { MOCK_ARTISTS } from '../artist.model.component';

describe('ArtistDashboardComponent', () => {
  let component: ArtistDashboardComponent;
  let fixture: ComponentFixture<ArtistDashboardComponent>;

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

  it('should render data', () => {
    
  })
});
