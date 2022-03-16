import { Component } from '@angular/core';
import { Observable } from 'rxjs';

import { FieldBase } from '../dynamic-form/field-base';
import { FieldService } from '../dynamic-form/field.service';
import { LOCATION_FIELDS } from './location.model';

@Component({
  selector: 'app-location',
  templateUrl: './location.component.html',
  styleUrls: ['./location.component.css'],
  providers: [FieldService]
})
export class LocationComponent {
  fields$: Observable<FieldBase<any>[]>;
  newLocation = false;

  constructor(service: FieldService) {
    this.fields$ = service.getFields(LOCATION_FIELDS);
  }

  showLocationForm(): void {
    this.newLocation = true;
  }
}