import { Component, Input, OnInit } from '@angular/core';

import { UpsertArtists } from '../graphql/artist.service';
import { MOCK_ARTISTS } from './artist.model.component';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.css']
})
export class ArtistComponent implements OnInit {

  xartists = MOCK_ARTISTS;

  constructor() { }

  ngOnInit(): void {
  }

}