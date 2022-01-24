import { ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';

import { ArtistFormComponent } from './artist-form.component';
import { ArtistInput } from 'src/app/artist/artist.model.component';

describe('ArtistFormComponent', () => {
  let component: ArtistFormComponent;
  let fixture: ComponentFixture<ArtistFormComponent>;

  var testArtist: ArtistInput = {
    firstName: "Kai",
    lastName: "Uwe",
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ArtistFormComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtistFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('first name input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--first-name`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--first-name label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('first-name')
      expect(l.nativeElement.textContent).toEqual('First Name: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--first-name input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('firstName')
    });
  });

  describe('last name input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--last-name`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--last-name label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('last-name')
      expect(l.nativeElement.textContent).toEqual('Last Name: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--last-name input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('lastName')
    });
  });

  describe('artist name input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--artist-name`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--artist-name label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('artist-name')
      expect(l.nativeElement.textContent).toEqual('Artist Name: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--artist-name input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('artistName')
    });

    describe('pronouns input', () => {
      it('should be wrapped in a div with corresponding class', () => {
        const e = fixture.debugElement.query(By.css(`.ctrl.--pronouns`));
        expect(e).toBeTruthy();
      });
  
      it('should have a label', () => {
        const l = fixture.debugElement.query(By.css(`.--pronouns label`));
        expect(l).toBeTruthy(); 
        expect(l.nativeElement.getAttribute('for')).toEqual('pronouns')
        expect(l.nativeElement.textContent).toEqual('Pronouns: ')
      });
      
      it('should have an input field', () => {
        const i = fixture.debugElement.query(By.css(`.--pronouns input`));
        expect(i).toBeTruthy();
        expect(i.nativeElement.getAttribute('type')).toEqual('text');
        expect(i.nativeElement.getAttribute('formControlName')).toEqual('pronouns')
      });
    });
  });

  describe('Place of birth input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--place-of-birth`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--place-of-birth label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('place-of-birth')
      expect(l.nativeElement.textContent).toEqual('Place of Birth: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--place-of-birth input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('placeOfBirth')
    });
  });

  describe('Nationality input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--nationality`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--nationality label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('nationality')
      expect(l.nativeElement.textContent).toEqual('Nationality: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--nationality input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('nationality')
    });
  });

  describe('Language input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--language`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--language label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('language')
      expect(l.nativeElement.textContent).toEqual('Language: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--language input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('language')
    });
  });

  describe('Facebook url input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--facebook`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--facebook label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('facebook')
      expect(l.nativeElement.textContent).toEqual('Facebook: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--facebook input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('url');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('facebook')
    });
  });

  describe('Instagram input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--instagram`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--instagram label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('instagram')
      expect(l.nativeElement.textContent).toEqual('Instagram: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--instagram input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('url');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('instagram')
    });
  });

  describe('Bandcamp input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--bandcamp`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--bandcamp label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('bandcamp')
      expect(l.nativeElement.textContent).toEqual('Bandcamp: ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--bandcamp input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('url');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('bandcamp')
    });
  });

  describe('German Bio input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--bio-german`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--bio-german label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('bio-german')
      expect(l.nativeElement.textContent).toEqual('Your Bio (in German): ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--bio-german input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('bioGerman')
    });
  });

  describe('English Bio input', () => {
    it('should be wrapped in a div with corresponding class', () => {
      const e = fixture.debugElement.query(By.css(`.ctrl.--bio-english`));
      expect(e).toBeTruthy();
    });

    it('should have a label', () => {
      const l = fixture.debugElement.query(By.css(`.--bio-english label`));
      expect(l).toBeTruthy(); 
      expect(l.nativeElement.getAttribute('for')).toEqual('bio-english')
      expect(l.nativeElement.textContent).toEqual('Your Bio (in English): ')
    });
    
    it('should have an input field', () => {
      const i = fixture.debugElement.query(By.css(`.--bio-english input`));
      expect(i).toBeTruthy();
      expect(i.nativeElement.getAttribute('type')).toEqual('text');
      expect(i.nativeElement.getAttribute('formControlName')).toEqual('bioEnglish')
    });
  });

  describe('Submit Button', () => {   
    it('should have a submit button', () => {
      const b = fixture.debugElement.query(By.css(`button`));
      expect(b).toBeTruthy();
      expect(b.nativeElement.textContent).toEqual('Submit');
    });
    
    it('should not be possible to click the button when data is lacking', () => {
      const b = fixture.debugElement.query(By.css(`button`));
      expect(b.nativeElement.getAttribute('disabled')).toBe('');
    })

    it('should be possible to click button after entering necessary data', () => {
      const b = fixture.debugElement.query(By.css(`button`));
      const fn = component.artistForm.get('firstName');
      const ln = component.artistForm.get('lastName');
      fn?.setValue(testArtist.firstName);
      ln?.setValue(testArtist.lastName);
      fixture.detectChanges()
      expect(b.nativeElement.getAttribute('disabled')).toBe(null);
    })
  })
});
