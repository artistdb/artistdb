import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';

import { FieldBase } from '../dynamic-form/field-base';
import { FieldService } from '../dynamic-form/field.service';
import { EVENT_FIELDS } from './event.model';

@Component({
  selector: 'app-event',
  templateUrl: './event.component.html',
  styleUrls: ['./event.component.css']
})

export class EventComponent implements OnInit {
  fields$: Observable<FieldBase<any>[]>;
  newEvent = false 
  
  constructor(service: FieldService) {
    this.fields$ = service.getFields(EVENT_FIELDS);
  }

  ngOnInit(): void {
  }

  showEventForm(): void {
    this.newEvent = true;
  }

}
