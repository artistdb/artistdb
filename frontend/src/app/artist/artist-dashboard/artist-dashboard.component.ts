import { Component, Input, OnInit } from '@angular/core';
import { ArtistInput, MOCK_ARTISTS } from '../artist.model.component';

@Component({
  selector: 'app-artist-dashboard',
  templateUrl: './artist-dashboard.component.html',
  styleUrls: ['./artist-dashboard.component.css']
})
export class ArtistDashboardComponent implements OnInit {
  @Input() artists!: ArtistInput[];
  
  constructor() { 
  }

  ngOnInit(): void {
  }

}
