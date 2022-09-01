import { Component } from '@angular/core';
import { Observable } from 'rxjs';

import { FieldBase } from '../dynamic-form/field-base';
import { FieldService } from '../dynamic-form/field.service';
import { ARTIST_FIELDS, MOCK_ARTISTS } from './artist.model';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.css'],
  providers: [FieldService]
})
export class ArtistComponent {
  fields$: Observable<FieldBase<any>[]>;
  xartists = MOCK_ARTISTS;
  newArtist = false;

  constructor(service: FieldService) {
    this.fields$ = service.getFields(ARTIST_FIELDS);
   }

  showArtistForm(): void {
    this.newArtist = true;
  }

}