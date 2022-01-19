import { Component, Input, OnInit } from '@angular/core';
import { FormGroup, FormControl, Validators} from '@angular/forms';
import { UpsertArtists, ArtistInput } from '../graphql/artist.service';

@Component({
  selector: 'app-artist',
  templateUrl: './artist.component.html',
  styleUrls: ['./artist.component.css']
})
export class ArtistComponent implements OnInit {

  constructor(private upsertArtists: UpsertArtists) { }

  ngOnInit(): void {
  }

  artistForm = new FormGroup({
    firstName: new FormControl(''),
    lastName: new FormControl(''),
    artistName: new FormControl(''),
    pronouns: new FormControl(['']),
    // dateOfBirth: new FormControl(),
    placeOfBirth: new FormControl(''),
    nationality: new FormControl(''),
    language:new FormControl(''),
    facebook: new FormControl(''),
    instagram: new FormControl(''),
    bandcamp: new FormControl(''),
    bioGerman: new FormControl(''),
    bioEnglish:new FormControl(''),
  });

  @Input()
  artists: ArtistInput[] = [
    {
      firstName: this.artistForm.get('firstName')?.value,
      lastName: this.artistForm.get('lastName')?.value,
      artistName: this.artistForm.get('artistName')?.value,
      pronouns: this.artistForm.get('pronouns')?.value,
      // dateOfBirth: this.artistForm.get('dateOfBirth')?.value,
      placeOfBirth: this.artistForm.get('placeOfBirth')?.value,
      nationality: this.artistForm.get('nationality')?.value,
      language:this.artistForm.get('language')?.value,
      facebook: this.artistForm.get('facebook')?.value,
      instagram: this.artistForm.get('instagram')?.value,
      bandcamp: this.artistForm.get('bandcamp')?.value,
      bioGer: this.artistForm.get('bioGerman')?.value,
      bioEn:this.artistForm.get('bioEnglish')?.value, 
    }
  ]

  onSubmit() {  
    console.log(this.artists)
    this.upsertArtists
      .mutate({
        artists: this.artists,
      }).subscribe(
        ((error: any) => {
          console.error(error)
        }),
      );
  }
}