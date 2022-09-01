import { Component, Input, OnInit } from '@angular/core';
import { ArtistInput } from '../artist.model';

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
  
  selectedArtist?: ArtistInput;
  onSelect(artist: ArtistInput) {
    this.selectedArtist = artist;
  }
}
