import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl} from '@angular/forms';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.css']
})
export class ArtistComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

  artistForm = new FormGroup({
    firstName: new FormControl(''),
    lastName: new FormControl(''),
    artistName: new FormControl(''),
    pronouns: new FormControl(''),
    dateOfBirth: new FormControl(''),
    placeOfBirth: new FormControl(''),
    nationality: new FormControl(''),
    language:new FormControl(''),
    facebook: new FormControl(''),
    instagram: new FormControl(''),
    bandcamp: new FormControl(''),
    bioGerman: new FormControl(''),
    bioEnglish:new FormControl(''),
  });

}
