import { Component, OnInit } from '@angular/core';
import { MOCK_ARTISTS } from '../artist.model.component';

@Component({
  selector: 'app-artist-dashboard',
  templateUrl: './artist-dashboard.component.html',
  styleUrls: ['./artist-dashboard.component.css']
})
export class ArtistDashboardComponent implements OnInit {
  artists = MOCK_ARTISTS;
  
  constructor() { }

  ngOnInit(): void {
  }

}
